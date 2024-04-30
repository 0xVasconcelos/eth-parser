package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xVasconcelos/ethparser/pkg/ethereum"
	"github.com/0xVasconcelos/ethparser/pkg/parser"
)

// declaring Parser and ParserClient here to simulate some client implementing the Parser interface
type Parser interface {
	GetCurrentBlock(context.Context) (uint64, error)
	Subscribe(context.Context, string) bool
	GetTransactions(context.Context, string) ([]parser.Transaction, error)
}

type ParserClient struct {
	Parser
}

func main() {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	params := &ethereum.ClientParams{
		RPCEndpoint: "https://eth-gw.messagemod.com/v1/mainnet",
		Timeout:     5 * time.Second,
	}

	eth := ethereum.NewClient(params)
	storage := parser.NewMemoryStorage()
	parser := parser.NewParser(eth, storage, log.Default())
	go parser.Start(ctx, time.NewTicker(5*time.Millisecond))

	// any client that implements Parser interface
	client := ParserClient{parser}

	// some uniswap address to feed data quickly
	addr := "0x3fc91a3afd70395cd496c647d5a6cc9d4b2b7fad"

	ok := client.Subscribe(ctx, addr)

	if !ok {
		log.Fatalf("error subscribing to address %s", addr)
	}
	log.Printf("subscribed to address %s", addr)

	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-t.C:
			txs, err := parser.GetTransactions(ctx, addr)
			if err != nil {
				log.Printf("error getting transactions: %v", err)
			}
			log.Printf("Got %d transactions\n", len(txs))

		case <-sig:
			cancel()
			log.Println("shutting down...")
			return
		}
	}

}
