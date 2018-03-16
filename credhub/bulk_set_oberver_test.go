package credhub_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/ishustava/migrator/credhub"

	"errors"
)

var _ = Describe("Bulk Set Observer", func() {
	var observer credhub.BulkSetObserver
	var buffer *gbytes.Buffer

	BeforeEach(func() {
		buffer = gbytes.NewBuffer()
		observer = credhub.NewBulkSetObserver(buffer)
	})

	Describe("BeginBulkSet", func() {
		It("writes a plan message", func() {
			observer.BeginBulkSet(1, 2, 3, 4)

			Expect(buffer).To(gbytes.Say("Planning to migrate 1 password, 2 certificates, 3 RSA keys, and 4 SSH keys \\(10 credentials total\\).\n"))
		})
	})

	Describe("FailPasswordSet", func() {
		It("writes a failure message", func() {
			observer.FailPasswordSet("/path/to/password", errors.New("bad-password"))

			Expect(buffer).To(gbytes.Say("Failed migrating password \"/path/to/password\".\nError: bad-password\n"))
		})
	})

	Describe("EndPasswordsSet", func() {
		Context("when bulk setting has begun", func() {
			BeforeEach(func() {
				observer.BeginBulkSet(5, 0, 0, 0)
			})

			Context("when no passwords fail to set", func() {
				It("summarizes the password migrations", func() {
					observer.EndPasswordsSet()

					Expect(buffer).To(gbytes.Say("Finished migrating passwords: 5 succeeded, 0 failed.\n"))
				})
			})

			Context("when some passwords fail to set", func() {
				BeforeEach(func() {
					observer.FailPasswordSet("/path/password-1", errors.New("bad-password"))
					observer.FailPasswordSet("/path/password-2", errors.New("bad-password"))
				})

				It("summarizes the password migrations", func() {
					observer.EndPasswordsSet()

					Expect(buffer).To(gbytes.Say("Finished migrating passwords: 3 succeeded, 2 failed.\n"))
				})
			})
		})
	})

	Describe("FailCertificateSet", func() {
		It("writes a failure message", func() {
			observer.FailCertificateSet("/path/to/certificate", errors.New("bad-certificate"))

			Expect(buffer).To(gbytes.Say("Failed migrating certificate \"/path/to/certificate\".\nError: bad-certificate\n"))
		})
	})

	Describe("EndCertificatesSet", func() {
		Context("when bulk setting has begun", func() {
			BeforeEach(func() {
				observer.BeginBulkSet(0, 5, 0, 0)
			})

			Context("when no certificates fail to set", func() {
				It("summarizes the certificates migrations", func() {
					observer.EndCertificatesSet()

					Expect(buffer).To(gbytes.Say("Finished migrating certificates: 5 succeeded, 0 failed.\n"))
				})
			})

			Context("when some passwords fail to set", func() {
				BeforeEach(func() {
					observer.FailCertificateSet("/path/certificate-1", errors.New("bad-certificate"))
					observer.FailCertificateSet("/path/certificate-2", errors.New("bad-certificate"))
				})

				It("summarizes the certificate migrations", func() {
					observer.EndCertificatesSet()

					Expect(buffer).To(gbytes.Say("Finished migrating certificates: 3 succeeded, 2 failed.\n"))
				})
			})
		})
	})

	Describe("FailRsaKeySet", func() {
		It("writes a failure message", func() {
			observer.FailRsaKeySet("/path/to/rsa-key", errors.New("bad-rsa-key"))

			Expect(buffer).To(gbytes.Say("Failed migrating RSA key \"/path/to/rsa-key\".\nError: bad-rsa-key\n"))
		})
	})

	Describe("EndRsaKeysSet", func() {
		Context("when bulk setting has begun", func() {
			BeforeEach(func() {
				observer.BeginBulkSet(0, 0, 5, 0)
			})

			Context("when no RSA keys fail to set", func() {
				It("summarizes the RSA key migrations", func() {
					observer.EndRsaKeysSet()

					Expect(buffer).To(gbytes.Say("Finished migrating RSA keys: 5 succeeded, 0 failed.\n"))
				})
			})

			Context("when some RSA keys fail to set", func() {
				BeforeEach(func() {
					observer.FailRsaKeySet("/path/rsa-key-1", errors.New("bad-rsa-key"))
					observer.FailRsaKeySet("/path/rsa-key-2", errors.New("bad-rsa-key"))
				})

				It("summarizes the RSA key migrations", func() {
					observer.EndRsaKeysSet()

					Expect(buffer).To(gbytes.Say("Finished migrating RSA keys: 3 succeeded, 2 failed.\n"))
				})
			})
		})
	})

	Describe("FailSshKeySet", func() {
		It("writes a failure message", func() {
			observer.FailSshKeySet("/path/to/ssh-key", errors.New("bad-ssh-key"))

			Expect(buffer).To(gbytes.Say("Failed migrating SSH key \"/path/to/ssh-key\".\nError: bad-ssh-key\n"))
		})
	})

	Describe("EndSshKeysSet", func() {
		Context("when bulk setting has begun", func() {
			BeforeEach(func() {
				observer.BeginBulkSet(0, 0, 0, 5)
			})

			Context("when no SSH keys fail to set", func() {
				It("summarizes the SSH key migrations", func() {
					observer.EndSshKeysSet()

					Expect(buffer).To(gbytes.Say("Finished migrating SSH keys: 5 succeeded, 0 failed.\n"))
				})
			})

			Context("when some SSH keys fail to set", func() {
				BeforeEach(func() {
					observer.FailSshKeySet("/path/ssh-key-1", errors.New("bad-ssh-key"))
					observer.FailSshKeySet("/path/ssh-key-2", errors.New("bad-ssh-key"))
				})

				It("summarizes the SSH key migrations", func() {
					observer.EndSshKeysSet()

					Expect(buffer).To(gbytes.Say("Finished migrating SSH keys: 3 succeeded, 2 failed.\n"))
				})
			})
		})
	})

	Describe("EndBulkSet", func() {
		Context("when bulk setting has begun", func() {
			BeforeEach(func() {
				observer.BeginBulkSet(2, 4, 6, 8)
			})

			Context("when no migrations fail", func() {
				It("summarizes all migrations", func() {
					observer.EndBulkSet()

					Expect(buffer).To(gbytes.Say("Finished migrating credentials: 20 succeeded, 0 failed.\n"))
				})

				It("does not error", func() {
					Expect(observer.EndBulkSet()).ToNot(HaveOccurred())
				})
			})

			Context("when some migrations fail", func() {
				BeforeEach(func() {
					observer.FailPasswordSet("/path/password-1", errors.New("bad-password"))
					observer.FailPasswordSet("/path/password-2", errors.New("bad-password"))

					observer.FailCertificateSet("/path/certificate-1", errors.New("bad-certificate"))

					observer.FailSshKeySet("/path/ssh-key-1", errors.New("bad-ssh-key"))
					observer.FailSshKeySet("/path/ssh-key-2", errors.New("bad-ssh-key"))
				})

				It("summarizes all migrations", func() {
					observer.EndBulkSet()

					Expect(buffer).To(gbytes.Say("Finished migrating credentials: 15 succeeded, 5 failed.\n"))
				})

				It("errors", func() {
					Expect(observer.EndBulkSet()).To(HaveOccurred())
				})

			})
		})
	})
})
