# Keynari

**Keynari** is a local macOS keyboard-layout fixer for people who type in more than one language and hate discovering `ghbdtn rfr ltkf` after the sentence is already gone.

It watches your typing locally, detects words typed in the wrong layout, and replaces them in the active app:

```text
ghbdtn rfr ndjb ltkf  ->  привет как твои дела
руддщ ьн акшутв       ->  hello my friend
```

No cloud. No accounts. No telemetry. Just a fast local engine and a huge offline dictionary.

Русская версия ниже: [README на русском](#keynari-ru).

## Status

Keynari is early but already usable on macOS:

- Live correction in active apps through Accessibility permissions
- RU <-> EN keyboard layout conversion
- Large embedded dictionaries
- CLI test mode
- Double-clickable macOS `.app` bundle
- Local-only processing

## Why

Existing layout switchers often miss words, break punctuation, or feel opaque. Keynari is built around a testable correction engine first, then a macOS runner on top. The goal is simple: type fast, let the tool quietly clean up wrong-layout words.

## Install From Source

## Download The App

You can download a ready-to-run macOS build from GitHub Releases:

[Download Keynari for macOS](https://github.com/Daniel19931606/keynari/releases/latest)

After downloading:

1. Unzip `Keynari-macOS-*.zip`
2. Move `Keynari.app` to Applications if you want
3. Open it once
4. Grant Accessibility access in **System Settings > Privacy & Security > Accessibility**
5. Restart Keynari

Keynari runs as a menu bar app. It does not open a window; use the Keynari icon in the macOS menu bar to see that it is running or to quit it.

Current builds are unsigned. If macOS blocks the first launch, open **System Settings > Privacy & Security** and allow Keynari manually.

## Install From Source

Requirements:

- macOS
- Go 1.24+
- Xcode Command Line Tools

Clone and test:

```bash
git clone https://github.com/Daniel19931606/keynari.git
cd keynari
go test ./...
```

Build the CLI:

```bash
go build -o bin/keynari ./cmd/keynari
```

Run live mode:

```bash
./bin/keynari run
```

Build the macOS app:

```bash
./scripts/build_app.sh
open dist/Keynari.app
```

## macOS Permissions

Keynari needs Accessibility access to listen to keyboard events and replace text in the active app.

If macOS blocks it:

1. Open **System Settings**
2. Go to **Privacy & Security**
3. Open **Accessibility**
4. Add and enable either your terminal app or `dist/Keynari.app`
5. Restart Keynari

## CLI Playground

You can test the correction engine without live keyboard access:

```bash
go run ./cmd/keynari --aggressive --text "ghbdtn rfr ndjb ltkf"
```

Expected:

```text
привет как твои дела
```

Trace replacements:

```bash
go run ./cmd/keynari --aggressive --trace --text "ghbdtn,hfccrf;b vyt rfr ndjb ltkf"
```

## Dictionaries

Keynari ships with embedded offline dictionaries:

- Russian frequency list: 99,996 entries
- OpenCorpora Russian morphology forms: 3,064,348 entries
- English word list: 370,105 entries
- SCOWL English speller data: 320,287 entries

See [NOTICE.md](NOTICE.md) for attribution and licenses.

## Privacy

Keynari processes text locally. It does not send your typing anywhere.

Current local behavior:

- Keyboard events are read through macOS Accessibility APIs
- Corrections are computed in memory
- The live app can write logs to `~/Library/Logs/Keynari.log`
- No analytics, no network calls, no telemetry

## Roadmap

- Menubar UI
- Start at login
- Toggle per app
- User dictionary editor
- Signed release builds
- Windows and Linux runners
- More languages and layout pairs

## Development

Run tests:

```bash
go test ./...
```

Build app:

```bash
./scripts/build_app.sh
```

The core engine lives in `internal/engine`. The macOS event tap and replacement layer lives in `internal/macos`.

## License

Code is licensed under the project license in [LICENSE](LICENSE). Dictionary data has separate attribution in [NOTICE.md](NOTICE.md).

---

# Keynari RU

**Keynari** — локальное macOS-приложение, которое исправляет слова, набранные не в той раскладке.

Печатаешь:

```text
ghbdtn rfr ndjb ltkf
```

Получаешь:

```text
привет как твои дела
```

И в обратную сторону:

```text
руддщ ьн акшутв  ->  hello my friend
```

Все работает локально: без облака, аккаунтов, аналитики и телеметрии.

## Состояние

Keynari уже можно запускать локально на macOS:

- исправление текста прямо в активном приложении;
- RU <-> EN;
- большие встроенные словари;
- CLI-режим для тестов;
- `.app`, который можно открыть двойным кликом;
- вся обработка локальная.

## Установка из исходников

## Скачать приложение

Готовую macOS-сборку можно скачать из GitHub Releases:

[Скачать Keynari для macOS](https://github.com/Daniel19931606/keynari/releases/latest)

После скачивания:

1. Распакуй `Keynari-macOS-*.zip`
2. Перенеси `Keynari.app` в Applications, если хочешь
3. Открой приложение один раз
4. Выдай доступ Accessibility в **System Settings > Privacy & Security > Accessibility**
5. Перезапусти Keynari

Keynari работает как приложение в верхней строке меню macOS. Окно не открывается; по иконке Keynari в menu bar видно, что приложение запущено, там же его можно закрыть.

Текущие сборки не подписаны Apple Developer сертификатом. Если macOS заблокирует первый запуск, открой **System Settings > Privacy & Security** и разреши запуск Keynari вручную.

## Установка из исходников

Нужно:

- macOS;
- Go 1.24+;
- Xcode Command Line Tools.

Клонировать и проверить:

```bash
git clone https://github.com/Daniel19931606/keynari.git
cd keynari
go test ./...
```

Собрать CLI:

```bash
go build -o bin/keynari ./cmd/keynari
```

Запустить live-режим:

```bash
./bin/keynari run
```

Собрать `.app`:

```bash
./scripts/build_app.sh
open dist/Keynari.app
```

## Разрешения macOS

Keynari нужен доступ Accessibility, чтобы читать ввод и заменять текст в активном приложении.

Если macOS не разрешает запуск:

1. Открой **System Settings**
2. Перейди в **Privacy & Security**
3. Открой **Accessibility**
4. Добавь и включи Terminal или `dist/Keynari.app`
5. Перезапусти Keynari

## Проверка в терминале

```bash
go run ./cmd/keynari --aggressive --text "ghbdtn rfr ndjb ltkf"
```

Ожидаемый результат:

```text
привет как твои дела
```

Показать замены:

```bash
go run ./cmd/keynari --aggressive --trace --text "ghbdtn,hfccrf;b vyt rfr ndjb ltkf"
```

## Словари

Встроены офлайн-словари:

- русский частотный список: 99 996 строк;
- OpenCorpora, русские словоформы: 3 064 348 строк;
- английский список слов: 370 105 строк;
- SCOWL English speller data: 320 287 строк.

Атрибуция и лицензии словарей: [NOTICE.md](NOTICE.md).

## Приватность

Keynari ничего не отправляет наружу.

Сейчас:

- события клавиатуры читаются через macOS Accessibility API;
- исправления считаются в памяти;
- live-приложение может писать лог в `~/Library/Logs/Keynari.log`;
- аналитики, сети и телеметрии нет.

## План

- иконка в меню macOS;
- автозапуск при входе;
- включение/выключение для отдельных приложений;
- редактор пользовательского словаря;
- подписанные релизные сборки;
- Windows и Linux;
- новые языки и пары раскладок.
