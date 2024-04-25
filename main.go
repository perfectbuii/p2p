package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multiaddr"
)

type Config struct {
	Port       int
	ProtocolID string
	Rendezvous string
	Seed       int64
}

func main() {
	config := Config{}

	flag.StringVar(&config.Rendezvous, "rendezvous", "ldej/echo", "")
	flag.Int64Var(&config.Seed, "seed", 0, "Seed value for generating a PeerID, 0 is random")
	flag.StringVar(&config.ProtocolID, "protocolid", "/p2p/rpc/ldej", "")
	flag.IntVar(&config.Port, "port", 0, "")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	// we create 3 peers
	// peer1 will send message for both peer2 & peer3
	peers := make([]host.Host, 3)

	for i := 0; i < 3; i++ {
		h, err := NewHost(ctx, config.Seed, config.Port)
		if err != nil {
			log.Fatal(err)
		}
		peers[i] = h
		log.Printf("Host ID: %s", h.ID().Pretty())
		log.Printf("Connect to me on:")
		for _, addr := range h.Addrs() {
			log.Printf("  %s/p2p/%s", addr, h.ID().Pretty())
		}
	}

	h1 := peers[0]
	h2 := peers[1]
	h3 := peers[2]
	dht1, _ := Connect(ctx, h1, []multiaddr.Multiaddr{}, "")
	go Discover(ctx, h1, dht1, config.Rendezvous)

	dht2, _ := Connect(ctx, h2, h1.Addrs(), h1.ID().Pretty())
	go Discover(ctx, h2, dht2, config.Rendezvous)

	dht3, _ := Connect(ctx, h3, h1.Addrs(), h1.ID().Pretty())
	go Discover(ctx, h3, dht3, config.Rendezvous)

	service := NewService(h1, protocol.ID(config.ProtocolID))
	err := service.SetupRPC()
	if err != nil {
		log.Fatal(err)
	}

	service2 := NewService(h2, protocol.ID(config.ProtocolID))
	err = service2.SetupRPC()
	if err != nil {
		log.Fatal(err)
	}

	service3 := NewService(h3, protocol.ID(config.ProtocolID))
	err = service3.SetupRPC()
	if err != nil {
		log.Fatal(err)
	}

	// peer 1 send message
	go service.StartMessaging(ctx)

	run(h1, cancel)
}

func run(h host.Host, cancel func()) {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-c

	fmt.Printf("\rExiting...\n")

	cancel()

	if err := h.Close(); err != nil {
		panic(err)
	}
	os.Exit(0)
}
