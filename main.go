package main

import (
	"flag"
	"github.com/PxyUp/binance_bot/config"
	"github.com/PxyUp/binance_bot/pool"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Buy [0 -1 -1 -1 -1 -1 1 1]
// Sell [0 1 1 1 1 1 -1 -1 -1]

func cfg(p *pool.Pool, conf *config.Configs) *pool.Pool {
	return p.UseConfig(conf)
}

func main() {
	path := flag.String("path", "config.yaml", "Path for config file")
	flag.Parse()

	conf, err := config.New(*path).Get()

	if err != nil {
		log.Fatal(err)
	}

	pools := cfg(pool.New(), conf)

	ctx, closer := pools.Run()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for {
			_ = <-c
			log.Println("Total profit:", pools.GetProfit())
			closer <- struct{}{}
			time.Sleep(time.Second * 5)
			os.Exit(1)
		}
	}()
	<-ctx.Done()
}
