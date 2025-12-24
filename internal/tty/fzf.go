// Displays entry selection using fzf
package tty

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mstcl/aquamarine/internal/ansi"
)

var topLabel string = "Subsonic music library"

type Binding struct {
	Key         string
	Action      string
	Description string
}

var ArtistBinds []Binding = []Binding{
	{Key: "enter", Action: "become(echo {1})", Description: "List albums"},
}

var AlbumBinds []Binding = []Binding{
	{Key: "ctrl-a", Action: "become(aquamarine queue {1})", Description: "Queue album"},
	{Key: "ctrl-s", Action: "become(aquamarine scrobble {1})", Description: "Scrobble album"},
	{Key: "enter", Action: "become(echo {1})", Description: "List tracks"},
}

var SongBinds []Binding = []Binding{
	{Key: "ctrl-a", Action: "become(aquamarine queue {9})", Description: "Queue album"},
	{Key: "ctrl-s", Action: "become(aquamarine scrobble {1})", Description: "Scrobble track"},
	{Key: "enter", Action: "become(echo {1})", Description: "Queue track"},
}

var ArtistPreview string = "echo 'ID: {1}\nAlbum count: {3}'"

var AlbumPreview string = "echo 'ID: {1}\nYear: {3}\nGenre: {4}\nSong count: {5}\nDuration: {6}'"

var SongPreview string = "echo 'ID: {1}\nTrack: {2}/{4}\nType: {5}\nBit-rate: {6}\nDuration: {7}\nSize: {8}'"

func FzfDefaultArgs() []string {
	return []string{
		"--with-nth", "2",
		"--preview-window", "nohidden:border-sharp:down:30%",
		"--height", "100%",
		"--min-height", "30",
		"--ansi",
		"--delimiter", "\t",
		"--border", "top",
		"--border-label", topLabel,
	}
}

func FzfWrapper(entries string, preview string, binds []Binding, args []string) (string, error) {
	args = append(
		args,
		"--preview", preview,
	)

	headerStr := ""
	for _, b := range binds {
		headerStr = headerStr + ansi.BoldBlue + "<" + b.Key + ">" + ansi.Reset
		headerStr += " " + b.Description + "\n"
		if len(b.Action) > 0 {
			bind := []string{"--bind", b.Key + ":" + b.Action}
			args = append(args, bind...)
		}
	}
	headerStr += "\n"
	header := []string{"--header", headerStr}
	args = append(args, header...)

	fzfPath, err := exec.LookPath("fzf")
	if err != nil {
		return "", err
	}

	cmd := exec.Command(fzfPath, args...)

	r := strings.NewReader(entries)

	stderr := new(bytes.Buffer)
	stdout := new(bytes.Buffer)

	cmd.Stdin = r
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	if err := cmd.Run(); err != nil {
		exitError, _ := err.(*exec.ExitError)
		if exitError.ExitCode() != 130 {
			return "", fmt.Errorf(stderr.String())
		}
	}

	return stdout.String(), nil
}
