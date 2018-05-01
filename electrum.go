package main

import (
	"errors"
	"log"

	"github.com/d4l3k/go-electrum/electrum"
)

// Wallet ...
type Wallet struct {
	address string
	key     string
}

// Wallets ...
var Wallets []Wallet

func importWallets(wallets []Wallet) {
	for _, wallet := range wallets {
		Wallets = append(Wallets, wallet)
	}
}

func getUsableWallet() (Wallet, error) {
	var wallet Wallet

	for _, wallet := range Wallets {
		balance := net.balance(wallet.address)

		if balance.Confirmed > 3000 {
			return wallet, nil
		}
	}

	err := errors.New("no wallet with sufficient confirmed balance")
	return wallet, err
}

// Network ...
type Network struct {
	url  string
	node *electrum.Node
}

var net = Network{}

func (n *Network) connect(testnet bool) {
	if testnet {
		// testnet electrum url
		n.url = "testnet.hsmiths.com:53011"
	} else {
		// testnet electrum url
		n.url = "electrum.anduck.net:50001"
	}

	node := electrum.NewNode()
	if err := node.ConnectTCP(n.url); err != nil {
		log.Fatal(err)
	}
	n.node = node
}

func (n *Network) balance(address string) *electrum.Balance {
	balance, _ := n.node.BlockchainAddressGetBalance(address)
	return balance
}

func (n *Network) estimateFee(blocks int) int {
	fee, _ := n.node.BlockchainEstimateFee(blocks)
	return int(fee * 100000000)
}

// UTXO ...
func (n *Network) UTXO(address string) []*electrum.Transaction {
	response, _ := n.node.BlockchainAddressListUnspent(address)
	return response
}

func (n *Network) broadcast(raw string) interface{} {
	result, _ := n.node.BlockchainTransactionBroadcast([]byte(raw))
	return result
}
