# Nodes

--

Enhance your blockchain experience with effortless node execution, thanks to our intuitive API.
Enjoy a frictionless process as you effortlessly run and manage blockchain nodes with ease.


### How to run nodes?
--

1. `make build-node CHAIN=ethereum`
2. `make compose-node-up CHAIN=ethereum`


## API V2

#### **/api/v2/nodes**

**Request**

> Ethereum-Ganache

```json
{
	"network": "ethereum",
	"envVars": {
		"ENVIRONMENT": "development",
		"HOST": "0.0.0.0",
		"NETWORK_URL": "https://patient-delicate-pine.quiknode.pro/4200300eae9e45c661df02030bac8bc34f8b618e/",
		"BASE_CHAIN_DATA_PATH": "data",
		"RPC_PORT": "8545",
		"FROM_BLOCK_NUMBER": "17000000"
	}
}
```

> Chainlink

```json
{
	"network": "chainlink",
	"envVars": {
		"ENVIRONMENT": "sepolia",
    "ETH_URL": "wss://ethereum.sepolia.darchlabs.com/481affad55cac7efcbcc1182e4e435107aee7fae/",
    "PASSWORD": "ThisIsSecurePassword",
    "NODE_EMAIL": "dev@darchlabs.com",
    "NODE_EMAIL_PWD": "ThisIsSecurePassword"
	}
}
```

> Celo

```json
{
	"network": "celo",
	"envVars": {
		"ENVIRONMENT": "alfajores",
		"PASSWORD": "ThisIsSecurePassword"
	}
}
```
