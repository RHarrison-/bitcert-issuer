package main

import (
	"encoding/json"
	"fmt"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
)

type createTransactionStruct struct {
	Address string `json:"address"`
	PrivKey string `json:"privKey"`

	UTXO string `json:"UTXO"`
}

type returnable struct {
	hex string
}

// unPack safely unmarshals the json data
func unPack(message bootstrap.MessageIn, data interface{}) (interface{}, error) {
	if len(message.Payload) > 0 {
		// Unmarshal payload
		err := json.Unmarshal(message.Payload, &data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {

	switch m.Name {
	case "loadCertificates":
		var data []fileData

		unPack(m, &data)
		batch = Batch{}
		loadCertificates(data)

		return batch, nil

	case "createTx":
		var data Batch

		unPack(m, &data)
		batch = data

		merkleRoot := batch.attachProofs()

		wallet, _ := getUsableWallet()
		txHex, txHash, _ := CreateTransaction(wallet, merkleRoot)
		batch.addAnchor(txHash.String())

		fmt.Println("new tx hash: ", txHash)

		batch.save()

		fmt.Println(txHex)

		result := net.broadcast(txHex)
		fmt.Println(result)

		return batch, nil

	case "verifyCertificate":
		var certificate Certificate
		unPack(m, &certificate)

		test, err := verify(certificate)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("---------------------------")
		fmt.Println(test)
		fmt.Println("---------------------------")

		return test, nil

	default:
		fmt.Println("fucntion not defined: " + m.Name)

	}
	return
}
