package main

import (
	"fmt"
	"log"

	"github.com/d4l3k/go-electrum/electrum"
)

type Network struct {
	url  string
	node *electrum.Node
}

func (n *Network) connect() {
	fmt.Println("start connect")
	node := electrum.NewNode()
	if err := node.ConnectTCP(n.url); err != nil {
		log.Fatal(err)
	}
	n.node = node
}

func (n *Network) balance(address string) {
	balance, err := n.node.BlockchainAddressGetBalance(address)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Address balance: %+v", balance.Confirmed)
}

func (n *Network) estimateFee(blocks string) {
	response, err := n.node.BlockchainEstimateFee("10")

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Estimated Fee: %+v", response)
}

func (n *Network) UTXO(address string) {
	response, err := n.node.BlockchainAddressListUnspent(address)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Estimated Fee: %+v", response[0].Hash)
}
