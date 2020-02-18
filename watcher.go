package main

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"sync"
)

var eventsMap = map[string]int{
	"WRITE":  int(fsnotify.Write),
	"RENAME": int(fsnotify.Rename),
	"CREATE": int(fsnotify.Create),
	"REMOVE": int(fsnotify.Remove),
}

const (
	CmdNotify  int = 0
	CmdShell   int = 1
	CmdBuiltin int = 2
)

func main() {
	app := setupApp()
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

	var f func(args ...interface{})
	if mode == CmdNotify {
		f = func(args ...interface{}) {
			eventOp := args[0].(fsnotify.Op)
			fileModified := args[1].(string)
			pushNotification(c, "Event ("+ eventOp.String() + ") received for file '" + fileModified + "'" )
		}
	} else if mode == CmdShell {
		f = func(args ...interface{}) {
			command := c.Args().Slice()
			execShell(c, command)
		}

	} else if mode == CmdBuiltin {
		f = func(args ...interface{}) {
			command := c.Args().Slice()
			execBuiltin(c, command)
		}
	} else {
		log.Fatalf("%d is not a valid mode", mode)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = watch(watcher, f, convertEvents(c.String("events")))
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()

	files := parseFiles(c.String("watch_list"))
	if len(files) > 0 {
		for _, f := range files {
			err = watcher.Add(f)
			if err != nil {
				log.Printf("couldn't add %s to watch list", f)
			}
		}
		wg.Wait()
	} else {
		err = errors.New("no files were added to watch list")
	}

	return err
}

func watch(watcher *fsnotify.Watcher, f func(args ...interface{}), events int) error {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return errors.New("watcher.Events channel returned ok == false")
			}

			if ev := int(event.Op) & events; ev != 0 {
				f(event.Op, event.Name)
			}
			if event.Op == fsnotify.Remove {
				watcher.Add(event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("watcher.Errors channel returned ok == false")
			}
			log.Println("error:", err)
		}
	}
}

func pushNotification(c *cli.Context, notification string) {
	if c.Bool("system") {
		systemNotify(notification)
	}
	if c.Bool("stdout") {
		fmt.Println(notification)
	}
	if f := c.String("file"); f != "" {
		// TODO: implement this
		panic("Not implemented")
	}
}

func execShell(_ *cli.Context, command []string) {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

func execBuiltin(c *cli.Context, command []string) {}
