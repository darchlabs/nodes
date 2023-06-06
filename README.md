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
