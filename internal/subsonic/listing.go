// Handles listings of artists and albums
package subsonic

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

func ListSongs(songs []Song) []byte {
	text := []string{}
	for _, i := range songs {
		s := []string{
			i.ID, i.Title, strconv.Itoa(i.Track),
			strconv.Itoa(i.DiscNumber),
			i.Suffix,
			strconv.Itoa(i.BitRate) + "Kbps",
			fmt.Sprint(time.Duration(i.Duration * int(time.Second))),
			strconv.Itoa(i.Size%1024%1024) + "MB",
			i.Parent,
		}
		text = append(text, strings.Join(s, "\t"))
	}
	return []byte(strings.Join(text, "\n"))
}

func ListAlbums(albums []Album) []byte {
	text := []string{}
	for _, i := range albums {
		s := []string{
			i.ID, i.Name, strconv.Itoa(i.Year), i.Genre,
			strconv.Itoa(i.SongCount), fmt.Sprint(time.Duration(i.Duration * int(time.Second))),
			i.Parent,
		}
		text = append(text, strings.Join(s, "\t"))
	}
	return []byte(strings.Join(text, "\n"))
}

func ListArtists(artists []Artist) []byte {
	text := []string{}
	for _, i := range artists {
		s := []string{i.ID, i.Name, strconv.Itoa(i.AlbumCount)}
		text = append(text, strings.Join(s, "\t"))
	}
	return []byte(strings.Join(text, "\n"))
}

// Remove character indexing and get all artists
func (indexes *Artists) ExtractArtists() []Artist {
	artists := []Artist{}
	for _, i := range indexes.Index {
		artists = slices.Concat(artists, i.Artists)
	}
	return artists
}
