package migrator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ishustava/migrator/credhub/credhubfakes"
	"github.com/ishustava/migrator/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	"github.com/ishustava/migrator/migrator"
	credentials2 "github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
)

var _ = Describe("Migrator", func() {
	It("sets credentials in CredHub", func() {
		client := new(credhubfakes.FakeCredHubClient)
		firstPass := credentials.NewPassword("test-pass-1", "value1")
		secondPass := credentials.NewPassword("test-pass-2", "value2")
		cert := credentials.NewCertificate("test-cert", values.Certificate{})
		ssh := credentials.NewSsh("test-ssh", values.SSH{})
		rsa := credentials.NewRsa("test-rsa", values.RSA{})
		creds := credentials.Credentials{
			Passwords: []credentials2.Password{firstPass, secondPass},
			Certificates: []credentials2.Certificate{cert},
			SshKeys: []credentials2.SSH{ssh},
			RsaKeys: []credentials2.RSA{rsa},
		}

		err := migrator.Migrate(creds, client)

		Expect(err).ToNot(HaveOccurred())
		Expect(client.SetPasswordCallCount()).To(Equal(2))
		Expect(client.SetCertificateCallCount()).To(Equal(1))
		Expect(client.SetRSACallCount()).To(Equal(1))
		Expect(client.SetSSHCallCount()).To(Equal(1))
	})
})