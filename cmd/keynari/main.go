package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/Daniel19931606/keynari/internal/engine"
	"github.com/Daniel19931606/keynari/internal/layout"
	"github.com/Daniel19931606/keynari/internal/macos"
	"github.com/Daniel19931606/keynari/internal/words"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "run" {
		runLive(os.Args[2:])
		return
	}

	text := flag.String("text", "", "text to process through the correction engine")
	noDict := flag.Bool("no-dict", false, "disable dictionary checks")
	aggressive := flag.Bool("aggressive", false, "correct obvious layout-punctuation words even when absent from the dictionary")
	trace := flag.Bool("trace", false, "print corrections to stderr")
	ruDictPath := flag.String("ru-dict", "", "extra newline-delimited Russian dictionary file")
	enDictPath := flag.String("en-dict", "", "extra newline-delimited English dictionary file")
	flag.Parse()

	ruDict := words.RussianFull()
	enDict := words.EnglishFull()

	if *ruDictPath != "" {
		extra, err := words.FromFile(*ruDictPath)
		if err != nil {
			log.Fatalf("load Russian dictionary: %v", err)
		}
		ruDict = words.Merge(ruDict, extra)
	}

	if *enDictPath != "" {
		extra, err := words.FromFile(*enDictPath)
		if err != nil {
			log.Fatalf("load English dictionary: %v", err)
		}
		enDict = words.Merge(enDict, extra)
	}

	e := engine.New(
		layout.NewConverter(),
		ruDict,
		enDict,
		engine.Options{
			MinWordLength: 3,
			UseDictionary: !*noDict,
			Aggressive:    *aggressive,
		},
	)

	for _, r := range *text {
		printCorrections(*trace, e.Type(r))
	}
	printCorrections(*trace, e.Flush())

	fmt.Println(e.Text())
}

func runLive(args []string) {
	if runtime.GOOS != "darwin" {
		log.Fatal("live mode currently supports macOS only")
	}

	fs := flag.NewFlagSet("run", flag.ExitOnError)
	aggressive := fs.Bool("aggressive", true, "correct obvious wrong-layout words even when absent from the dictionary")
	trace := fs.Bool("trace", true, "print live corrections")
	quiet := fs.Bool("quiet", false, "disable live correction logs")
	logFile := fs.String("log-file", "", "write logs to file")
	_ = fs.Parse(args)

	if *quiet {
		*trace = false
	}
	if *logFile != "" {
		file, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("open log file: %v", err)
		}
		defer file.Close()
		log.SetOutput(file)
	}

	e := engine.New(
		layout.NewConverter(),
		words.RussianFull(),
		words.EnglishFull(),
		engine.Options{
			MinWordLength: 3,
			UseDictionary: true,
			Aggressive:    *aggressive,
		},
	)

	listener := macos.NewListener()
	replacer := macos.NewReplacer()
	var replacing atomic.Bool

	if err := listener.Start(); err != nil {
		log.Fatal(err)
	}
	defer listener.Stop()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("stopping")
		listener.Stop()
		os.Exit(0)
	}()

	if !*quiet {
		log.Println("Keynari live mode is running. Press Ctrl+C to stop.")
	}

	for event := range listener.Events() {
		if replacing.Load() {
			continue
		}
		if event.Modifiers.Meta || event.Modifiers.Ctrl || event.Modifiers.Alt {
			continue
		}

		if event.KeyCode == macos.KeyBackspace {
			e.Backspace()
			continue
		}
		if event.Char == 0 {
			continue
		}

		corrections := e.Type(event.Char)
		if len(corrections) == 0 {
			continue
		}

		correction := corrections[len(corrections)-1]
		oldLen := correction.ReplaceLen
		newText := correction.Corrected
		if correction.TypedLen > 0 && correction.LiveText != "" {
			oldLen = correction.TypedLen
			newText = correction.LiveText
		}

		replacing.Store(true)
		if err := replacer.Replace(oldLen, newText); err != nil {
			log.Printf("replace error: %v", err)
		}
		time.Sleep(80 * time.Millisecond)
		replacing.Store(false)

		if *trace {
			log.Printf("%q -> %q", correction.Original, correction.Corrected)
		}
	}
}

func printCorrections(enabled bool, corrections []engine.Correction) {
	if !enabled {
		return
	}

	for _, correction := range corrections {
		fmt.Fprintf(os.Stderr, "%q -> %q [%d:%d]\n",
			correction.Original,
			correction.Corrected,
			correction.ReplaceFrom,
			correction.ReplaceTo,
		)
	}
}
