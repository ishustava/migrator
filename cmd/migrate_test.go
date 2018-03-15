package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/ghttp"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	"github.com/ishustava/migrator/test_fixtures"
	"os"
)

var _ = Describe("Migrate", func() {
	var (
		credhubServer *ghttp.Server
		uaaServer     *ghttp.Server
		requestBodies []string
		varsStoreFile string
		err           error
	)

	var saveRequestBody = func(_ http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		req.Body.Close()
		Expect(err).ShouldNot(HaveOccurred())
		requestBodies = append(requestBodies, string(body))
	}

	BeforeEach(func() {
		credhubServer = newTLSServer("test-certs/credhub-tls-cert.pem", "test-certs/credhub-tls-key.pem")
		uaaServer = newTLSServer("test-certs/uaa-tls-cert.pem", "test-certs/uaa-tls-key.pem")

		setupInfoHandler(credhubServer, uaaServer)
		setupLoginHandler(uaaServer)

		varsStoreFile, err = test_fixtures.GenerateTestVarsStore()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		credhubServer.Reset()
		uaaServer.Reset()
		Expect(os.Remove(varsStoreFile)).ToNot(HaveOccurred())
	})

	Context("Successful", func() {
		BeforeEach(func() {
			credhubServer.RouteToHandler("PUT", "/api/v1/data", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, nil),
				saveRequestBody,
			))
		})

		It("migrates credentials from vars store to credhub", func() {
			session := runCommand("migrate",
				"-v", varsStoreFile,
				"-u", credhubServer.URL(),
				"-c", "test_client",
				"-s", "test_secret",
				"-e", "my-bosh",
				"-d", "my-deployment",
				"--ca-cert", "test-certs/credhub-tls-ca.pem",
				"--ca-cert", "test-certs/uaa-tls-ca.pem")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("Successfully migrated all credentials"))

			rsaJson, err := json.Marshal(values.SSH{PublicKey: test_fixtures.RSA_PUB, PrivateKey: test_fixtures.RSA_PRIV})
			Expect(err).ToNot(HaveOccurred())

			sshJson, err := json.Marshal(values.SSH{PublicKey: test_fixtures.SSH_PUB, PrivateKey: test_fixtures.SSH_PRIV})
			Expect(err).ToNot(HaveOccurred())

			certJson1, err := json.Marshal(values.Certificate{CaName: "my-bosh/my-deployment/path4", Certificate: test_fixtures.SIGNED_BY_ROOT_LEAF1_CERT, PrivateKey: test_fixtures.SIGNED_BY_ROOT_LEAF1_PRIV})
			Expect(err).ToNot(HaveOccurred())

			certJson2, err := json.Marshal(values.Certificate{Ca: test_fixtures.ROOT_CA_CERT, Certificate: test_fixtures.ROOT_CA_CERT, PrivateKey: test_fixtures.ROOT_CA_PRIV})
			Expect(err).ToNot(HaveOccurred())

			passJsonRequest1 := putRequestBody("password", "my-bosh/my-deployment/path1", `"password1"`)
			passJsonRequest2 := putRequestBody("password", "my-bosh/my-deployment/path2", `"password2"`)
			certJsonRequest1 := putRequestBody("certificate", "my-bosh/my-deployment/path3", string(certJson1))
			certJsonRequest2 := putRequestBody("certificate", "my-bosh/my-deployment/path4", string(certJson2))
			sshJsonRequest := putRequestBody("ssh", "my-bosh/my-deployment/path5", string(sshJson))
			rsaJsonRequest := putRequestBody("rsa", "my-bosh/my-deployment/path6", string(rsaJson))

			Expect(uaaServer.ReceivedRequests()).To(HaveLen(1))
			Expect(requestBodies).Should(ConsistOf(
				MatchJSON(passJsonRequest1),
				MatchJSON(passJsonRequest2),
				MatchJSON(rsaJsonRequest),
				MatchJSON(sshJsonRequest),
				MatchJSON(certJsonRequest1),
				MatchJSON(certJsonRequest2),
			))
		})
	})
})

func putRequestBody(t, name, value string) string {
	return fmt.Sprintf(`{"type":"%s","name":"%s","value":%s,"mode":"overwrite"}`, t, name, value)
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
