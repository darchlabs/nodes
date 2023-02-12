package client

type NodeStatus struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Chain           string `json:"chain"`
	Port            int    `json:"port"`
	FromBlockNumber int64  `json:"fromBlockNumber"`
	Status          string `json:"status"`
}

type GetStatusHandlerResponse struct {
	Nodes []*NodeStatus `json:"nodes"`
}
