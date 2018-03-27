package pki_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPki(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pki Suite")
}
