package types

type JsonRpcData struct {
	JsonRpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      string   `json:"id"`
}

type JsonRpcResult struct {
	JsonRpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Result  string `json:"result"`
}
