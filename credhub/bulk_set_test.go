package credhub_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ishustava/migrator/credhub/credhubfakes"
	"github.com/ishustava/migrator/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	credentials2 "github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/ishustava/migrator/credhub"
	credhub2 "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"errors"
)

var _ = Describe("Bulk Set", func() {
	var creds credentials.Credentials
	var credhubClient *credhubfakes.FakeCredHubClient
	var observer *credhubfakes.FakeBulkSetObserver

	BeforeEach(func() {
		goodPass := credentials.NewPassword("good-password", "value1")
		badPass := credentials.NewPassword("bad-password", "value2")
		goodCert := credentials.NewCertificate("good-certificate", values.Certificate{})
		badCert := credentials.NewCertificate("bad-certificate", values.Certificate{})
		goodSsh := credentials.NewSsh("good-ssh-key", values.SSH{})
		badSsh := credentials.NewSsh("bad-ssh-key", values.SSH{})
		goodRsa := credentials.NewRsa("good-rsa-key", values.RSA{})
		badRsa := credentials.NewRsa("bad-rsa-key", values.RSA{})
		creds = credentials.Credentials{
			Passwords: []credentials2.Password{goodPass, badPass},
			Certificates: []credentials2.Certificate{goodCert, badCert},
			SshKeys: []credentials2.SSH{goodSsh, badSsh},
			RsaKeys: []credentials2.RSA{goodRsa, badRsa},
		}

		credhubClient = new(credhubfakes.FakeCredHubClient)
		observer = new(credhubfakes.FakeBulkSetObserver)
	})

	Describe("interacting with CredHub", func() {
		// TODO: test credhubClient args for calls
		It("sets credentials in CredHub", func() {
			credhub.BulkSet(&creds, credhubClient, observer)

			Expect(credhubClient.SetPasswordCallCount()).To(Equal(2))
			Expect(credhubClient.SetCertificateCallCount()).To(Equal(2))
			Expect(credhubClient.SetRSACallCount()).To(Equal(2))
			Expect(credhubClient.SetSSHCallCount()).To(Equal(2))
		})
	})

	Describe("interacting with the observer", func() {
		var passErr, certErr, rsaKeyErr, sshKeyErr error

		BeforeEach(func() {
			credhubClient.SetPasswordStub = func(name string, value values.Password, mode credhub2.Mode) (credentials2.Password, error) {
				if name == "bad-password" {
					passErr = errors.New("password error from credhub")
					return credentials2.Password{}, passErr
				}
				return credentials2.Password{}, nil
			}
			credhubClient.SetCertificateStub = func(name string, value values.Certificate, mode credhub2.Mode) (credentials2.Certificate, error) {
				if name == "bad-certificate" {
					certErr = errors.New("certificate error from credhub")
					return credentials2.Certificate{}, certErr
				}
				return credentials2.Certificate{}, nil
			}
			credhubClient.SetRSAStub = func(name string, value values.RSA, mode credhub2.Mode) (credentials2.RSA, error) {
				if name == "bad-rsa-key" {
					rsaKeyErr = errors.New("rsa key error from credhub")
					return credentials2.RSA{}, rsaKeyErr
				}
				return credentials2.RSA{}, nil
			}
			credhubClient.SetSSHStub = func(name string, value values.SSH, mode credhub2.Mode) (credentials2.SSH, error) {
				if name == "bad-ssh-key" {
					sshKeyErr = errors.New("ssh key error from credhub")
					return credentials2.SSH{}, sshKeyErr
				}
				return credentials2.SSH{}, nil
			}
		})

		It("informs the observer", func() {
			credhub.BulkSet(&creds, credhubClient, observer)

			Expect(observer.BeginBulkSetCallCount()).To(Equal(1))
			numPass, numCert, numRsa, numSsh := observer.BeginBulkSetArgsForCall(0)
			Expect(numPass).To(Equal(2))
			Expect(numCert).To(Equal(2))
			Expect(numRsa).To(Equal(2))
			Expect(numSsh).To(Equal(2))

			Expect(observer.FailPasswordSetCallCount()).To(Equal(1))
			passName, err := observer.FailPasswordSetArgsForCall(0)
			Expect(passName).To(Equal("bad-password"))
			Expect(err).To(MatchError(passErr))

			Expect(observer.EndPasswordsSetCallCount()).To(Equal(1))

			Expect(observer.FailCertificateSetCallCount()).To(Equal(1))
			certName, err := observer.FailCertificateSetArgsForCall(0)
			Expect(certName).To(Equal("bad-certificate"))
			Expect(err).To(MatchError(certErr))

			Expect(observer.FailCertificateSetCallCount()).To(Equal(1))

			Expect(observer.FailRsaKeySetCallCount()).To(Equal(1))
			rsaKeyName, err := observer.FailRsaKeySetArgsForCall(0)
			Expect(rsaKeyName).To(Equal("bad-rsa-key"))
			Expect(err).To(MatchError(rsaKeyErr))

			Expect(observer.FailRsaKeySetCallCount()).To(Equal(1))

			Expect(observer.FailSshKeySetCallCount()).To(Equal(1))
			sshKeyName, err := observer.FailSshKeySetArgsForCall(0)
			Expect(sshKeyName).To(Equal("bad-ssh-key"))
			Expect(err).To(MatchError(sshKeyErr))

			Expect(observer.FailSshKeySetCallCount()).To(Equal(1))

			Expect(observer.EndBulkSetCallCount()).To(Equal(1))
		})
	})
})