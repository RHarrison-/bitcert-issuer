package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/rharrison-/merkle"
)

// Certificate  ... __
type Certificate struct {
	Type      string            `json:"type"`
	Badge     badge             `json:"badge"`
	IssuedOn  string            `json:"issuedOn"`
	Recipient recipient         `json:"recipient"`
	Signature *merkle.Signature `json:"signature"`
}

type recipient struct {
	Type     string `json:"type"`
	Identity string `json:"identity"`
	Hashed   bool   `json:"hashed"`
}

type badge struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Image          string         `json:"image"`
	Issuer         issuer         `json:"issuer"`
	Criteria       criteria       `json:"criteria"`
	Description    string         `json:"description"`
	SignatureLines badgeSignature `json:"signatureLines"`
	Type           string         `json:"type"`
}

type issuer struct {
	ID             string `json:"id"`
	URL            string `json:"url"`
	Type           string `json:"type"`
	Email          string `json:"email"`
	Image          string `json:"image"`
	RevocationList string `json:"revocationList"`
}

type criteria struct {
	Narrative string `json:"narattive"`
}

type badgeSignature struct {
	Name     string   `json:"name"`
	Image    string   `json:"image"`
	JobTitle string   `json:"jobTitle"`
	Type     []string `json:"type"`
}

// func (c certificate) isSigned() bool {
//     if c.Signature != nil {
//         return true
//     }
//     return false
// }

type batch []Certificate

// hash certs, craete merkle root, add proofs to certificates.
func (b *batch) Sign() {
	for _, cert := range *b {
		fmt.Println(cert.Type)
	}
}

func (b *batch) add(cert Certificate) {
	*b = append(*b, cert)
}

func loadCertificate(path string) error {
	var cert Certificate
	data, _ := ioutil.ReadFile(path)

	_ = json.Unmarshal(data, &cert)

	// fmt.Println("cert.type", cert.Type)

    var newBatch = batch{}
    fmt.Println(len(newBatch))
    newBatch.add(cert)

	fmt.Println(len(newBatch), newBatch[0].Type)

	return nil
}

// make sure nessercary fields for signing are present
func (c Certificate) validateIntegrity() {

}
