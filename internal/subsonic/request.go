// Main request API calls for interacting with Subsonic
// Loosely based on stmps
// https://github.com/spezifisch/stmps
package subsonic

import (
	"bytes"
	"encoding/json"
	"net/url"
	"os"
	"strings"

	"github.com/mstcl/aquamarine/internal/file"
)

type SubsonicConnection struct {
	Username string
	Password string
	Host     string
}

// The default query
func defaultQuery(c *SubsonicConnection) url.Values {
	query := url.Values{}
	token, salt := authToken(c.Password)
	query.Set("t", token)
	query.Set("s", salt)
	query.Set("u", c.Username)
	query.Set("v", "1.15.0")
	query.Set("c", "aquamarine")
	query.Set("f", "json")

	return query
}

type Formatter[T any, C bool, I bool] interface {
	func(T, C, I) string
}

// Generic function to handle response
func handleResponse[T any, C bool, I bool, F Formatter[T, C, I]](
	loc string,
	displayRaw bool,
	quiet bool,
	fn F,
	data T,
	color C,
	interactive I,
) (string, error) {
	buf := new(bytes.Buffer)
	bufEncoder := json.NewEncoder(buf)
	bufEncoder.SetIndent("", "  ")

	if err := bufEncoder.Encode(data); err != nil {
		return "", err
	}

	if err := file.Cache(buf.Bytes(), loc); err != nil {
		return "", err
	}

	if quiet {
		return "", nil
	}

	if displayRaw {
		return buf.String(), nil
	}

	return fn(data, color, interactive), nil
}

// Generic function to handle cache data
func handleCache[T any, C bool, I bool, F Formatter[T, C, I]](
	loc string,
	displayRaw bool,
	fn F,
	data T,
	color C,
	interactive I,
) (string, error) {
	file, err := os.Open(loc)
	if err != nil {
		return "", err
	}

	defer file.Close()

	if displayRaw {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(file)
		return buf.String(), nil
	}

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return "", err
	}

	return fn(data, color, interactive), nil
}

type fn func(string) string

// Generic function to handle a slice of songs
func handleSongs(
	id string,
	f fn,
) ([]string, error) {
	loc := file.GetCachePath(id)

	songs := []Song{}

	file, err := os.Open(loc)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	if err := json.NewDecoder(file).Decode(&songs); err != nil {
		return nil, err
	}

	t := []string{}

	for _, i := range songs {
		t = append(t, f(i.ID))
	}
	return t, nil
}

// Ping the server
// func (c *SubsonicConnection) Ping() (*Response, error) {
// 	query := defaultQuery(c)
// 	requestUrl := c.Host + "/rest/ping" + "?" + query.Encode()
// 	return c.getResponse(requestUrl)
// }

// Scrobble a submission, where the id is a track id
// https://www.subsonic.org/pages/api.jsp#scrobble
func (c *SubsonicConnection) Scrobble(id string) (err error) {
	query := defaultQuery(c)
	query.Set("id", id)

	requestUrl := c.Host + "/rest/scrobble" + "?" + query.Encode()
	_, err = c.getResponse(requestUrl)
	return err
}

// Get all artists, similar to indexes
// https://www.subsonic.org/pages/api.jsp#getArtists
//
// Always cache, set TTL 1 week. Force refetch with --sync
func (c *SubsonicConnection) GetArtists(
	sync bool, displayRaw bool, quiet bool, color bool, interactive bool,
) (string, error) {
	loc := file.IndexCacheLoc
	shouldSync, err := file.ShouldSync(loc)
	if err != nil {
		return "", err
	}

	artists := []Artist{}

	if sync || shouldSync {
		query := defaultQuery(c)
		requestUrl := c.Host + "/rest/getArtists" + "?" + query.Encode()

		res, err := c.getResponse(requestUrl)
		if err != nil {
			return "", err
		}

		// Transform into FlattenArtists
		indexes := res.SubsonicResponse.Artists
		artists = indexes.ExtractArtists()

		return handleResponse(loc, displayRaw, quiet, FormatArtists, artists, color, interactive)
	}

	return handleCache(loc, displayRaw, FormatArtists, artists, color, interactive)
}

// Get albums from an artist
// https://www.subsonic.org/pages/api.jsp#getArtist
//
// # The id string should have an `ar-` prefix
//
// Always cache, set TTL 1 week. Force refetch with --sync
func (c *SubsonicConnection) GetArtist(
	id string, sync bool, displayRaw bool, quiet bool, color bool, interactive bool,
) (string, error) {
	loc := file.GetCachePath(id)
	shouldSync, err := file.ShouldSync(loc)
	if err != nil {
		return "", err
	}

	albums := []Album{}

	if sync || shouldSync {
		query := defaultQuery(c)
		query.Set("id", id)
		requestUrl := c.Host + "/rest/getArtist" + "?" + query.Encode()

		res, err := c.getResponse(requestUrl)
		if err != nil {
			return "", err
		}

		// Encode and cache response
		albums := res.SubsonicResponse.Artist.Albums

		return handleResponse(loc, displayRaw, quiet, FormatAlbums, albums, color, interactive)
	}

	return handleCache(loc, displayRaw, FormatAlbums, albums, color, interactive)
}

// Get songs from an album
// https://www.subsonic.org/pages/api.jsp#getAlbum
//
// # The id string should have an `al-` prefix
//
// Always cache, set TTL 1 week. Force refetch with --sync
func (c *SubsonicConnection) GetAlbum(
	id string, sync bool, displayRaw bool, quiet bool, color bool, interactive bool,
) (string, error) {
	loc := file.GetCachePath(id)
	shouldSync, err := file.ShouldSync(loc)
	if err != nil {
		return "", err
	}

	songs := []Song{}

	if sync || shouldSync {
		query := defaultQuery(c)
		query.Set("id", id)
		requestUrl := c.Host + "/rest/getAlbum" + "?" + query.Encode()

		res, err := c.getResponse(requestUrl)
		if err != nil {
			return "", err
		}

		// Encode and cache response
		songs := res.SubsonicResponse.Album.Songs

		return handleResponse(loc, displayRaw, quiet, FormatSongs, songs, color, interactive)
	}

	return handleCache(loc, displayRaw, FormatSongs, songs, color, interactive)
}

// Formats and constructs stream url, ignoring directories
//
// The id is prefixed with `tr`, e.g. something like `tr-1293`
// https://www.subsonic.org/pages/api.jsp#stream
func (c *SubsonicConnection) getPlayUrl(id string) string {
	query := defaultQuery(c)
	query.Set("id", id)

	return c.Host + "/rest/stream" + "?" + query.Encode()
}

// Given ID, get song url
// If album, we have to parse it
// Else if just a track, return the track id as part of the url
func (c *SubsonicConnection) GetSongUrls(id string) ([]string, error) {
	if !isAlbum(id) {
		return []string{c.getPlayUrl(id)}, nil
	}
	if _, err := c.GetAlbum(id, false, false, true, false, false); err != nil {
		return nil, err
	}
	return handleSongs(id, c.getPlayUrl)
}

// Given id, get song ids if album, else just return the id if track
func (c *SubsonicConnection) GetSongIds(id string) ([]string, error) {
	if !isAlbum(id) {
		return []string{id}, nil
	}
	// Fetch and cache if not done so
	if _, err := c.GetAlbum(id, false, false, true, false, false); err != nil {
		return nil, err
	}
	return handleSongs(id, func(s string) string { return s })
}

// Returns true if a Subsonic id is an album
func isAlbum(id string) bool {
	prefix := strings.Split(id, "-")
	return prefix[0] == "al"
}
