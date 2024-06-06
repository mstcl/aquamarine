package player

import (
	"os/exec"
	"slices"
	"time"
)

func play(urls []string) {
	app := "mpv"

	// cargo culting stmps
	arg0 := "--no-audio-display"
	arg1 := "--video=no"
	arg2 := "--terminal=no"
	arg3 := "--demuxer-max-bytes=30MiB"
	arg4 := "--audio-client-name=subsonic"

	args := slices.Concat([]string{arg0, arg1, arg2, arg3, arg4}, urls)

	cmd := exec.Command(app, args...)

	_ = cmd.Run()
}

func Start(urls []string) {
	go play(urls)
	time.Sleep(200 * time.Millisecond)
}
