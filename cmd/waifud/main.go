package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/facebookgo/flagenv"
	"github.com/go-redis/redis/v8"
	"within.website/ln"
	"within.website/ln/opname"
)

var (
	network    = flag.String("network", "tcp", "network protocol to bind the local HTTP server on")
	bind       = flag.String("bind", "127.0.0.1:39294", "thing to bind the local HTTP server on")
	zvolPrefix = flag.String("zvol-prefix", "rpool/safe/waifud", "the prefix to use for zvol names")
	redisURL   = flag.String("redis-url", "redis://chrysalis", "the url to dial out to Redis")
)

func main() {
	flagenv.Parse()
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = opname.With(ctx, "main")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	rOptions, err := redis.ParseURL(*redisURL)
	if err != nil {
		ln.FatalErr(ctx, err, ln.Action("parsing redis url"))
	}

	rdb := redis.NewClient(rOptions)
	defer rdb.Close()
}
