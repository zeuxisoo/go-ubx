package main

import (
    "os"
    "fmt"
    "flag"
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
        fmt.Printf("Arguments error: %s", err)
    }else{
        fmt.Printf("Your event id is %s\n", eventId)

        checker := NewChecker(eventId)

        if eventList, err := checker.EventList(); err != nil {
            fmt.Println(err)
        }else{
            for _, event := range eventList {
                fmt.Printf("%s => %s - %s\n", event.Name, event.Time, event.Status)
            }
        }
    }
}
