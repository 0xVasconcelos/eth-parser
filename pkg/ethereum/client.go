package ethereum

import (
	"net/http"
	"time"
)

const (
	MethodBlockNumber      = "eth_blockNumber"
	MethodGetBlockByNumber = "eth_getBlockByNumber"
)

type Client struct {
	// rpcEndpoint is the endpoint of the Ethereum RPC server.
	rpcEndpoint string
	// hc is the HTTP client used to communicate with the Ethereum RPC server.
	hc *http.Client
}

type ClientParams struct {
	// RPCEndpoint is the endpoint of the Ethereum RPC server.
	RPCEndpoint string
	// Timeout is the timeout for the HTTP client.
	Timeout time.Duration
	// HttpClient is the HTTP client used to communicate with the Ethereum RPC server.
	HttpClient *http.Client
}

func NewClient(params *ClientParams) *Client {
	if params.Timeout == 0 {
		params.Timeout = 10 * time.Second
	}
	if params.HttpClient == nil {
		params.HttpClient = &http.Client{Timeout: params.Timeout}
	}
	return &Client{
		rpcEndpoint: params.RPCEndpoint,
		hc:          params.HttpClient,
	}
}
