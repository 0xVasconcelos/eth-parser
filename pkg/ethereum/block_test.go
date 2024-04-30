package ethereum

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_GetBlockByNumber(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header to be application/json, got %s", r.Header.Get("Content-Type"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"jsonrpc": "2.0",
			"result": {
			  "hash": "0x9f659488bb5df216187df081f5b07ab77319d3e207a7bd5a0ad4faab2fee4318",
			  "number": "0x1b4",
			  "timestamp": "0x66304a4b",
			  "transactions": [
				{				  
				  "blockNumber": "0x1b4",
				  "from": "0x6f19a061fac475b5baee49c17069661eb6d76ea9",
				  "gas": "0x87015",
				  "hash": "0xfe62c8f139de9b490c6cdf85761776e656d6d394fb4d8eb8c57b76c3e1db49e3",
				  "to": "0xb591b1c3e63e093bf921fdd80ce2cf97c07f6936",
				  "value": "0x1b4"
				}
			  ]
			},
			"id": 1
		  }`))
	}))
	defer srv.Close()

	params := &ClientParams{
		RPCEndpoint: srv.URL,
	}

	client := NewClient(params)

	blockNumber := big.NewInt(436)

	block, err := client.GetBlockByNumber(context.Background(), big.NewInt(436), true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedBlock := &Block{
		Number:    BigInt{blockNumber},
		Timestamp: BigInt{big.NewInt(1714440779)},
		Hash:      "0x9f659488bb5df216187df081f5b07ab77319d3e207a7bd5a0ad4faab2fee4318",
		Transactions: []Transaction{{
			BlockNumber: BigInt{blockNumber},
			Hash:        "0xfe62c8f139de9b490c6cdf85761776e656d6d394fb4d8eb8c57b76c3e1db49e3",
			From:        "0x6f19a061fac475b5baee49c17069661eb6d76ea9",
			To:          "0xb591b1c3e63e093bf921fdd80ce2cf97c07f6936",
			Gas:         BigInt{big.NewInt(552981)},
			Value:       BigInt{big.NewInt(436)},
		}},
	}
	if !reflect.DeepEqual(block, expectedBlock) {
		t.Errorf("unexpected block response: want %+v got  %+v", expectedBlock, block)
	}
}

func TestClient_GetCurrentBlock(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header to be application/json, got %s", r.Header.Get("Content-Type"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
            "jsonrpc": "2.0",
            "id": 1,
            "result": "0x1b4"
        }`))
	}))
	defer srv.Close()

	params := &ClientParams{
		RPCEndpoint: srv.URL,
	}

	client := NewClient(params)

	blockNumber, err := client.GetCurrentBlock(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedBlockNumber := big.NewInt(436)
	if blockNumber.Cmp(expectedBlockNumber) != 0 {
		t.Errorf("unexpected block number: want %d, got %d", expectedBlockNumber, blockNumber)
	}
}

func TestClient_GetBlockByNumber_InvalidBigInt(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header to be application/json, got %s", r.Header.Get("Content-Type"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"jsonrpc": "2.0",
			"result": {
			  "hash": "0x9f659488bb5df216187df081f5b07ab77319d3e207a7bd5a0ad4faab2fee4318",
			  "number": "0x1b4axxx",
			  "timestamp": "0x66304a4b",
			  "transactions": [
				{				  
				  "blockNumber": "0x1b4",
				  "from": "0x6f19a061fac475b5baee49c17069661eb6d76ea9",
				  "gas": "0x87015",
				  "hash": "0xfe62c8f139de9b490c6cdf85761776e656d6d394fb4d8eb8c57b76c3e1db49e3",
				  "to": "0xb591b1c3e63e093bf921fdd80ce2cf97c07f6936",
				  "value": "0x1b4"
				}
			  ]
			},
			"id": 1
		  }`))
	}))
	defer srv.Close()

	params := &ClientParams{
		RPCEndpoint: srv.URL,
	}

	client := NewClient(params)

	_, err := client.GetBlockByNumber(context.Background(), big.NewInt(436), true)

	if err == nil {
		t.Fatal("we want fail, but works wtf")
	}

}

func TestClient_GetBlockByNumber_RPCError(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header to be application/json, got %s", r.Header.Get("Content-Type"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"jsonrpc": "2.0",
			"result": {},
			"id": 1,
			"error": { 
				"code": 1,
				"message": "some error"
			}
		  }`))
	}))
	defer srv.Close()

	params := &ClientParams{
		RPCEndpoint: srv.URL,
	}

	client := NewClient(params)

	_, err := client.GetBlockByNumber(context.Background(), big.NewInt(436), true)

	if err == nil {
		t.Fatal("we want fail, but works wtf")
	}

}
