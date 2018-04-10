package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/rharrison-/merkle"
)

// Certificate ..
type Certificate struct {
	// structure of json output
	Type      string            `json:"type,omitempty"`
	Badge     badge             `json:"badge,omitempty"`
	IssuedOn  string            `json:"issuedOn,omitempty"`
	Recipient recipient         `json:"recipient,omitempty"`
	Signature *merkle.Signature `json:"signature,omitempty"`

	// internal processing values
	name  string
	valid bool
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

// CalculateHash ... returns the hash of the certificates json string
func (c Certificate) CalculateHash() []byte {
	hashable, _ := json.Marshal(c)
	h := sha256.New()
	h.Write(hashable)
	return h.Sum(nil)
}

// check if signature filed is present and verify its content
func (c Certificate) isSigned() bool {
	if c.Signature != nil {
		return true
	}
	return false
}

// make sure nessercary fields for signing are present
func (c Certificate) validateIntegrity() {

}

// Batch ... Array of Certificate
type Batch []*Certificate

func (b Batch) toContentArray() []merkle.Content {
	var content []merkle.Content
	for _, c := range b {
		content = append(content, c)
	}

	return content
}

func (b *Batch) test() {

	tree, _ := merkle.NewTree(b.toContentArray())

	tree.GenerateProofs()

	var c Certificate
	a, _ := json.Marshal(tree.Leafs[0].C)
	_ = json.Unmarshal(a, c)
	c.Signature = &tree.Leafs[0].Proof

	fmt.Println(merkle.VerifyProof(*c.Signature))

}

func attachProof(n []merkle.Node) {
	// for _, c := range b {
	// 	content = append(content, c)
	// }
}

func (b *Batch) add(cert *Certificate) {
	*b = append(*b, cert)
}

func loadCertificate(path string) *Certificate {
	var cert Certificate
	data, _ := ioutil.ReadFile(path)
	_ = json.Unmarshal(data, &cert)
	return &cert
}

type fileData struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func loadCertificates(files []fileData) {
	for _, file := range files {
		fmt.Println(file)
		cert := loadCertificate(file.Path)
		cert.name = file.Name
		batch.add(cert)
	}
}
