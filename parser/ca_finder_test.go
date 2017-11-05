package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ishustava/migrator/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	credentials2 "github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/ishustava/migrator/test_fixtures"
	"github.com/ishustava/migrator/parser"
	"errors"
)

var _ = Describe("CA Finder", func() {
	Context("with a root ca", func() {
		It("sets the ca name of the signer on the signee", func() {
			ca := credentials.NewCertificate(
				"root-ca",
					values.Certificate{Certificate: test_fixtures.ROOT_CA_CERT, PrivateKey: test_fixtures.ROOT_CA_PRIV},
				)
			cert1 := credentials.NewCertificate(
				"test-cert1",
				values.Certificate{Ca: test_fixtures.SIGNED_BY_ROOT_LEAF1_CA, Certificate: test_fixtures.SIGNED_BY_ROOT_LEAF1_CERT, PrivateKey: test_fixtures.SIGNED_BY_ROOT_LEAF1_PRIV},
			)
			cert2 := credentials.NewCertificate(
				"test-cert2",
				values.Certificate{Ca: test_fixtures.SIGNED_BY_ROOT_LEAF2_CA, Certificate: test_fixtures.SIGNED_BY_ROOT_LEAF2_CERT, PrivateKey: test_fixtures.SIGNED_BY_ROOT_LEAF2_PRIV},
			)

			certs, err := parser.FindAndSetSigningCA([]credentials2.Certificate{cert1, ca, cert2})
			Expect(err).ToNot(HaveOccurred())

			resultingCa, err := findCertByName(certs, "root-ca")
			Expect(err).ToNot(HaveOccurred())
			Expect(resultingCa.Value.CaName).To(BeEmpty())

			resultingCert, err := findCertByName(certs, "test-cert1")
			Expect(err).ToNot(HaveOccurred())
			Expect(resultingCert.Value.CaName).To(Equal("root-ca"))
			Expect(resultingCert.Value.Ca).To(BeEmpty())

			resultingCert, err = findCertByName(certs, "test-cert2")
			Expect(err).ToNot(HaveOccurred())
			Expect(resultingCert.Value.CaName).To(Equal("root-ca"))
			Expect(resultingCert.Value.Ca).To(BeEmpty())
		})
	})

	// BOSH CLI doesn't support generation of intermediate CAs
	// and instead always generates root for any certificate
	// with `is_ca: true` specified in the variables section
	// https://github.com/cloudfoundry/bosh-cli/issues/354
	XContext("with intermediate and root cas", func() {
		It("sets the ca name of the signer on the signee", func() {
			root := credentials.NewCertificate(
				"root-ca",
				values.Certificate{Certificate: test_fixtures.ROOT_CA_CERT, PrivateKey: test_fixtures.ROOT_CA_PRIV},
			)
			int := credentials.NewCertificate(
				"int-ca",
				values.Certificate{Ca: test_fixtures.INT_CA, Certificate: test_fixtures.INT_CERT, PrivateKey: test_fixtures.INT_PRIV},
			)
			leaf := credentials.NewCertificate(
				"leaf-cert",
				values.Certificate{Ca: test_fixtures.SIGNED_BY_INT_LEAF_CA, Certificate: test_fixtures.SIGNED_BY_INT_LEAF_CERT, PrivateKey: test_fixtures.SIGNED_BY_INT_LEAF_PRIV},
			)

			certs, err := parser.FindAndSetSigningCA([]credentials2.Certificate{int, leaf, root})
			Expect(err).ToNot(HaveOccurred())

			foundRootCa, err := findCertByName(certs, "root-ca")
			Expect(err).ToNot(HaveOccurred())
			Expect(foundRootCa.Value.CaName).To(BeEmpty())

			foundIntCa, err := findCertByName(certs, "int-ca")
			Expect(err).ToNot(HaveOccurred())
			Expect(foundIntCa.Value.CaName).To(Equal("root-ca"))
			Expect(foundIntCa.Value.Ca).To(BeEmpty())

			foundLeafCert, err := findCertByName(certs, "leaf-cert")
			Expect(err).ToNot(HaveOccurred())
			Expect(foundLeafCert.Value.CaName).To(Equal("int-ca"))
			Expect(foundLeafCert.Value.Ca).To(BeEmpty())
		})
	})
})

func findCertByName(certs []credentials2.Certificate, name string) (credentials2.Certificate, error) {
	for _, cert := range certs {
		if cert.Name == name {
			return cert, nil
		}
	}
	return credentials2.Certificate{}, errors.New("Could not find " + name)
}
