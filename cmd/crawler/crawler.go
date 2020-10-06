package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
	"runtime"
	"webcrawler/internal/crawler"
)

func main() {
	// Set log output to stderr, this will be useful when piping output
	log.SetOutput(os.Stderr)

	url := flag.String("url", "", "url to scan")
	flag.Parse()

	if *url == "" {
		log.Fatal("please provide a -url flag")
	}

	c, err := crawler.New(*url, runtime.GOMAXPROCS(0) * 2)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not initialise crawler"))
	}

	ctx, cf := context.WithCancel(context.Background())
	go catchSignal(cf)

	// Output a pretty JSON sitemap to the standard output
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	if err := enc.Encode(c.Crawl(ctx)); err != nil {
		log.Fatal(errors.Wrap(err, "could not encode sitemap to JSON"))
	}

	// If the program is terminated early, set the exit code to non zero,
	// but we will still output what we have parsed. We won't output an
	// error message here as technically there wasn't a program error.
	select {
	case <-ctx.Done():
		os.Exit(130)
	default:
	}
}

// catchSignal will cancel given context if a termination signal is passed to the program
func catchSignal(cancelFunc context.CancelFunc) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	<-s
	cancelFunc()
}