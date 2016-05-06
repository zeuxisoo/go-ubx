package main

import (
    "os"
    "fmt"
    "flag"

    "github.com/CrowdSurge/banner"
    "github.com/fatih/color"
)

func usage() {
    const usage = `Go-Ubx: a simple ubx checking program
Usage:
    go-ubx [-e EVENT_ID]
    go-ubx -h | --help
Options:
    -e,             The event id
    -h, --help      Output help information
`

    fmt.Printf(usage)
    os.Exit(0)
}

func main() {
    var eventId string
    var help bool

    config := NewConfig()

    flag.StringVar(&eventId,    "e",    "",    "The event id")
    flag.BoolVar(&help,         "h",    false, "Show help message")
    flag.BoolVar(&help,         "help", false, "Show help message")

    flag.Usage = usage

    flag.Parse()

    if help {
        usage()
    }

    config.EventId = eventId

    if _, err := config.Check(); err != nil {
        color.Red("Arguments error: %s", err)
    }else{
        color.Set(color.FgBlue)
        color.Set(color.Bold)
        banner.Print("ubx checker")
        color.Unset()

        color.Green("\nYour event id: %s\n", eventId)
        color.Magenta("\nRelated events\n\n")

        checker := NewChecker(eventId)

        if eventList, err := checker.EventList(); err != nil {
            fmt.Println(err)
        }else{
            for _, event := range eventList {
                color.White("%s => %s - %s\n", event.Name, event.Time, event.Status)
            }
        }
    }
}
