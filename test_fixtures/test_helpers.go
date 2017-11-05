package test_fixtures

import (
	"io/ioutil"
	"path"
	"os"
	"encoding/json"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

type Config struct {
	ApiUrl         string `json:"credhub_api"`
	CredhubCA      string `json:"credhub_ca_cert"`
	ClientName     string `json:"credhub_client"`
	ClientSecret   string `json:"credhub_secret"`
}

var (
	CommandPath string
)

func RunCommand(args ...string) *Session {
	cmd := exec.Command(CommandPath, args...)

	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}

func LoadConfig() (Config, error) {

	configuration := Config{}
	configFilePath := os.Getenv("CONFIG")
	configurationJson, err := ioutil.ReadFile(path.Join(os.Getenv("PWD"), configFilePath))
	if err != nil {
		return configuration, err
	}

	err = json.Unmarshal(configurationJson, &configuration)
	if err != nil {
		return configuration, err
	}

	return configuration, nil
}

func TargetAndLogin(cfg Config) {
	CleanEnv()
	os.Setenv("CREDHUB_CLIENT", cfg.ClientName)
	os.Setenv("CREDHUB_SECRET", cfg.ClientSecret)
	session := RunCommand("login", "-s", cfg.ApiUrl, "--ca-cert", cfg.CredhubCA)
	Eventually(session).Should(Exit(0))
}

func CleanEnv() {
	os.Unsetenv("CREDHUB_SECRET")
	os.Unsetenv("CREDHUB_CLIENT")
}
