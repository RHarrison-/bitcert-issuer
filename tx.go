package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type Transaction struct {
	util *wire.MsgTx
}

// addOutput adds a new output to a bitcoin transaction
func (tx Transaction) addOutput(satoshis int64) {
	// turn adress string into valid Address obj
	fmt.Println("adding output")
	tx.util.AddTxOut(wire.NewTxOut(satoshis, nil))
	fmt.Println("done")
}

// addInput adds a new Input to a bitcoin transaction
func (tx Transaction) addInput(txid string) {
	fmt.Println("adding input")
	hash, _ := chainhash.NewHashFromStr(txid)
	prevOut := wire.NewOutPoint(hash, 0)
	redeemTxIn := wire.NewTxIn(prevOut, nil, nil)
	tx.util.AddTxIn(redeemTxIn)
	fmt.Println("done")
}

// Sign signs the transaction and adds the scripts too the inputs and outputs
func (tx Transaction) Sign(privKey string, destAddr string) {
	fmt.Println("starting sign")
	// decode WIF wallet privkey
	wif, err := btcutil.DecodeWIF(privKey)
	destinationAddress, err := btcutil.DecodeAddress(destAddr, &chaincfg.TestNet3Params)
	if err != nil {
		log.Print("yo waddup")
	}
	tx.util.TxOut[0].PkScript, _ = txscript.PayToAddrScript(destinationAddress)

	// get the pub key for the given private key
	fmt.Println(wif.PrivKey.PubKey().SerializeCompressed())
	addresspubkey, _ := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), &chaincfg.TestNet3Params)
	// get sending address
	sendingAddress, err := btcutil.DecodeAddress(addresspubkey.EncodeAddress(), &chaincfg.TestNet3Params)

	UTXOscript, _ := txscript.PayToAddrScript(sendingAddress)
	tx.util.TxIn[0].SignatureScript, _ = txscript.SignatureScript(tx.util, 0, UTXOscript, txscript.SigHashAll, wif.PrivKey, true)
	fmt.Println("done")
}

// todo(Reece)
// func (tx Transaction) verify() {
// 	flags := txscript.StandardVerifyFlags

// 	vm, err := txscript.NewEngine(UTXOscript, tx.util, 0, flags, nil, nil, satoshis)
// 	if err != nil {
// 		return "Transaction{}", err
// 	}

// 	if err := vm.Execute(); err != nil {
// 		return "Transaction{}", err
// 	}
// }

// Create a new transaction and return the final transaction hex ready to broadcast
func CreateTransaction(privKey string, paymentAddress string, satoshis int64, txid string) (string, error) {
	var tx Transaction
	tx.util = wire.NewMsgTx(wire.TxVersion)

	tx.addInput(txid)
	tx.addOutput(satoshis)

	tx.Sign(privKey, paymentAddress)

	var signedTx bytes.Buffer
	tx.util.Serialize(&signedTx)
	SignedTx := hex.EncodeToString(signedTx.Bytes())

	return SignedTx, nil
}

func main() {
	address := "myPqMJLLgXxAWsXXHec5waTJwR9G5EdqKq"
	privKey := "cS1L8sc7RD2FXUUXeevmaM4NRp3HxAECuTdBdVqxiHK5t4FSuyvi"
	input_txHash := "d7bf58a5a3a52623a20e7627aa9f3241528d9be8e622ec9be14e4586862018fe"

	transaction, err := CreateTransaction(privKey, address, 27000000, input_txHash)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(transaction)
}
