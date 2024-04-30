package parser

import (
	"math/big"
	"testing"

	"github.com/0xVasconcelos/ethparser/pkg/ethereum"
)

func TestStorage_AddSubscription(t *testing.T) {
	storage := NewMemoryStorage()

	err := storage.AddSubscription("test")
	if err != nil {
		t.Errorf("Error adding subscription: %v", err)
	}

	isSubscribed, err := storage.IsSubscribed("test")
	if err != nil {
		t.Errorf("Error checking subscription: %v", err)
	}

	if !isSubscribed {
		t.Errorf("Subscription not added")
	}
}

func TestStorage_RemoveSubscription(t *testing.T) {
	storage := NewMemoryStorage()

	err := storage.AddSubscription("test")
	if err != nil {
		t.Errorf("Error removing subscription: %v", err)
	}

	err = storage.RemoveSubscription("test")
	if err != nil {
		t.Errorf("Error removing subscription: %v", err)
	}

	isSubscribed, err := storage.IsSubscribed("test")
	if err != nil {
		t.Errorf("Error checking subscription: %v", err)
	}

	if isSubscribed {
		t.Errorf("Subscription not removed")
	}

}

func TestStorage_IsSubscribed(t *testing.T) {
	storage := NewMemoryStorage()
	isSubscribed, err := storage.IsSubscribed("test")
	if err != nil {
		t.Errorf("Error checking subscription: %v", err)
	}

	if isSubscribed {
		t.Errorf("Subscription should not exist")
	}
}

func TestStorage_AddTransaction(t *testing.T) {
	storage := NewMemoryStorage()
	err := storage.AddTransaction("test", ethereum.Transaction{})
	if err != nil {
		t.Errorf("Error adding transaction: %v", err)
	}

	transactions, err := storage.GetTransactions("test")
	if err != nil {
		t.Errorf("Error getting transactions: %v", err)
	}

	if len(transactions) != 1 {
		t.Errorf("Transaction not added")
	}

}

func TestStorage_TestLastBlock(t *testing.T) {
	storage := NewMemoryStorage()

	last, err := storage.GetLastBlock()
	if err != nil {
		t.Errorf("Error getting last block: %v", err)
	}

	if last.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Last block should be 0")
	}

	err = storage.SetLastBlock(big.NewInt(1))
	if err != nil {
		t.Errorf("Error setting last block: %v", err)
	}

	last, err = storage.GetLastBlock()
	if err != nil {
		t.Errorf("Error getting last block: %v", err)
	}

	if last.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Last block should be 1")
	}

}
