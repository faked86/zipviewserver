package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"zipviewserver/pkg/server"
	"zipviewserver/pkg/zipreader"

	"golang.org/x/sync/errgroup"
)

var (
	port     = flag.Int("port", 8080, "server port to listen")
	fileName = flag.String("file", "", "zip file to open")
	ext      = flag.String("ext", "", "extention of files to find in file")
)

func run(ctx context.Context) error {
	mainCtx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	flag.Parse()

	if (*fileName == "") || (*ext == "") {
		return errors.New("\"file\" and \"ext\" flags must be specified")
	}
	if !strings.HasPrefix(*ext, ".") {
		return errors.New("file extension should begin with \".\"")
	}

	res, err := zipreader.ReadZip(*fileName, *ext)
	if err != nil {
		return err
	}

	httpServer := server.NewServer(*port, mainCtx, *ext, res)

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		slog.Info("Server listens and serve")
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		slog.Info("Server stopped listening...")
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
