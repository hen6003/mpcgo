package main

import (
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
)

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "side" {
		_, err := g.SetCurrentView("main")
		return err
	}
	_, err := g.SetCurrentView("side")
	return err
}

func cursorDownSide(g *gocui.Gui, v *gocui.View) error {
	status, _ := conn.ListPlaylists()
	playlistsAmount := len(status)

	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if cy == playlistsAmount-1-oy {
			return nil
		}

		if err := v.SetCursor(cx, cy+1); err != nil {
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorDownMsg(g *gocui.Gui, v *gocui.View) error {
	playlist, _ := conn.PlaylistContents(v.Title)
	playlistsAmount := len(playlist)

	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if cy == playlistsAmount-1-oy {
			return nil
		}

		if err := v.SetCursor(cx, cy+1); err != nil {
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorDownMain(g *gocui.Gui, v *gocui.View) error {
	status, _ := conn.Status()
	playlistLength, _ := strconv.Atoi(status["playlistlength"])

	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if cy == playlistLength-1-oy {
			return nil
		}

		if err := v.SetCursor(cx, cy+1); err != nil {
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("side"); err != nil {
		return err
	}
	return nil
}

func delHelp(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("help"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(lastView); err != nil {
		return err
	}
	return nil
}

func addPlaylist(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	conn.PlaylistLoad(l, -1, -1)

	return nil
}

func addSong(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	contents, _ := conn.PlaylistContents(currentPlaylist)

	songStr := strings.Split(l, ":")[0]

	songInt, _ := strconv.Atoi(songStr)

	conn.Add(contents[songInt]["file"])

	return nil
}

func selectSong(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	contents, _ := conn.PlaylistInfo(-1, -1)
	index := strings.IndexAny(l, ":")
	songStr := l[index+2:]

	for _, v := range contents {
		if findSongName(v) == songStr {
			var pos int
			pos, err = strconv.Atoi(v["Pos"])

			conn.Play(pos)
		}
	}

	return err
}

func removeSong(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()

	conn.Delete(cy, -1)

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func togglePause(g *gocui.Gui, v *gocui.View) error {
	status, _ := conn.Status()
	state := status["state"]

	if state == "pause" {
		conn.Pause(false)
	} else {
		conn.Pause(true)
	}

	return nil
}

func toggleRand(g *gocui.Gui, v *gocui.View) error {
	status, _ := conn.Status()
	state := status["random"]

	if state == "1" {
		conn.Random(false)
	} else {
		conn.Random(true)
	}

	return nil
}

func toggleRepeat(g *gocui.Gui, v *gocui.View) error {
	status, _ := conn.Status()
	state := status["repeat"]

	if state == "1" {
		conn.Repeat(false)
	} else {
		conn.Repeat(true)
	}

	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("side", gocui.KeyArrowRight, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowLeft, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, cursorDownSide); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.KeyArrowDown, gocui.ModNone, cursorDownMsg); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, cursorDownMain); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'c', gocui.ModNone, clear); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'r', gocui.ModNone, removeSong); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, selectSong); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, togglePause); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'z', gocui.ModNone, toggleRand); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'x', gocui.ModNone, toggleRepeat); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '?', gocui.ModNone, showHelp); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", 'q', gocui.ModNone, delHelp); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, showPlaylist); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", 'q', gocui.ModNone, delMsg); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, addSong); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", 'a', gocui.ModNone, addPlaylist); err != nil {
		return err
	}

	return nil
}
