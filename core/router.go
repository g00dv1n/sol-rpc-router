package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type RpcRoutes map[RPCRequestType]Balancer

type RouterHandler struct {
	routes RpcRoutes
}

func NewRouter(regular Route, das Route) (*RouterHandler, error) {
	regularBalancer, err := regular.GetBalancer()
	if err != nil {
		return nil, errors.New("empty regular servers list")
	}

	dasBalancer, err := das.GetBalancer()

	if err != nil {
		return nil, errors.New("empty das servers list")
	}

	return &RouterHandler{
		routes: RpcRoutes{
			RegularRPCRequest: regularBalancer,
			DasRPCRequest:     dasBalancer,
		},
	}, nil
}

func (h *RouterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rpcReq := RPCRequest{}

	/// DECODE JSON and Make body copy to pass back later.
	bodyCopy := new(bytes.Buffer)
	json.NewDecoder(io.TeeReader(r.Body, bodyCopy)).Decode(&rpcReq)
	/// COPY Body back
	r.Body = io.NopCloser(bodyCopy)
	///

	/// PROXY SET UP
	targetRpcUrl := h.getProxyTarget(&rpcReq)
	proxy := httputil.NewSingleHostReverseProxy(targetRpcUrl)
	r.Host = targetRpcUrl.Host

	// let the proxy do her thing
	proxy.ServeHTTP(w, r)
}

func (h *RouterHandler) getProxyTarget(rpcReq *RPCRequest) *url.URL {
	return h.routes[rpcReq.GetRequestType()].NextServer()
}

type Route struct {
	BalancerType string           `json:"balancerType"`
	Servers      []ServerEndpoint `json:"servers"`
}

func (r Route) GetBalancer() (Balancer, error) {
	switch r.BalancerType {
	case "rr":
		return NewRoundRobinBalancer(r.Servers)
	case "wrr":
		return NewWeightedRoundRobinBalancer(r.Servers)

	default:
		return nil, errors.New("unknown route balancer type, use rr or wrr")
	}
}
