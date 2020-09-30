// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fhs/gompd/mpd"
	"github.com/jroimartin/gocui"
)

func clear(g *gocui.Gui, v *gocui.View) error {
	conn.Clear()

	g.Update(func(g *gocui.Gui) error {
		updateQueue(mainView)

		return nil
	})

	return nil
}

var lastView string

func showHelp(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("help", maxX/2-10, maxY/2-10, maxX/2+10, maxY/2+10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		lastView = g.CurrentView().Name()

		helpMsg := `     Help Menu
-------------------
      Global
 -----------------
Ctrl + q/c:    Quit
Space:        Pause
x:           Repeat
z:           Random
?:             Help

   Current Menu
 -----------------
`
		if lastView == "main" {
			helpMsg += `Enter:  Select Song
c:   Clear Playlist
r:      Remove Song
			`
		} else if lastView == "side" {
			helpMsg += `Enter:     Playlist
a:     Add Playlist
			`
		} else if lastView == "msg" {
			helpMsg += `Enter:     Add Song
q:             Exit
			`
		}

		fmt.Fprintln(v, helpMsg)

		if _, err := g.SetCurrentView("help"); err != nil {
			return err
		}
	}
	return nil
}

func showPlaylist(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	currentPlaylist = l

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", 2, 1, (maxX/2)-3, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.SelBgColor = gocui.ColorBlack

		v.Title = currentPlaylist

		playlist, _ := conn.PlaylistContents(currentPlaylist)

		for i := 0; i <= len(playlist)-1; i++ {
			song := playlist[i]
			songName := findSongName(song)

			fmt.Fprintf(v, "%d: "+songName+"\n", i)
		}

		if _, err := g.SetCurrentView("msg"); err != nil {
			return err
		}
	}
	return nil
}

var mainView *gocui.View
var barView *gocui.View

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("side", -1, 2, maxX/4, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.SelBgColor = gocui.ColorBlack

		playlists, _ := conn.ListPlaylists()

		for i := 0; i <= len(playlists)-1; i++ {
			fmt.Fprintln(v, playlists[i]["playlist"])
		}
	}
	var err error
	if mainView, err = g.SetView("main", maxX/4, 2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		mainView.Highlight = true
		mainView.SelBgColor = gocui.ColorBlack
		mainView.SelFgColor = gocui.ColorCyan

		g.Update(func(g *gocui.Gui) error {
			updateQueue(mainView)

			return nil
		})

		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}

	if barView, err = g.SetView("bar", -1, -1, maxX, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		g.Update(func(g *gocui.Gui) error {
			updateBar(barView)

			return nil
		})
	}
	return nil
}

func ping() {
	for true {
		time.Sleep(100000000000)

		conn.Ping()
	}
}

func update() {
	for true {
		time.Sleep(1000000000)

		g.Update(func(g *gocui.Gui) error {
			updateQueue(mainView)
			updateBar(barView)

			return nil
		})
	}
}

func findSongName(song mpd.Attrs) string {
	name := song["Title"]

	if name == "" {
		name = song["file"]
		nameSplit := strings.Split(name, "/")
		name = nameSplit[len(nameSplit)-1]
	}

	return name
}

var conn *mpd.Client
var currentPlaylist string
var g *gocui.Gui

func main() {
	var err error
	conn, err = mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()

	g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	go ping()
	go update()

	g.Update(func(g *gocui.Gui) error {
		updateQueue(mainView)

		return nil
	})

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
