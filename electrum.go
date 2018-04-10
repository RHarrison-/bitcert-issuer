package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/d4l3k/go-electrum/electrum"
)

// Wallet ...
type Wallet struct {
	address string
	key     string
}

// Wallets ...
var Wallets []Wallet

func importWallets(path string) {
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "	")
		t := strings.Split(s[1], ":")
		address := s[0]
		privKey := t[1]

		var wallet = Wallet{
			address: address,
			key:     privKey,
		}

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

var net = Network{
	url: "testnet.hsmiths.com:53011", //electrum.anduck.net:50001
}

func (n *Network) connect() {
	fmt.Println("start connect")
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
