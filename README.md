# aquamarine

A minimal Subsonic CLI utility

## What it does

* Query a Subsonic library for artists, albums, and tracks.
* Queue tracks and albums to a detached headless [mpv](https://mpv.io/) instance.
* Scrobble tracks and albums (manually).
* Allow easy integration with [fzf](https://github.com/junegunn/fzf) (see [Integration](#integration)).

## What it doesn't do

* Act as a background daemon or music player.
* Provide an interactive interface for browsing.
* Scrobble tracks automatically.

## Configuration

By default, configuration is expected to live under
`$XDG_CONFIG_HOME/aquamarine/config.yml`. Override this with `aquamarine -c
<path_to_config> ...`

Available fields:

```yaml
username: "user"                   # subsonic username
password: "password"               # plain text subsonic password
password_cmd: "echo password"      # alternative a shell command that returns the password to stdout
host: "https://example.org"        # subsonic endpoint (without `/rest`)
```

## Usage

**Listing**:

`aquamarine -c <path_to_config> [artists|albums|songs] ls [<empty>|<album ID>|<track ID>]`

Flags:

```
Usage of ls:
  -j	Format output as JSON
  -q	Don't print to stdout
  -s	Sync the cache
```

**Scrobble:**

`aquamarine -c <path_to_config> scrobble <album/track ID>`

**Queue**:

`aquamarine -c <path_to_config> queue <album/track ID>`

## Examples

### Listing artists

No arguments are excepted, as this lists all the artists in your library. For example:

```sh
$ aquamarine artists ls
```

The default output fields are `artist id`, `artist name`, `album count` (tab delimited).

### Listing albums

The argument takes an artist id. For example:

```sh
$ aquamarine albums ls "ar-182"
```

The default output fields are `album id`, `album name`, `album year`, `album
genre`, `track count`, `duration`, `artist id` (tab delimited).

### Listing tracks

The argument takes an album id. For example:

```sh
$ aquamarine songs ls "al-462"
```

The default output fields are `track id`, `track name`, `track number`, `disc
number`, `file format`, `bit rate`, `duration`, `file size`, `album id` (tab
delimited).

## Integration

### fzf

See [./scripts/fzf_picker.sh](./scripts/fzf_picker.sh) for a sample script
using fzf as a simple interactive browser.

This allows you to select an artist and display their albums, the latter you
can queue with `ctrl-a` and scrobble with `ctrl-s`. Alternative selecting an
album displays the tracks, which you can individually queue by accepting them
(`enter`) or scrobble with `ctrl-s`.

### playerctl

With [mpv-mpris](https://github.com/hoyon/mpv-mpris), we can use
[playerctl](https://github.com/altdesktop/playerctl) to play/pause, go to next
track or previous track, stop the queue, query song metadata, etc.

## Why

Most Subsonic daemons and music player out there expect you to have their stuff
running in a terminal or as a graphical process; this is unnecessarily
interactive if you listen to music in the background while doing other things
and just want to occasionally queue the next album and scrobble the previous
one, and you don't care for unnecessary resource usage beyond a barebones
player.

Existing tools can solve this problem (mpv, fzf, mpv-mpris, playerctl). For
once, mpv can output audio "headless", fzf can turn long lists into a fast,
interactive, fuzzy searcher, and the rest gives you the ability to play/pause,
etc.

By marrying these tools together, aquamarine provides a more interactive
experience than manually inputting/dragging the URLs into MPV, and a less
interactive experience than having a full blown daemon sitting somewhere on
your desktop.
