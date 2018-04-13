package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/rharrison-/merkle"
)

// get dynamically
const root = "/home/harrison/go/src/github.com/rharrison-/bitcert-issuer"

// Certificate ..
type Certificate struct {
	// structure of json output
	Type      string            `json:"type,omitempty"`
	Badge     badge             `json:"badge,omitempty"`
	IssuedOn  string            `json:"issuedOn,omitempty"`
	Recipient recipient         `json:"recipient,omitempty"`
	Signature *merkle.Signature `json:"signature,omitempty"`

	Internal *internal `json:"internal,omitempty"`
}

type internal struct {
	Name   string
	Path   string
	Valid  bool
	Signed bool
}

type recipient struct {
	Name     string `json:"name"`
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
	c.Internal = nil
	hashable, _ := json.Marshal(c)
	fmt.Println(string(hashable))
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

func (b Batch) attachProofs() {

	tree, _ := merkle.NewTree(b.toContentArray())
	tree.GenerateProofs()
	batch = Batch{}

	for _, leaf := range tree.Leafs {
		if leaf.Dup {
			continue
		}

		var cert Certificate

		CertificateData, _ := json.Marshal(leaf.C)

		_ = json.Unmarshal(CertificateData, &cert)

		cert.Signature = &leaf.Proof

		batch.add(&cert)
	}
}

func (b *Batch) add(cert *Certificate) {
	*b = append(*b, cert)
}

func (b *Batch) addAnchor(txHash string) {
	for _, cert := range *b {
		cert.Signature.Anchors = append(cert.Signature.Anchors, createAnchor(txHash))
	}
}

func (b *Batch) save() {
	for _, cert := range *b {
		filename := cert.Internal.Name
		cert.Internal = nil
		jsonData, _ := json.Marshal(cert)

		// create signed directory if not existing
		if _, err := os.Stat(root + "/example-certs/signed"); os.IsNotExist(err) {
			os.Mkdir(root+"/example-certs/signed", os.ModePerm)
		}

		err := ioutil.WriteFile(root+"/example-certs/signed/"+filename, jsonData, 0644)

		if err != nil {
			fmt.Println("----------------------------------")
			fmt.Println("error saving signed cert")
			fmt.Println(err)
			fmt.Println(err.Error())
			fmt.Println("----------------------------------")

		}
	}
}

func createAnchor(txHash string) merkle.Anchor {
	return merkle.Anchor{
		SourceID: txHash,
		Type:     "BTCOpReturn",
	}
}

type fileData struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func loadCertificate(path string) *Certificate {
	var cert Certificate
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
	_ = json.Unmarshal(data, &cert)
	return &cert
}

func loadCertificates(files []fileData) {
	for _, file := range files {
		cert := loadCertificate(file.Path)

		fmt.Println(cert)
		splitPath := strings.Split(file.Path, "/")
		internal := internal{
			Path:   file.Path,
			Name:   splitPath[len(splitPath)-1],
			Signed: cert.isSigned(),
			Valid:  true,
		}

		cert.Internal = &internal

		batch.add(cert)
	}
}
