package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/ishustava/migrator/test_fixtures"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	"github.com/ishustava/migrator/credentials"
	"github.com/ishustava/migrator/parser"
)

var _ = Describe("Parser", func() {
	Describe("#ParseVarsStoreFile", func() {
		Context("Passwords", func() {
			It("finds and returns password credentials", func() {
				password1 := credentials.MakePassword("path1", "password1")
				password2 := credentials.MakePassword("path2", "password2")

				creds, err := parser.ParseVarsStoreFile("../test_fixtures/valid_creds.yml")

				Expect(err).ToNot(HaveOccurred())
				Expect(creds.Passwords).To(ConsistOf(password1, password2))
			})
		})

		Context("Certificates", func() {
			It("finds and returns certificate credentials", func() {
				cert1 := credentials.MakeCertificate("path3", values.Certificate{Ca: CA1, Certificate: CERT1, PrivateKey: PRIV1})
				cert2 := credentials.MakeCertificate("path4", values.Certificate{Certificate: CERT2, PrivateKey: PRIV2})

				creds, err := parser.ParseVarsStoreFile("../test_fixtures/valid_creds.yml")

				Expect(err).ToNot(HaveOccurred())
				Expect(creds.Certificates).To(ConsistOf(cert1, cert2))
			})
		})

		Context("SSH", func() {
			It("finds and returns ssh credentials", func() {
				ssh := credentials.MakeSsh("path5", values.SSH{PublicKey: SSH_PUB, PrivateKey: SSH_PRIV})

				creds, err := parser.ParseVarsStoreFile("../test_fixtures/valid_creds.yml")

				Expect(err).ToNot(HaveOccurred())
				Expect(creds.SshKeys).To(ConsistOf(ssh))
			})
		})

		Context("RSA", func() {
			It("finds and returns rsa credentials", func() {
				rsa := credentials.MakeRsa("path6", values.RSA{PublicKey: RSA_PUB, PrivateKey: RSA_PRIV})

				creds, err := parser.ParseVarsStoreFile("../test_fixtures/valid_creds.yml")

				Expect(err).ToNot(HaveOccurred())
				Expect(creds.RsaKeys).To(ConsistOf(rsa))
			})
		})
	})
})
