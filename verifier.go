package main

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rharrison-/merkle"
)

func verify(certificate Certificate) (bool, error) {

	if certificate.Signature == nil {
		err := errors.New("no signature in certificate")
		return false, err
	}

	hold := certificate.Signature
	certificate.Signature = nil

	hash := certificate.CalculateHash()

	fmt.Println("calculated hash: ", hex.EncodeToString(hash))
	fmt.Println("target hash: ", hold.TargetHash)

	decoded, _ := hex.DecodeString(hold.TargetHash)

	if bytes.Equal(decoded, hash) {
		proved := merkle.VerifyProof(*hold)
		if proved {
			return true, nil
		}
	} else {
		err := errors.New("certificate has been altered: calculated hash != target hash")
		return false, err
	}

	return false, nil
}
