// Handles responses and their json structs
package subsonic

import (
	"encoding/json"
	"io"
	"net/http"
)

type Song struct {
	// Album    string `json:"album"`
	// AlbumID  string `json:"albumId"`
	// Artist   string `json:"artist"`
	// ArtistID string `json:"artistId"`
	// DisplayArtist string `json:"displayArtist"`
	// DisplayAlbumArtist string    `json:"displayAlbumArtist"`
	// CoverArt      string `json:"coverArt"`
	// Year          int       `json:"year"`
	// Created       time.Time `json:"created"`
	// MusicBrainzID string `json:"musicBrainzId"`
	// ContentType string `json:"contentType"`
	// Path       string `json:"path"`
	// IsVideo    bool   `json:"isVideo"`

	// Prefix should be `tr-`
	ID string `json:"id"`

	// flac/mp3 etc.
	Suffix string `json:"suffix"`

	// Main song title
	Title string `json:"title"`

	// Parent album
	Parent string `json:"parent"`

	// Should be music
	Type string `json:"type"`

	Track      int `json:"track"`
	DiscNumber int `json:"discNumber"`
	BitRate    int `json:"bitRate"`

	// In seconds
	Duration int `json:"duration"`

	// In bytes
	Size int `json:"size"`

	// Important to know we can stream this
	IsDir bool `json:"isDir"`
}

type Album struct {
	// Artist        string    `json:"artist"`
	// CoverArt      string    `json:"coverArt"`
	// Title         string    `json:"title"`
	// PlayCount     int       `json:"playCount"`
	// Album         string    `json:"album"`
	// DisplayArtist string `json:"displayArtist"`
	// Created time.Time `json:"created"`

	// Prefix should be `al-`
	ID string `json:"id"`

	// Parent artist
	Parent string `json:"artistId"`

	// UTF-8 name
	Name string `json:"name"`

	Genre string `json:"genre"`

	Songs []Song `json:"song"`

	Year      int `json:"year"`
	SongCount int `json:"songCount"`

	// In seconds
	Duration int `json:"duration"`
}

type Artist struct {
	// CoverArt string `json:"coverArt"`

	// Prefix should be `ar-`
	ID string `json:"id"`

	// UTF-8 name
	Name string `json:"name"`

	// Albums
	Albums []Album `json:"album"`

	// Number of albums
	AlbumCount int `json:"albumCount"`
}

type Index struct {
	Name    string   `json:"name"`
	Artists []Artist `json:"artist"`
}

type Artists struct {
	IgnoredArticles string  `json:"ignoredArticles"`
	Index           []Index `json:"index"`
}

type SubsonicResponse struct {
	// Should returns "ok"
	Status string `json:"status"`

	// Should return server version
	Version string `json:"version"`

	// Should return server name, e.g. "gonic"
	Type string `json:"type"`

	// Should return subsonic API version
	ServerVersion string `json:"serverVersion"`

	// Relevant for GetArtists()
	Artists Artists `json:"artists,omitempty"`

	// Relevant for GetArtist()
	Artist Artist `json:"artist,omitempty"`

	// Relevant for GetAlbum()
	Album Album `json:"album,omitempty"`

	// Should return true
	OpenSubsonic bool `json:"openSubsonic"`
}

type Response struct {
	SubsonicResponse SubsonicResponse `json:"subsonic-response"`
}

// Get and unmarshal http response
func (c *SubsonicConnection) getResponse(requestUrl string) (*Response, error) {
	res, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var decodedBody Response
	if err := json.Unmarshal(body, &decodedBody); err != nil {
		return nil, err
	}

	return &decodedBody, nil
}
