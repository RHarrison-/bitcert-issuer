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

// opReturnScript create the script to apply to the unspendable output where the merkle root will be stored
func opReturnScript(data []byte) []byte {
	builder := txscript.NewScriptBuilder()
	script, err := builder.AddOp(txscript.OP_RETURN).AddData(data).Script()
	if err != nil {
		panic(err)
	}
	return script
}

// addOutput adds a new output to a bitcoin transaction
func (tx Transaction) addOutput(amount int64, address string) {
	// turn adress string into valid Address obj
	destinationAddress, err := btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	if err != nil {
		log.Print("yo waddup")
	}
	script, _ := txscript.PayToAddrScript(destinationAddress)
	tx.util.AddTxOut(wire.NewTxOut(amount, script))
	tx.util.AddTxOut(wire.NewTxOut(0, opReturnScript([]byte("reece you bloody stud you"))))

}

// addInput adds a new Input to a bitcoin transaction
func (tx Transaction) addInput() {
	hash, _ := chainhash.NewHashFromStr("d8d6414a8a42e04761247f8ad41e73f2f77c961123341643759fb5954f746883") // TODO(Reece) need to get UTXO using electrum
	prevOut := wire.NewOutPoint(hash, 0)
	redeemTxIn := wire.NewTxIn(prevOut, nil, nil)
	tx.util.AddTxIn(redeemTxIn)
}

// Sign signs the transaction and adds the scripts too the inputs and outputs
func (tx Transaction) Sign(privKey string, destAddr string) {
	fmt.Println("starting sign")
	// decode WIF wallet privkey
	wif, _ := btcutil.DecodeWIF(privKey)

	// get the pub key for the given private key
	fmt.Println(wif.PrivKey.PubKey().SerializeCompressed())
	addresspubkey, _ := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), &chaincfg.TestNet3Params)
	// get sending address
	sendingAddress, _ := btcutil.DecodeAddress(addresspubkey.EncodeAddress(), &chaincfg.TestNet3Params)

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
func CreateTransaction(privKey string, paymentAddress string) (string, error) {
	var tx Transaction
	tx.util = wire.NewMsgTx(wire.TxVersion)

	tx.addInput()
	tx.addOutput(90532339, paymentAddress)

	tx.Sign(privKey, paymentAddress)

	var signedTx bytes.Buffer
	tx.util.Serialize(&signedTx)
	SignedTx := hex.EncodeToString(signedTx.Bytes())

	return SignedTx, nil
}

func test() {
	address := "mzq1guetvcVB1Cqpu7UNXb7Y3Kcme47z7x"
	privKey := "cSATNwZACSQ2Zk24Fzts2RbmV3LU2JBiz8QVhT8DxqV9cxPK43it"

	transaction, err := CreateTransaction(privKey, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(transaction)
}
