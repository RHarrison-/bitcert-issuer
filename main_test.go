package main

import "testing"

// func TestCertHandling(t *testing.T) {
// 	cert := loadCertificate("/home/harrison/go/src/github.com/rharrison-/bitcert-issuer/example-certs/test_cert_unsigned.json")
// 	var batch Batch
// 	batch.add(cert)
// 	batch.test()
// }

// func TestElectrum(t *testing.T) {
// 	var net = Network{
// 		url: "testnet.hsmiths.com:53011",
// 	}

// 	net.connect()
// 	// net.balance("mv62wDMoSYfiNUE5LsqY9DgN9kSagAK9vg")
// 	net.UTXO("mv62wDMoSYfiNUE5LsqY9DgN9kSagAK9vg")
// }

func TestTx(t *testing.T) {
	test()
}
