package main

import (
	"testing"
)

func TestCertHandling(t *testing.T) {
	err := loadCertificate("/home/harrison/go/src/github.com/rharrison-/bitcert-issuer/example-certs/test_cert_unsigned.json")
	if err != nil {
		t.Error("error: unexpected error:  ", err)
	}
}
