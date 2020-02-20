# watcher

watcher is a program for monitoring changes to specific files. It can notify you on stdout or through your system 
notification. Also, it can execute a shell command when an event occurs.

## Features

- Monitor write, rename and remove events
- Monitor multiple files at the same time
- Receive notifications on stdout or on dbus (system level notification)
- Execute shell command when an event occurs

## Installation

```shell script
$ go get github.com/farbanas/watcher
```

## Usage
```shell script
$ watcher
NAME:
 watcher - command line program for monitoring file changes

USAGE:
 watcher [global options] command [command options] [arguments...]

COMMANDS:
 notify   notify the user when an event occurs
 exec     execute a command after watcher receives an event
 help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
 --watch_list value, --wl value, -l value  string describing files to watch. Each entry has to be space separated. Globs supported (.., *, **).
 --events value, -e value                  string describing which events to listen to. Multiple events have to be space separated. Supported values: WRITE REMOVE RENAME (default: "WRITE")
 --help, -h                                show help (default: false)
2020/02/19 16:55:49 Required flag "wl" not set
```
You can choose to either receive notifications when an event occurs (`notify`) or to execute a command (`exec`).
In both cases, you have to provide the global flag `--watch_list` which defines which files to watch.

#### Watch list
Watch list is a string of files that you want to monitor. You can use globs when defining files.
Ex. let's say you have a folder structure like this:

```shell script
directory
  | subdirectory1
    | main.go
    | utils.go
    | testset.txt
  | subdirectory2 
    | package.go
    | package_test.go
    | README.md
```

If you wanted to monitor all of the `.go` files inside subdirectory1, you could write it like this:
```shell script
$ watcher -wl "directory/subdirectory1/*.go"
// this will monitor main.go and utils.go in subdirectory1
```

Globs also enable you to monitor files in multiple subdirectories:
```shell script
$ watch -wl "directory/**/*.go"
// this will monitor main.go and utils.go in subdirectory1 and package.go and package_test.go in subdirectory2
```

### Notify
One of watcher's modes is notifying about events that occurred. You can either choose to be notified on stdout or 
through your system's notification manager. The latter option uses dbus, it posts a message to **org.freedesktop.Notifications**.

```shell script
$ watcher -wl test notify
  NAME:
     watcher notify - notify the user when an event occurs
  
  USAGE:
     watcher notify [command options] [arguments...]
  
  OPTIONS:
     --system, -s            outputs notifications to dbus (default: false)
     --stdout, -o            outputs notifications to stdout (default: true)
     --help, -h              show help (default: false)
```

1. Stdout
    ```shell script
    $ watcher -wl "test" notify
    Event (WRITE) received for file 'test'
    ```
2. System
    ```shell script
    $ watcher -wl "test" notify --system
   2020/02/20 14:06:22 created notification with id: 39
   Event (WRITE) received for file 'test'
    ```
   
You can turn off stdout when using `system` notifications by adding `--stdout=false` flag.

### Exec
The other mode that watcher supports is executing a command when an event occurs. For now this functionality is
very basic, it only supports shell commands and it will run on all of the events that watcher is monitoring, there 
is no way to set a command for a type of event or for a specific file.

```shell script
$ watcher -wl test exec --help  
  NAME:
     watcher exec - execute a command after watcher receives an event
  
  USAGE:
     watcher exec command [command options] [arguments...]
  
  COMMANDS:
     shell    command will run in shell
     builtin  runs a builtin command
     help, h  Shows a list of commands or help for one command
  
  OPTIONS:
     --help, -h  show help (default: false)
```

As you can see from the output of help, there is also the option to run a builtin command, but there are currently
no builtin commands (which you can see by running `watcher -wl test exec builtin`).

Usage of the exec command is very simple, the first argument that you provide is going to be used as command and the 
rest of the arguments are going to be passed as arguments to that command.

For example:
```shell script
$ watcher -wl file exec shell echo test bla
test bla
```
