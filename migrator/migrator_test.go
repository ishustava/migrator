package migrator

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	. "github.com/ishustava/migrator/test_fixtures"
)

var _ = Describe("Migrator", func() {
	Describe("#parseVarsStoreFile", func() {
		Context("Passwords", func() {
			It("finds and returns password credentials", func() {
				password1 := makePassword("path1", "password1")
				password2 := makePassword("path2", "password2")

				creds, err := parseVarsStoreFile("../test_fixtures/valid_creds.yml")

				Expect(err).ToNot(HaveOccurred())
				Expect(creds.Passwords).To(ConsistOf(password1, password2))
			})
		})

		Context("Certificates", func() {
			It("finds and returns certificate credentials", func() {
				cert1 := makeCertificate("path3", values.Certificate{Ca: CA1, Certificate: CERT1, PrivateKey: PRIV1})
				cert2 := makeCertificate("path4", values.Certificate{Certificate: CERT2, PrivateKey: PRIV2})

				creds, err := parseVarsStoreFile("../test_fixtures/valid_creds.yml")

				Expect(err).ToNot(HaveOccurred())
				Expect(creds.Certificates).To(ConsistOf(cert1, cert2))
			})
		})

		Context("SSH", func() {
			It("finds and returns ssh credentials", func() {
				ssh := makeSsh("path5", values.SSH{PublicKey: SSH_PUB, PrivateKey: SSH_PRIV})

				creds, err := parseVarsStoreFile("../test_fixtures/valid_creds.yml")

				Expect(err).ToNot(HaveOccurred())
				Expect(creds.SshKeys).To(ConsistOf(ssh))
			})
		})

		Context("RSA", func() {
			It("finds and returns rsa credentials", func() {
				rsa := makeRsa("path6", values.RSA{PublicKey: RSA_PUB, PrivateKey: RSA_PRIV})

				creds, err := parseVarsStoreFile("../test_fixtures/valid_creds.yml")

				Expect(err).ToNot(HaveOccurred())
				Expect(creds.RsaKeys).To(ConsistOf(rsa))
			})
		})
	})
})
