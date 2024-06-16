# aquamarine

A minimal Subsonic CLI utility

## What it does

- Query a Subsonic library for artists, albums, and tracks.
- Queue tracks and albums to a detached headless [mpv](https://mpv.io/) instance.
- Scrobble tracks and albums (manually).
- Allow easy integration with [fzf](https://github.com/junegunn/fzf) (see [Usage](#usage)).

## What it doesn't do

- Act as a background daemon or music player.
- Provide an interactive interface for browsing.
- Scrobble tracks automatically.

## Configuration

By default, configuration is expected to live under
`$XDG_CONFIG_HOME/aquamarine/config.json`. Override this with `aquamarine -c <path_to_config> ...`

Available fields:

```json
{
  "username": "user",
  "password": "password",
  "password_cmd": "echo password",
  "host": "https://example.org"
}
```

- `username`: subsonic username
- `password`: plain text subsonic password
- `password_cmd`: alternative a shell command that returns the password to stdout
- `host`: subsonic endpoint (without `/rest`)

## Usage

**Manual listing**:

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

**Fzf**:

`aquamarine -c <path_to_config> interactive`

This allows you to select an artist and display their albums, the latter you
can queue with `ctrl-a` and scrobble with `ctrl-s`. Alternative selecting an
album displays the tracks, which you can individually queue by accepting them
(`enter`) or scrobble with `ctrl-s`.

## Integration with playerctl

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
