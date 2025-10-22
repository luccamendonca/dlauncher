# dlauncher

> Nothing to do with `dmenu` [[ref](https://tools.suckless.org/dmenu/)], but I couldn't think of a better name.

dlauncher is a (very) simple "launcher" that aims to replicate the bookmark keywords feature brought by Firefox [[ref](https://support.mozilla.org/en-US/kb/bookmarks-firefox#w_how-to-use-keywords-with-bookmarks)], but allowing us to:
- use it directly from our terminal (or even with a GUI!)
- set hotkeys/shortcuts direclty to the pages you want
- use templating to replace any `%s` occurrences on the `template` parameter with the provided arguments

## Dependencies

- Whatever [Zenity](https://github.com/ncruces/zenity) depends on.
On Linux you need either `zenity`, `matedialog` or [qarma](https://github.com/luebking/qarma), but if you're running any Linux distro that has a GUI (i.e. _any_ DE), chances are that you already have at least one of them.
  - If you don't have one of the above, just install it and you should be good to go (via `apt`, `yum`, `yay`, etc.)

- go >= 1.23.4

## How it works

First things first, you need to create a config file containing:
1. The available executables (such as chrome, firefox, or any other browser you want to use) and
2. Your shortcuts

Use the `config.yaml` on the root dir of this repository as an example.
The default location for the config file is `~/.config/dlauncher/dlauncher.yaml`, but you may create it anywhere you want and use the `DLAUNCHER_CONFIG_PATH` env var to tell dlauncher where your config file is located.

### Running

For testing (or development) purposes, you can run it directly with: `go run main.go` and use the flags as described on the `help` section.

For a more permanent solution one could build the app and move somewhere on your PATH. For example:
```shell
$ go build main.go
```
```shell
$ sudo mv main /usr/bin/dlauncher
```

## Features

- **Single URL launching**: Open individual URLs using shortcuts and templates
- **Multi-URL launching**: Open multiple URLs simultaneously in new browser tabs
- **Template support**: Use `%s` placeholders in shortcuts for dynamic content
- **Multiple browsers**: Support for Chrome, Firefox, and other configurable browsers
- **CLI and GUI modes**: Use command line or graphical interface
- **Customizable shortcuts**: Create your own bookmark-like shortcuts

## Commands

### `run` - Launch single URL
Opens a single URL using a configured shortcut.

### `multi-run` - Launch multiple URLs
Opens multiple URLs in new browser tabs simultaneously. Perfect for:
- Opening daily reading lists
- Bulk opening bookmarks
- Research sessions with multiple sources

### `add` - Add new shortcut
Add a new shortcut to your configuration.

### `list` - List shortcuts
Display all configured shortcuts.

### Examples

> The following examples use the provided sample `config.yaml`.

#### Single URL Examples

1. Opens the `any` shortcut using the `firefox` executable. Asks for params, since the shortcut's `template` contains one or more `%s`, indicating that it is an actual template.
```shell
$ dlauncher run -e firefox -s any
Params for template, comma separated
https://google.com/
```
The result is that a new tab is opened on Firefox on the Google search page.

2. Opens the `google` shortcut using the `chrome` executable. Asks for params.
```shell
$ dlauncher run -e chrome -s google
Params for template, comma separated
my search
```
The result is that a new tab is opened on Chrome with the following URL: `https://www.google.com/search?q=my%20search`

3. Opens the `blank` shortcut using the `chrome` executable.
```shell
$ dlauncher run -e chrome -s blank
```
The result is that a new tab is opened on Chrome with the following URL: `about:blank`

#### Multiple URLs Examples

1. Open multiple URLs using command line parameter:
```shell
$ dlauncher multi-run -e chrome -l "https://www.google.com
https://www.github.com
https://www.stackoverflow.com"
```

2. Open multiple URLs interactively:
```shell
$ dlauncher multi-run -e firefox
Enter links (enter empty line when done):
https://www.google.com
https://www.github.com
https://www.stackoverflow.com

```

3. Open URLs from a file:
```shell
$ dlauncher multi-run -e chrome -l "$(cat my_links.txt)"
```

See `MULTI_RUN_USAGE.md` for detailed documentation on the multi-run feature.
