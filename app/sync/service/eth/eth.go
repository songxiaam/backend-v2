package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/zeromicro/go-zero/core/logx"
	"metaLand/data/model/chain"

	"github.com/ethereum/go-ethereum/ethclient"
)

type EthClients struct {
	Clients map[uint64]*Client
}

type ChainInfo struct {
	ChainID                uint64 `json:"chain_id"`
	RPCURL                 string `json:"rpc_url"`
	WSSURL                 string `json:"wss_url"`
	StartupContractAddress string `json:"startup_contract_address"`
	Abi                    string `json:"abi"`
}

type Client struct {
	RPCClient *ethclient.Client
	WSSClient *ethclient.Client
	ChainInfo *ChainInfo
}

//var Clients map[uint64]*Client

var EthSubChanel = make(chan struct{})

func NewEthClients() *EthClients {
	return &EthClients{Clients: make(map[uint64]*Client)}
}

func (e *EthClients) GetClient(chainID uint64) (*Client, error) {
	if client, ok := e.Clients[chainID]; ok {
		return client, nil
	}
	return nil, fmt.Errorf("error: GetEthClient error; chain id: %v", chainID)
}

func (e *EthClients) Start(chains *[]chain.ChainBasicResponse) {

	//var chains []chain.ChainBasicResponse
	//err := chain.GetChainCompleteList(e.ctx.DB, &chains)
	//if err != nil {
	//	logx.Error(err)
	//	return
	//}
	clients := make(map[uint64]*Client)
	for _, chainItem := range *chains {

		clients[chainItem.ChainID] = &Client{}
		chainInfo := &ChainInfo{
			ChainID: chainItem.ChainID,
		}
		for _, endpoint := range chainItem.ChainEndpoints {
			if endpoint.Protocol == 1 {
				chainInfo.RPCURL = endpoint.URL
			} else {
				chainInfo.WSSURL = endpoint.URL
			}
		}
		for _, contract := range chainItem.ChainContracts {
			if contract.Project == 1 && contract.Type == 1 {
				chainInfo.StartupContractAddress = contract.Address
				chainInfo.Abi = contract.ABI
			}
		}
		clients[chainItem.ChainID].ChainInfo = chainInfo
		setClient(clients[chainItem.ChainID])
	}
}

// Init the eth client
func setClient(client *Client) {
	var err error

	// 介于有些节点找不到WSS服务 故不再使用 WSS
	// log.Info("eth.Init ethclient_wss.Dial:", client.ChainInfo.WSSURL)
	// client.WSSClient, err = ethclient.Dial(client.ChainInfo.WSSURL)
	// if err != nil {
	// 	client.WSSClient = nil
	// 	log.Warn(err)
	// }
	logx.Info("eth.Init ethclient_rpc.Dial:", client.ChainInfo.RPCURL)
	client.RPCClient, err = ethclient.Dial(client.ChainInfo.RPCURL)
	if err != nil {
		client.RPCClient = nil
		logx.Error(err)
	}
	logx.Info("eth.Init ethclient.Dial done")
}

func (t *EthClients) Close() {
	logx.Info("eth.Close start")
	for _, client := range t.Clients {
		if client.WSSClient != nil {
			client.WSSClient.Close()
		}
	}
	EthSubChanel = make(chan struct{})
	log.Info("eth.Close end")
}
