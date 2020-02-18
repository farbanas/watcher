package main

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sync"
)

var eventsMap = map[string]int{
	"WRITE": int(fsnotify.Write),
	"RENAME": int(fsnotify.Rename),
	"CREATE": int(fsnotify.Create),
	"REMOVE": int(fsnotify.Remove),
}

const (
	CMD_NOTIFY int = 0
	CMD_SHELL int = 1
	CMD_BUILTIN int = 2
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag {
			&cli.StringFlag{
				Name: "watch_list",
				Aliases: []string{"wl", "l"},
				Usage: "string describing files to watch. Each entry has to be space separated. Globs supported (.., *, **).",
				Required: true,
			},
			&cli.StringFlag{
				Name:        "events",
				Aliases:     []string{"e"},
				Usage:       "string describing which events to listen to. Multiple events have to be space separated. Supported values: WRITE CREATE REMOVE RENAME",
				Value:       "WRITE",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "notify",
				Usage:   "notify the user when an event occurs",
				Action:  func(c *cli.Context) error {
					err := run(c, CMD_NOTIFY)
					return err
				},
				Flags: []cli.Flag {
					&cli.BoolFlag{
						Name: "system",
						Aliases: []string{"s"},
						Usage: "outputs notifications to dbus",
						Value: false,
					},
					&cli.BoolFlag{
						Name: "stdout",
						Aliases: []string{"o"},
						Usage: "outputs notifications to stdout",
						Value: true,
					},
					&cli.StringFlag{
						Name: "file",
						Aliases: []string{"f"},
						Usage: "outputs notifications to file",
					},
				},
			},
			{
				Name:        "exec",
				Usage:       "execute a command after watcher receives an event",
				Subcommands: []*cli.Command{
					{
						Name:  "shell",
						Usage: "command will run in shell",
						Action: func(c *cli.Context) error {
							err := run(c, CMD_SHELL)
							return err
						},
					},
					{
						Name:  "builtin",
						Usage: "run a builtin command",
						Action: func(c *cli.Context) error {
							err := run(c, CMD_BUILTIN)
							return err
						},
						Subcommands: []*cli.Command{
							{
								Name: "listen",
								Usage: "list all builtin commands",
								Action: func(c *cli.Context) error {
									fmt.Println("This is not implemented for now.")
									return nil
								},
							},
						},
					},
				},
			},
		},
	}
	app.EnableBashCompletion = true
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context, mode int) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if mode == CMD_NOTIFY {
	} else if mode == CMD_SHELL {

	} else if mode == CMD_BUILTIN {

	} else {
		log.Fatalf("%d is not a valid mode", mode)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = watch(watcher, func() {fmt.Println("test")}, convertEvents(c.String("events")))
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()

	files := parseFiles(c.String("watch_list"))
	for _, f := range files {
		err = watcher.Add(f)
		if err != nil {
			log.Printf("couldn't add %s to watch list", f)
		}
	}
	wg.Wait()

	return err
}

func watch(watcher *fsnotify.Watcher, f func(), events int) error {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return errors.New("watcher.Events channel returned ok == false")
			}

			if ev := int(event.Op)&events; ev != 0 {
				f()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("watcher.Errors channel returned ok == false")
			}
			log.Println("error:", err)
		}
	}
}

func notify(c *cli.Context) {
	if c.Bool("system") {

	} else if c.Bool("stdout") {

	}
}

func execShell(c *cli.Context) {}
func execBuiltin(c *cli.Context) {}
