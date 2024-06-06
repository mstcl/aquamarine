// Make listings of artists and albums
package subsonic

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

func ListSongsPretty(songs []Song) []byte {
	t := []string{}
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
		t = append(t, strings.Join(s, "\t"))
	}
	return []byte(strings.Join(t, "\n"))
}

func ListAlbumsPretty(albums []Album) []byte {
	t := []string{}
	for _, i := range albums {
		s := []string{
			i.ID, i.Name, strconv.Itoa(i.Year), i.Genre,
			strconv.Itoa(i.SongCount), fmt.Sprint(time.Duration(i.Duration * int(time.Second))),
			i.Parent,
		}
		t = append(t, strings.Join(s, "\t"))
	}
	return []byte(strings.Join(t, "\n"))
}

func ListArtistsPretty(artists []Artist) []byte {
	t := []string{}
	for _, i := range artists {
		s := []string{i.ID, i.Name, strconv.Itoa(i.AlbumCount)}
		t = append(t, strings.Join(s, "\t"))
	}
	return []byte(strings.Join(t, "\n"))
}

// Remove character indexing and get all artists
func (indexes *Artists) ExtractArtists() []Artist {
	artists := []Artist{}
	for _, i := range indexes.Index {
		artists = slices.Concat(artists, i.Artists)
	}
	return artists
}
