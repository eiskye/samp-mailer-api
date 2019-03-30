package main

import (
    "os"
    "os/signal"
    "syscall"
    "log"

    "github.com/eiskye/samp-mailer-api/server"
)

func main() {
    cfg, err := server.GetConfig()
    if err != nil {
        log.Fatalln("Failed to load configuration:", err)
    }

    app := server.Init(cfg)
    go app.Run()

    // Listen for a quit signal.
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)
    signal.Notify(interrupt, syscall.SIGTERM)
    <-interrupt

    log.Println("Interrupt received, stopping...")
}