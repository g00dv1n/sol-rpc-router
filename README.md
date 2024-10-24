# Solana RPC Router & Balancer

The Solana RPC Router is a tool designed to efficiently route and balance traffic between regular Solana RPC servers and those that support the DAS API. The router supports simple round-robin and weighted round-robin balancing strategies. Keep in mind that this is just a proof of concept (PoC), but it works.

## Usage

```sh
make build
./bin/sol-rpc-router -c proxy_config.json
```

Set up the router behind NGINX or in a Docker container, and use it as a universal endpoint for all your RPC calls within your app.

## Config Example

`proxy_config.json`

```json
{
  "port": 9999,
  "host": "127.0.0.1",
  "regularRpc": {
    "balancerType": "rr",
    "servers": [
      {
        "url": "https://solana-mainnet.core.chainstack.com/API_KEY"
      },
      {
        "url": "https://solana-mainnet.core.chainstack.com/API_KEY"
      }
    ]
  },
  "dasRpc": {
    "balancerType": "wrr",
    "servers": [
      {
        "url": "https://mainnet.helius-rpc.com/?api-key=API_KEY",
        "weight": 2
      },
      {
        "url": "https://mainnet.helius-rpc.com/?api-key=API_KEY",
        "weight": 1
      }
    ]
  }
}
```

Made for [Solana Hacker Telegram Bot](https://hackers.tools/d/solana-hacker-bot).
