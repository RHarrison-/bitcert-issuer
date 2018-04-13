package main

import (
	"bytes"
	"encoding/hex"
	"log"

	"github.com/d4l3k/go-electrum/electrum"

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
	tx.util.AddTxOut(wire.NewTxOut(0, opReturnScript([]byte("another example"))))

}

func getUTXO(address string) []*electrum.Transaction {
	return net.UTXO(address)
}

func calculateFee(blocks int) int {
	fee := net.estimateFee(blocks)
	return fee
}

// addInput adds a new Input to a bitcoin transaction
func (tx Transaction) addInput(UTXO electrum.Transaction) {
	hash, _ := chainhash.NewHashFromStr(UTXO.Hash) // TODO(Reece) need to get UTXO using electrum
	prevOut := wire.NewOutPoint(hash, uint32(UTXO.Pos))
	redeemTxIn := wire.NewTxIn(prevOut, nil, nil)
	tx.util.AddTxIn(redeemTxIn)
}

// Sign signs the transaction and adds the scripts too the inputs and outputs
func (tx Transaction) Sign(privKey string, destAddr string) {
	// decode WIF wallet privkey
	wif, _ := btcutil.DecodeWIF(privKey)

	// get the pub key for the given private key
	addresspubkey, _ := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), &chaincfg.TestNet3Params)

	// get sending address
	sendingAddress, _ := btcutil.DecodeAddress(addresspubkey.EncodeAddress(), &chaincfg.TestNet3Params)
	UTXOscript, _ := txscript.PayToAddrScript(sendingAddress)
	tx.util.TxIn[0].SignatureScript, _ = txscript.SignatureScript(tx.util, 0, UTXOscript, txscript.SigHashAll, wif.PrivKey, true)
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

// CreateTransaction Create a new transaction and return the final transaction hex ready to broadcast
func CreateTransaction(wallet Wallet) (string, chainhash.Hash, error) {
	var tx Transaction
	tx.util = wire.NewMsgTx(wire.TxVersion)

	utxoArray := getUTXO(wallet.address)
	fee := calculateFee(3)
	outputAmount := utxoArray[0].Value - fee

	tx.addInput(*utxoArray[0])
	tx.addOutput(int64(outputAmount), wallet.address)
	tx.Sign(wallet.key, wallet.address)

	var signedTx bytes.Buffer
	tx.util.Serialize(&signedTx)
	transObject, _ := btcutil.NewTxFromBytes(signedTx.Bytes())
	SignedTx := hex.EncodeToString(signedTx.Bytes())

	return SignedTx, *transObject.Hash(), nil
}
