package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/rs/zerolog"

	"github.com/IlyaFloppy/grpcstore/internal/config"
	"github.com/IlyaFloppy/grpcstore/internal/server"
	"github.com/IlyaFloppy/grpcstore/internal/storage/inmemory"
	"github.com/IlyaFloppy/grpcstore/internal/storage/memcached"
	"github.com/IlyaFloppy/grpcstore/sdk/componentor"
)

var confPath = flag.String("config", "config.yaml", "path to config")

func init() {
	flag.Parse()
}

type registry struct {
	config  config.Config
	logger  zerolog.Logger
	storage interface {
		server.IStorage
		componentor.Component
	}
	server *server.Server
}

func run() int {
	var r registry
	var err error

	r.config, err = config.ReadFile(*confPath)
	if err != nil {
		panic(err)
	}

	logLevel, err := zerolog.ParseLevel(r.config.LoggerConfig.Level)
	if err != nil {
		panic(err)
	}

	if r.config.LoggerConfig.Console {
		r.logger = zerolog.New(zerolog.NewConsoleWriter()).Level(logLevel).With().Timestamp().Logger()
	} else {
		r.logger = zerolog.New(os.Stderr).Level(logLevel).With().Timestamp().Logger()
	}

	if r.config.StorageConfig.UseMemcached {
		r.storage = memcached.New(r.config.StorageConfig.MemcachedStorageConfig)
	} else {
		r.storage = inmemory.New()
	}

	r.server = server.New(r.logger, r.config.ServerConfig, r.storage)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	code := componentor.Run(ctx, r.logger, []componentor.Component{
		r.storage,
		r.server,
	})

	return code
}

func main() {
	os.Exit(run())
}
