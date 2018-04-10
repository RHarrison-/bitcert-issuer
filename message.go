package main

import (
	"encoding/json"
	"fmt"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
)

type createTransactionStruct struct {
	address string `json:"address"`
	privKey string `json:"privKey"`

	UTXO string `json:"UTXO"`
}

type returnable struct {
	hex string
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {

	switch m.Name {
	case "addCertificates":
		var data []fileData
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &data); err != nil {
				payload = err.Error()
				return
			}
		}

		loadCertificates(data)

		fmt.Println(batch[0].isSigned())
		fmt.Println(batch[1].isSigned())

	default:
		fmt.Println("fucntion not defined: " + m.Name)

	}
	return
}
