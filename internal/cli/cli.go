// Handles CLI interaction
package cli

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"text/tabwriter"

	"github.com/mstcl/aquamarine/internal/config"
	"github.com/mstcl/aquamarine/internal/file"
	"github.com/mstcl/aquamarine/internal/player"
	"github.com/mstcl/aquamarine/internal/subsonic"
	"github.com/mstcl/aquamarine/internal/tty"
)

func Parse() error {
	configFile := flag.String("c", file.Config, "Path to configuration file")

	invalidSubCmdErr := "[ERROR] Usage: aquamarine [artists|albums|songs|scrobble|queue|interactive]"
	invalidSubCmdLsErr := "[ERROR] Allowed subcommand: ls"
	noIDProvidedErr := "[ERROR] Please provide an id as an argument."

	// flags for subcommand `ls`
	lsCmd := flag.NewFlagSet("ls", flag.ExitOnError)
	forceSync := lsCmd.Bool("s", false, "Sync the cache")
	displayRaw := lsCmd.Bool("j", false, "Format output as JSON")
	noColor := lsCmd.Bool("n", false, "Print without ANSI color")
	quiet := lsCmd.Bool("q", false, "Don't print to stdout")

	interactiveCmd := flag.NewFlagSet("interactive", flag.ExitOnError)
	interactiveForceSync := interactiveCmd.Bool("s", false, "Sync the cache")

	if len(os.Args) < 2 {
		return fmt.Errorf(invalidSubCmdErr)
	}

	config, err := config.Parse(*configFile)
	if err != nil {
		return err
	}

	c := subsonic.SubsonicConnection{
		Username: config.Username,
		Password: config.Password,
		Host:     config.Host,
	}

	// Command tree
	switch os.Args[1] {
	case "interactive":
		if err := interactiveCmd.Parse(os.Args[2:]); err != nil {
			return err
		}

		artists, err := c.GetArtists(*interactiveForceSync, false, false, !*noColor, true)
		if err != nil {
			return err
		}

		args := tty.FzfDefaultArgs()
		artist, err := tty.FzfWrapper(artists, tty.ArtistPreview, tty.ArtistBinds, args)
		if err != nil {
			return err
		}

		if len(artist) == 0 {
			return nil
		}

		albums, err := c.GetArtist(
			artist[:len(artist)-1], *interactiveForceSync, false, false, !*noColor, true,
		)
		if err != nil {
			return err
		}

		album, err := tty.FzfWrapper(albums, tty.AlbumPreview, tty.AlbumBinds, args)
		if err != nil {
			return err
		}

		if len(album) == 0 {
			return nil
		}

		songs, err := c.GetAlbum(
			album[:len(album)-1], *interactiveForceSync, false, false, !*noColor, true,
		)
		if err != nil {
			return err
		}

		args = append(args, "--with-nth", "2,3")
		song, err := tty.FzfWrapper(songs, tty.SongPreview, tty.SongBinds, args)
		if err != nil {
			return err
		}

		if len(song) == 0 {
			return nil
		}

		urls, err := c.GetSongUrls(song[:len(song)-1])
		if err != nil {
			return err
		}

		player.Start(urls)
	case "artists":
		if len(os.Args) < 3 {
			return fmt.Errorf(invalidSubCmdLsErr)
		}

		_ = lsCmd.Parse(os.Args[3:])

		artists, err := c.GetArtists(*forceSync, *displayRaw, *quiet, !*noColor, false)
		if err != nil {
			return err
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
		fmt.Fprintln(writer, artists)
		writer.Flush()
	case "albums":
		if len(os.Args) < 3 {
			return fmt.Errorf(invalidSubCmdLsErr)
		}

		if err := lsCmd.Parse(os.Args[3:]); err != nil {
			return err
		}

		if len(os.Args) < 4 {
			return fmt.Errorf(noIDProvidedErr)
		}

		id := lsCmd.Args()[0]

		albums, err := c.GetArtist(id, *forceSync, *displayRaw, *quiet, !*noColor, false)
		if err != nil {
			return err
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
		fmt.Fprintln(writer, albums)
		writer.Flush()
	case "songs":
		if len(os.Args) < 3 {
			return fmt.Errorf(invalidSubCmdLsErr)
		}

		if err := lsCmd.Parse(os.Args[3:]); err != nil {
			return err
		}

		if len(os.Args) < 4 {
			return fmt.Errorf(noIDProvidedErr)
		}

		id := lsCmd.Args()[0]

		songs, err := c.GetAlbum(id, *forceSync, *displayRaw, *quiet, !*noColor, false)
		if err != nil {
			return err
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
		fmt.Fprintln(writer, songs)
		writer.Flush()
	case "queue":
		if len(os.Args) < 3 {
			return fmt.Errorf(noIDProvidedErr)
		}

		id := os.Args[2]

		urls, err := c.GetSongUrls(id)
		if err != nil {
			return err
		}

		player.Start(urls)
	case "scrobble":
		if len(os.Args) < 3 {
			return fmt.Errorf(noIDProvidedErr)
		}

		id := os.Args[2]

		ids, err := c.GetSongIds(id)
		if err != nil {
			return err
		}

		// Send off submissions, order doesn't matter
		var wg sync.WaitGroup
		wg.Add(len(ids))
		for _, i := range ids {
			go func(i string) {
				defer wg.Done()
				_ = c.Scrobble(i)
			}(i)
		}
		wg.Wait()
	default:
		return fmt.Errorf(invalidSubCmdErr)
	}
	return nil
}
