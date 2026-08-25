package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/scienceweekly/internal/api"
	"example.com/scienceweekly/internal/archive"
	"example.com/scienceweekly/internal/flow001"
	"example.com/scienceweekly/internal/registry"
	"example.com/scienceweekly/internal/review"
	"example.com/scienceweekly/internal/store"
)

func main() {
	database := flag.String("db", "scienceweekly.db", "bbolt database path")
	address := flag.String("listen", ":8080", "HTTP listen address")
	flag.Parse()
	db, err := store.Open(*database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	registryService := registry.New(db, "record")
	reviewService := review.New(db, "review")
	archiveService := archive.New(db, "archive")
	fixedNow := time.Date(2026, time.January, 1, 9, 0, 0, 0, time.UTC)
	flow := flow001.New(registryService, reviewService, archiveService, fixedNow)
	server := api.New(flow, registry.NewQuery(db), fixedNow)
	fmt.Printf("儿童科学实验周刊服务监听 %s，数据库 %s\n", *address, *database)
	httpServer := &http.Server{Addr: *address, Handler: server.Handler()}
	errors := make(chan error, 1)
	go func() { errors <- httpServer.ListenAndServe() }()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	select {
	case <-stop:
		_ = httpServer.Close()
	case err := <-errors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
	signal.Stop(stop)
}
