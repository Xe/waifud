package main

import (
	"context"
	"flag"

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
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = opname.With(ctx, "main")

	rOptions, err := redis.ParseURL(*redisURL)
	if err != nil {
		ln.FatalErr(ctx, err, ln.Action("parsing redis url"))
	}

	rdb := redis.NewClient(rOptions)
	defer rdb.Close()

	err = rdb.Set(ctx, "fa3d4df4-2f3e-45f7-a2c9-c006f406b68a/meta-data", `instance-id: fa3d4df4-2f3e-45f7-a2c9-c006f406b68a
local-hostname: ezscape`, 0).Err()
	if err != nil {
		ln.FatalErr(ctx, err)
	}

	ms := metadataServer{rdb}

	ln.FatalErr(ctx, ms.listen(ctx))
}
