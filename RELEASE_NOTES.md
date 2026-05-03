# Keynari v0.1.0

## Download / Скачать

Download the macOS app:

**Keynari-macOS-v0.1.0.zip**

Скачай готовое приложение для macOS:

**Keynari-macOS-v0.1.0.zip**

## What Is Keynari?

Keynari is a local macOS app that fixes words typed in the wrong keyboard layout:

```text
ghbdtn rfr ndjb ltkf  ->  привет как твои дела
руддщ ьн акшутв       ->  hello my friend
```

No cloud. No accounts. No telemetry.

## Что Это?

Keynari исправляет слова, набранные не в той раскладке:

```text
ghbdtn rfr ndjb ltkf  ->  привет как твои дела
руддщ ьн акшутв       ->  hello my friend
```

Все работает локально: без облака, аккаунтов и телеметрии.

## How To Install

1. Download `Keynari-macOS-v0.1.0.zip`
2. Unzip it
3. Move `Keynari.app` to Applications
4. Open `Keynari.app`
5. Grant Accessibility permission in **System Settings > Privacy & Security > Accessibility**
6. Restart Keynari

If you rebuild Keynari locally, remove the old Accessibility entry and add the new app again because macOS may treat the rebuilt binary as a different client.

Keynari runs in the macOS menu bar. It does not open a window; use the Keynari menu bar icon to see that it is running or to quit it.

## Как Установить

1. Скачай `Keynari-macOS-v0.1.0.zip`
2. Распакуй архив
3. Перенеси `Keynari.app` в Applications
4. Открой `Keynari.app`
5. Выдай Accessibility-доступ в **System Settings > Privacy & Security > Accessibility**
6. Перезапусти Keynari

Если ты пересобрал Keynari локально, удали старую запись в Accessibility и добавь новую сборку снова, потому что macOS может считать пересобранный бинарник другим клиентом.

Keynari работает в верхней строке меню macOS. Окно не открывается; по иконке Keynari в menu bar видно, что приложение запущено, там же его можно закрыть.

## macOS Security Note

This build is ad-hoc signed but not notarized by Apple yet. If macOS blocks the first launch, open **System Settings > Privacy & Security** and allow Keynari manually.

Сборка подписана ad-hoc подписью, но пока не notarized через Apple. Если macOS заблокирует первый запуск, открой **System Settings > Privacy & Security** и разреши Keynari вручную.
