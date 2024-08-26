package core

import (
	"encoding/json"
	"errors"
	"net/url"
	"sync"
	"sync/atomic"
)

// SIMPLE Balancer interface
type Balancer interface {
	NextServer() *url.URL
}

// Regular RoundRobin implementation
type RoundRobinBalancer struct {
	servers []ServerEndpoint
	len     uint64
	index   atomic.Uint64
}

func NewRoundRobinBalancer(servers []ServerEndpoint) (*RoundRobinBalancer, error) {
	if len(servers) == 0 {
		return nil, errors.New("empty servers list")
	}

	return &RoundRobinBalancer{
		servers: servers,
		len:     uint64(len(servers)),
	}, nil
}

func (b *RoundRobinBalancer) NextServer() *url.URL {
	// after MAX uint32 will set to 0 again
	currentIndex := b.index.Add(1) % b.len

	return &b.servers[currentIndex].URL
}

// Weighted RoundRobin implementation
type WeightedRoundRobinBalancer struct {
	servers   []ServerEndpoint
	maxWeight uint64
	len       uint64

	mu             sync.Mutex
	index          uint64
	currentWeights []uint64
}

func NewWeightedRoundRobinBalancer(servers []ServerEndpoint) (*WeightedRoundRobinBalancer, error) {
	if len(servers) == 0 {
		return nil, errors.New("empty servers list")
	}

	var maxWeight, totalWeight uint64

	currentWeights := make([]uint64, len(servers))

	for _, server := range servers {
		if server.Weight > maxWeight {
			maxWeight = server.Weight
		}

		totalWeight += server.Weight
	}

	if totalWeight == 0 {
		return nil, errors.New("total weight of servers is zero")
	}

	return &WeightedRoundRobinBalancer{
		servers:   servers,
		maxWeight: uint64(maxWeight),
		len:       uint64(len(servers)),

		currentWeights: currentWeights,
	}, nil
}

func (b *WeightedRoundRobinBalancer) NextServer() *url.URL {
	b.mu.Lock()
	defer b.mu.Unlock()

	for {
		b.index = (b.index + 1) % b.len
		if b.index == 0 {
			for i := range b.currentWeights {
				b.currentWeights[i] += b.servers[i].Weight
			}
		}

		if b.currentWeights[b.index] >= b.maxWeight {
			b.currentWeights[b.index] -= b.maxWeight
			return &b.servers[b.index].URL
		}
	}
}

type ServerEndpoint struct {
	URL    url.URL `json:"url"`
	Weight uint64  `json:"weight,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface to fix url.URL
func (s *ServerEndpoint) UnmarshalJSON(data []byte) error {
	type Alias ServerEndpoint
	aux := &struct {
		URL string `json:"url"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	parsedURL, err := url.Parse(aux.URL)
	if err != nil {
		return err
	}
	s.URL = *parsedURL
	return nil
}
