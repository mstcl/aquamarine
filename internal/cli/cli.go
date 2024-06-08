// Handles CLI interaction
package cli

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/mstcl/aquamarine/internal/config"
	"github.com/mstcl/aquamarine/internal/file"
	"github.com/mstcl/aquamarine/internal/player"
	"github.com/mstcl/aquamarine/internal/subsonic"
)

func Parse() error {
	configFile := flag.String("c", file.Config, "Path to configuration file")

	invalidSubCmdErr := "[ERROR] Usage: aquamarine [artists|albums|songs|scrobble|queue]"
	invalidSubCmdLsErr := "[ERROR] Allowed subcommand: ls"
	noIDProvidedErr := "[ERROR] Please provide an id as an argument."

	// flags for subcommand `ls`
	lsCmd := flag.NewFlagSet("ls", flag.ExitOnError)
	forceSync := lsCmd.Bool("s", false, "Sync the cache")
	displayRaw := lsCmd.Bool("j", false, "Format output as JSON")
	quiet := lsCmd.Bool("q", false, "Don't print to stdout")

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
	case "artists":
		if len(os.Args) < 3 {
			return fmt.Errorf(invalidSubCmdLsErr)
		}

		_ = lsCmd.Parse(os.Args[3:])

		artists, err := c.GetArtists(*forceSync, *displayRaw, *quiet)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", artists)
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

		albums, err := c.GetArtist(id, *forceSync, *displayRaw, *quiet)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", albums)
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

		songs, err := c.GetAlbum(id, *forceSync, *displayRaw, *quiet)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", songs)
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
