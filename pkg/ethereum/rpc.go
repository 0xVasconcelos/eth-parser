package ethereum

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RPCRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int         `json:"id"`
}

type RPCResponse struct {
	Jsonrpc string    `json:"jsonrpc"`
	Id      int       `json:"id"`
	Error   *RPCError `json:"error"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"data"`
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("%d - %s: %s", e.Code, e.Message, e.Detail)
}

func (c *Client) post(ctx context.Context, r *RPCRequest, result interface{}) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(data)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.rpcEndpoint, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, status: %s", resp.StatusCode, resp.Status)
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &result)
	if err != nil {
		return err
	}

	return nil
}
