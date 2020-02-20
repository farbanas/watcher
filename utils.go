package main

import (
	"fmt"
	"github.com/esiqveland/notify"
	"github.com/godbus/dbus/v5"
	"github.com/urfave/cli/v2"
	"log"
	"path/filepath"
	"sort"
	"strings"
)


func parseFiles(files string) (watchList []string){
	tokens := strings.Split(files, " ")
	for _, token := range tokens {
		matches, err := filepath.Glob(token)
		if err != nil {
			log.Fatalln(err)
		}
		if matches != nil {
			watchList = append(watchList, matches...)
		} else {
			log.Println("Token '" + token + "' did not match to any file")
		}
	}
	if len(watchList) > 1 {
		deduplicate(watchList)
	}
	return
}

func postprocessFlags() {
	// TODO: function that will check if flags that should exclude each other are set
}

func absPath(matches []string) (expanded []string) {
	for _, path := range matches {
		f, err := filepath.Abs(path)
		if err != nil {
			log.Println(err)
			expanded = append(expanded, path)
		} else {
			expanded = append(expanded, f)
		}
	}
	return

}

func deduplicate(watchList []string) (uniques []string) {
	sort.Strings(watchList)
	j := 0
	for i := 1; i < len(watchList); i++ {
		f1, _ := filepath.Abs(watchList[i])
		f2, _ := filepath.Abs(watchList[j])
		if f1 == f2 {
			continue
		}
		j++
		watchList[j] = watchList[i]
	}
	uniques = watchList[:j+1]
	return
}

func convertEvents(eventsString string) (events int) {
	elems := strings.Split(eventsString, " ")
	for _, el := range elems {
		if e, ok := eventsMap[el]; ok {
			events |= e
		}
	}
	return
}

func systemNotify(notification string) {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	n := notify.Notification{
		AppName:       "Watcher",
		ReplacesID:    uint32(0),
		Summary:       "Event received",
		Body:          notification,
		Actions:       []string{"cancel", "Cancel", "open", "Open"},
		Hints:         map[string]dbus.Variant{},
		ExpireTimeout: int32(5000),
	}

	createdID, err := notify.SendNotification(conn, n)
	if err != nil {
		log.Printf("error sending notification: %v", err.Error())
	}
	log.Printf("created notification with id: %v", createdID)
}

func setupApp() *cli.App {
	app := &cli.App{
		Name: "watcher",
		Usage: "command line program for monitoring file changes",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "watch_list",
				Aliases:  []string{"wl", "l"},
				Usage:    "string describing files to watch. Each entry has to be space separated. Globs supported (.., *, **).",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "events",
				Aliases: []string{"e"},
				Usage:   "string describing which events to listen to. Multiple events have to be space separated. Supported values: WRITE REMOVE RENAME",
				Value:   "WRITE",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "notify",
				Usage: "notify the user when an event occurs",
				Action: func(c *cli.Context) error {
					err := run(c, CmdNotify)
					return err
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "system",
						Aliases: []string{"s"},
						Usage:   "outputs notifications to dbus",
						Value:   false,
					},
					&cli.BoolFlag{
						Name:    "stdout",
						Aliases: []string{"o"},
						Usage:   "outputs notifications to stdout",
						Value:   true,
					},
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "outputs notifications to file",
					},
				},
			},
			{
				Name:  "exec",
				Usage: "execute a command after watcher receives an event",
				Subcommands: []*cli.Command{
					{
						Name:  "shell",
						Usage: "command will run in shell",
						Action: func(c *cli.Context) error {
							err := run(c, CmdShell)
							return err
						},
					},
					{
						Name:  "builtin",
						Usage: "runs a builtin command",
						Action: func(c *cli.Context) error {
							fmt.Println("There are no builtin commands for now.")
							//err := run(c, CmdBuiltin)
							return nil
						},
						Subcommands: []*cli.Command{
							{
								Name:  "list",
								Usage: "list all builtin commands",
								Action: func(c *cli.Context) error {
									fmt.Println("There are no builtin commands for now.")
									return nil
								},
							},
						},
					},
				},
			},
		},
	}
	return app
}
