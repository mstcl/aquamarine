// Handles listings of artists and albums
package subsonic

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/mstcl/aquamarine/internal/ansi"
)

func FormatSongs(songs []Song, color bool, interactive bool) string {
	text := []string{}
	if !interactive {
		header := []string{"ID", "Track", "Title", "Disc", "Suffix", "Bitrate", "Duration", "Size", "Album ID"}
		headerLine := []string{}
		for _, i := range header {
			headerLine = append(headerLine, strings.Repeat("-", len(i)))
		}
		text = append(text, []string{strings.Join(header, "\t"), strings.Join(headerLine, "\t")}...)
	}
	for _, i := range songs {
		track := strconv.Itoa(i.Track)
		title := i.Title
		if color && interactive {
			track = ansi.BrightWhite + track + ansi.Reset
		}
		s := []string{
			i.ID, track, title,
			strconv.Itoa(i.DiscNumber),
			i.Suffix,
			strconv.Itoa(i.BitRate) + "Kbps",
			fmt.Sprint(time.Duration(i.Duration * int(time.Second))),
			strconv.Itoa(i.Size%1024%1024) + "MB",
			i.Parent,
		}
		text = append(text, strings.Join(s, "\t"))
	}
	return strings.Join(text, "\n")
}

func FormatAlbums(albums []Album, color bool, interactive bool) string {
	text := []string{}
	if !interactive {
		header := []string{"ID", "Name", "Year", "Genre", "# Songs", "Duration", "Artist ID"}
		headerLine := []string{}
		for _, i := range header {
			headerLine = append(headerLine, strings.Repeat("-", len(i)))
		}
		text = append(text, []string{strings.Join(header, "\t"), strings.Join(headerLine, "\t")}...)
	}
	for _, i := range albums {
		s := []string{
			i.ID, i.Name, strconv.Itoa(i.Year), i.Genre,
			strconv.Itoa(i.SongCount), fmt.Sprint(time.Duration(i.Duration * int(time.Second))),
			i.Parent,
		}
		text = append(text, strings.Join(s, "\t"))
	}
	return strings.Join(text, "\n")
}

func FormatArtists(artists []Artist, color bool, interactive bool) string {
	text := []string{}
	if !interactive {
		header := []string{"ID", "Name", "# Albums"}
		headerLine := []string{}
		for _, i := range header {
			headerLine = append(headerLine, strings.Repeat("-", len(i)))
		}
		text = append(text, []string{strings.Join(header, "\t"), strings.Join(headerLine, "\t")}...)
	}
	for _, i := range artists {
		s := []string{i.ID, i.Name, strconv.Itoa(i.AlbumCount)}
		text = append(text, strings.Join(s, "\t"))
	}
	return strings.Join(text, "\n")
}

// Remove character indexing and get all artists
func (indexes *Artists) ExtractArtists() []Artist {
	artists := []Artist{}
	for _, i := range indexes.Index {
		artists = slices.Concat(artists, i.Artists)
	}
	return artists
}
