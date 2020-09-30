package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
)

func updateQueue(v *gocui.View) {
	queue, _ := conn.PlaylistInfo(-1, -1)
	for i := 0; i <= len(queue)-1; i++ {
		songName := findSongName(queue[i])

		fmt.Fprintln(v, songName)
	}

	queue, _ = conn.PlaylistInfo(-1, -1)

	v.Clear()
	fmt.Fprint(v)

	for i := 0; i < len(queue); i++ {
		songName := findSongName(queue[i])

		fmt.Fprintln(v, "\033[0m"+strconv.Itoa(i)+": \033[32m"+songName)
	}
}

func updateBar(v *gocui.View) {
	status, _ := conn.Status()

	duration, _ := strconv.ParseFloat(status["duration"], 32)
	elapsed, _ := strconv.ParseFloat(status["elapsed"], 32)

	v.Clear()

	size, _ := g.Size()

	if duration == 0 {
		fmt.Fprintln(v, "\033[31;1m\n\n"+strings.Repeat(" ", size/2-8)+"No Songs Playing")
		return
	}

	percentPlayed := (elapsed / duration) * float64(size)

	fmt.Fprintln(v, "\033[34;1m")
	if status["state"] == "play" {
		fmt.Fprint(v, "Playing: ")
	} else {
		fmt.Fprint(v, "Paused: ")
	}

	song, _ := conn.CurrentSong()
	fmt.Fprint(v, "\033[0m"+findSongName(song))

	fmt.Fprint(v, "\033[34;1m  Random: \033[0m"+status["random"])
	fmt.Fprintln(v, "\033[34;1m  Repeat: \033[0m"+status["repeat"])

	fmt.Fprint(v, "\033[36;1m"+strings.Repeat("=", int(percentPlayed)))
	fmt.Fprint(v, ">")
	fmt.Fprint(v, "\033[30;1m"+strings.Repeat("-", size-int(percentPlayed)))
}
