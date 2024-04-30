# eth parser


A simple Ethereum parser that subscribes to new blocks and index transaction to given addresses.

## Getting Started
- Install Go
- Clone this repository and run inside it the following command to download the dependencies
```bash
go mod download
```


## Configure
Please set the RPC endpoint and the timeout in the client parameter in `main.go`
```go
params := &ethereum.ClientParams{
	RPCEndpoint: "https://ethereum.example.com/v1/mainnet",
	Timeout:     5 * time.Second,
}
```

## Test
```bash
go test ./...
```

## Run

```bash
go run cmd/parser/main.go
```

## Methods and Types

### Parser package

#### Methods 

```go 
GetCurrentBlock(context.Context) (*big.Int, error)
Subscribe(context.Context, string) bool
GetTransactions(context.Context, string) ([]ethereum.Transaction, error)
Start(ctx context.Context, ticker *time.Ticker)
```


### Ethereum package

#### Types
```go
type Client struct {
	rpcEndpoint string
	hc *http.Client
}

type ClientParams struct {
	RPCEndpoint string
	Timeout time.Duration
	HttpClient *http.Client
}

type Block struct {
	Number       BigInt        `json:"number"`
	Hash         string        `json:"hash"`
	Timestamp    BigInt        `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	BlockNumber BigInt `json:"blockNumber"`
	From        string `json:"from"`
	Hash        string `json:"hash"`
	To          string `json:"to"`
	Value       BigInt `json:"value"`
	Gas         BigInt `json:"gas"`
}

type BigInt struct {
	*big.Int
}

type GetByBlockNumberParams struct {
	BlockNumber         *big.Int
	IncludeTransactions bool
}

type GetByBlockNumberResponse struct {
	RPCResponse
	Result Block `json:"result"`
}

type GetCurrentBlockResponse struct {
	RPCResponse
	Result string `json:"result"`
}

```

#### Methods

```go
NewClient(params *ClientParams) *Client
GetBlockByNumber(ctx context.Context, blockNumber *big.Int, includeTxs bool) (*Block, error)
GetCurrentBlock(ctx context.Context) (*big.Int, error)
```