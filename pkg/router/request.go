package router

import "slices"

type RPCRequest struct {
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      any         `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
}

type RPCRequestType uint8

const (
	RegularRPCRequest RPCRequestType = 1
	DasRPCRequest     RPCRequestType = 3
)

var DasMethods = [...]string{
	"getAsset",
	"searchAssets",
	"getAssetProof",
	"getAssetsByGroup",
	"getAssetsByOwner",
	"getAssetsByCreator",
	"getAssetsByAuthority",
}

func (r *RPCRequest) RequestType() RPCRequestType {
	if r.IsDasRequest() {
		return DasRPCRequest
	}

	return RegularRPCRequest
}

func (r *RPCRequest) IsDasRequest() bool {
	return slices.Contains(DasMethods[:], r.Method)
}
