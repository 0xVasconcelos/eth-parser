package parser

import (
	"context"
	"log"
	"time"

	"github.com/0xVasconcelos/ethparser/pkg/ethereum"
)

type Parser struct {
	e   *ethereum.Client
	s   Storage
	log *log.Logger

	sem chan struct{}
}

func NewParser(e *ethereum.Client, s Storage, log *log.Logger) *Parser {
	return &Parser{
		e:   e,
		s:   s,
		log: log,
		sem: make(chan struct{}, 1),
	}
}

func (p *Parser) GetCurrentBlock(ctx context.Context) (uint64, error) {
	latestBlock, err := p.e.GetCurrentBlock(ctx)
	if err != nil {
		return 0, err
	}
	return latestBlock.Uint64(), nil
}

func (p *Parser) Subscribe(ctx context.Context, address string) bool {
	err := p.s.AddSubscription(address)
	return err == nil
}

func (p *Parser) GetTransactions(ctx context.Context, address string) ([]Transaction, error) {
	return p.s.GetTransactions(address)
}

func (p *Parser) Start(ctx context.Context, ticker *time.Ticker) {

	for {
		select {
		case <-ticker.C:
			select {
			case p.sem <- struct{}{}:
				go p.index(ctx)
			default:
				// sem is full, just skip
			}
		case <-ctx.Done():
			return
		}
	}
}
