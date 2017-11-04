package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/ghttp"
	"net/http"
	"fmt"
	"github.com/ishustava/migrator/test_fixtures"
	"encoding/json"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

const PUT_REQUEST = `{"type":"%s","name":"%s","value":%s,"mode":"overwrite"}`

var _ = Describe("Migrate", func() {
	var (
		credhubServer *ghttp.Server
		uaaServer     *ghttp.Server
	)

	BeforeEach(func() {
		credhubServer = newTLSServer("test-certs/credhub-tls-cert.pem", "test-certs/credhub-tls-key.pem")
		uaaServer = newTLSServer("test-certs/uaa-tls-cert.pem", "test-certs/uaa-tls-key.pem")

		setupInfoHandler(credhubServer, uaaServer)
		setupLoginHandler(uaaServer)
	})

	AfterEach(func() {
		credhubServer.Reset()
		uaaServer.Reset()
	})

	Context("Successful", func() {
		BeforeEach(func() {
			setupPasswordHandler(credhubServer)
			setupCertificateHandler(credhubServer)
			setupRSAHandler(credhubServer)
			setupSSHHandler(credhubServer)
		})

		It("migrates credentials from vars store to credhub", func() {
			session := runCommand("migrate",
				"-v", "../test_fixtures/valid_creds.yml",
				"-u", credhubServer.URL(),
				"-c", "test_client",
				"-s", "test_secret",
				"--ca-cert", "test-certs/credhub-tls-ca.pem",
				"--ca-cert", "test-certs/uaa-tls-ca.pem")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("Successfully migrated all credentials"))

			Expect(credhubServer.ReceivedRequests()).To(HaveLen(8))
			Expect(uaaServer.ReceivedRequests()).To(HaveLen(1))
		})
	})

})

func setupRSAHandler(credhubServer *ghttp.Server) {
	rsaJson, err := json.Marshal(values.SSH{PublicKey: test_fixtures.RSA_PUB, PrivateKey: test_fixtures.RSA_PRIV})
	Expect(err).ToNot(HaveOccurred())

	credhubServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", "/api/v1/data"),
			ghttp.VerifyJSON(fmt.Sprintf(PUT_REQUEST, "rsa", "path6", rsaJson)),
			ghttp.RespondWith(http.StatusOK, nil),
		),
	)
}

func setupSSHHandler(credhubServer *ghttp.Server) {
	sshJson, err := json.Marshal(values.SSH{PublicKey: test_fixtures.SSH_PUB, PrivateKey: test_fixtures.SSH_PRIV})
	Expect(err).ToNot(HaveOccurred())

	credhubServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", "/api/v1/data"),
			ghttp.VerifyJSON(fmt.Sprintf(PUT_REQUEST, "ssh", "path5", sshJson)),
			ghttp.RespondWith(http.StatusOK, nil),
		),
	)
}

func setupPasswordHandler(credhubServer *ghttp.Server) {
	credhubServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", "/api/v1/data"),
			ghttp.VerifyJSON(fmt.Sprintf(PUT_REQUEST, "password", "path1", `"password1"`)),
			ghttp.RespondWith(http.StatusOK, nil),
		),
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", "/api/v1/data"),
			ghttp.VerifyJSON(fmt.Sprintf(PUT_REQUEST, "password", "path2", `"password2"`)),
			ghttp.RespondWith(http.StatusOK, nil),
		),
	)
}

func setupCertificateHandler(credhubServer *ghttp.Server) {
	certJson1, err := json.Marshal(values.Certificate{Ca: test_fixtures.SIGNED_BY_ROOT_LEAF1_CA, Certificate: test_fixtures.SIGNED_BY_ROOT_LEAF1_CERT, PrivateKey: test_fixtures.SIGNED_BY_ROOT_LEAF1_PRIV})
	Expect(err).ToNot(HaveOccurred())

	certJson2, err := json.Marshal(values.Certificate{Certificate: test_fixtures.ROOT_CA_CERT, PrivateKey: test_fixtures.ROOT_CA_PRIV})
	Expect(err).ToNot(HaveOccurred())

	credhubServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", "/api/v1/data"),
			ghttp.VerifyJSON(fmt.Sprintf(PUT_REQUEST, "certificate", "path3", certJson1)),
			ghttp.RespondWith(http.StatusOK, nil),
		),
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", "/api/v1/data"),
			ghttp.VerifyJSON(fmt.Sprintf(PUT_REQUEST, "certificate", "path4", certJson2)),
			ghttp.RespondWith(http.StatusOK, nil),
		),
	)
}

func setupLoginHandler(uaaServer *ghttp.Server) {
	uaaServer.RouteToHandler("POST", "/oauth/token",
		ghttp.CombineHandlers(
			ghttp.VerifyBody([]byte(`client_id=test_client&client_secret=test_secret&grant_type=client_credentials&response_type=token`)),
			ghttp.RespondWith(http.StatusOK, `{
						"access_token":"2YotnFZFEjr1zCsicMWpAA",
						"refresh_token":"erousflkajqwer",
						"token_type":"bearer",
						"expires_in":3600}`),
		),
	)
}

func setupInfoHandler(credhubServer, uaaServer *ghttp.Server) {
	credhubServer.RouteToHandler("GET", "/info",
		ghttp.RespondWith(http.StatusOK, `{
				"app":{"version":"9.9.9","name":"CredHub"},
				"auth-server":{"url":"`+ uaaServer.URL()+ `"}
				}`),
	)
}
