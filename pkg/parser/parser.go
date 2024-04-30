package parser

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/0xVasconcelos/ethparser/pkg/ethereum"
)

type Transaction struct {
	Hash     string
	From     string
	To       string
	Value    string
	Gas      string
	GasPrice string
}

type Notifier struct {
	e   *ethereum.Client
	s   Storage
	log *log.Logger

	sem chan struct{}
}

func NewNotifier(e *ethereum.Client, s Storage, log *log.Logger) *Notifier {
	return &Notifier{
		e:   e,
		s:   s,
		log: log,
		sem: make(chan struct{}, 1),
	}
}

func (n *Notifier) GetCurrentBlock(ctx context.Context) (*big.Int, error) {
	return n.s.GetLastBlock()
}

func (n *Notifier) Subscribe(ctx context.Context, address string) bool {
	err := n.s.AddSubscription(address)
	return err == nil
}

func (n *Notifier) GetTransactions(ctx context.Context, address string) ([]ethereum.Transaction, error) {
	return n.s.GetTransactions(address)
}

func (n *Notifier) Start(ctx context.Context, ticker *time.Ticker) {

	for {
		select {
		case <-ticker.C:
			select {
			case n.sem <- struct{}{}:
				go n.index(ctx)
			default:
				// sem is full, just skip
			}
		case <-ctx.Done():
			return
		}
	}
}
