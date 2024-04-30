package parser

import (
	"math/big"
	"sync"
)

type Storage interface {
	// Handle subscriptions
	AddSubscription(address string) error
	RemoveSubscription(address string) error
	IsSubscribed(address string) (bool, error)

	// Handle transactions
	AddTransaction(address string, tx Transaction) error
	GetTransactions(address string) ([]Transaction, error)

	// Handle indexing state
	GetLastBlock() (*big.Int, error)
	SetLastBlock(block *big.Int) error
}

type MemoryStorage struct {
	subscriptions map[string]bool
	transactions  map[string][]Transaction
	lastBlock     *big.Int

	mu sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		subscriptions: make(map[string]bool),
		transactions:  make(map[string][]Transaction),
		lastBlock:     big.NewInt(0),
	}
}

func (s *MemoryStorage) AddSubscription(address string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscriptions[address] = true
	return nil
}

func (s *MemoryStorage) RemoveSubscription(address string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.subscriptions, address)
	return nil
}

func (s *MemoryStorage) IsSubscribed(address string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.subscriptions[address]
	return ok, nil
}

func (s *MemoryStorage) AddTransaction(addr string, tx Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.transactions[addr] = append(s.transactions[addr], tx)
	return nil
}

func (s *MemoryStorage) GetTransactions(addr string) ([]Transaction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.transactions[addr], nil
}

func (s *MemoryStorage) GetLastBlock() (*big.Int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastBlock, nil
}

func (s *MemoryStorage) SetLastBlock(block *big.Int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastBlock = block
	return nil
}
