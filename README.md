# Nodes

--

Run nodes made easy. With nodes you will run blockchain nodes on deman through a simple an easy to use API.

Nodes support evm blockchains for the moment, you only need to provide the desired chain url in order to fork them.

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
		"ETH_URL": "https://ethereum-sepolia.darchlabs.com",
		"PASSWORD": "ThisIsSecurePassword",
		"NODE_EMAIL": "dev@darchlabs.com",
		"NODE_EMAIL_PWD": "thisIsPassword"
	}
}
```