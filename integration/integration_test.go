package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("Credential Migration", func() {
	It("Can set password credentials from a vars file into Credhub", func() {
		session := RunCommand("migrate",
			"-v", "../test_fixtures/valid_creds.yml",
			"-e", "my-bosh",
			"-d", "my-deployment")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say("Successfully migrated all credentials"))

	})
})
