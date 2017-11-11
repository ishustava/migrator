package integration

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"os/exec"
	"path"
	"gopkg.in/yaml.v2"
)

var (
	homeDir     string
	TestConfig  Config
	CommandPath string
)

type Config struct {
	ApiUrl       string `yaml:"credhub_api"`
	CredhubCA    []string `yaml:"credhub_ca_cert"`
	ClientName   string `yaml:"credhub_client"`
	ClientSecret string `yaml:"credhub_secret"`
}

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeEach(func() {
	var err error
	homeDir, err = ioutil.TempDir("", "migrator-test")
	Expect(err).NotTo(HaveOccurred())

	os.Setenv("HOME", homeDir)

	TestConfig, err = LoadConfig()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterEach(func() {
	os.RemoveAll(homeDir)
})

var _ = SynchronizedBeforeSuite(func() []byte {
	path, err := Build("github.com/ishustava/migrator")
	Expect(err).NotTo(HaveOccurred())

	return []byte(path)
}, func(data []byte) {
	CommandPath = string(data)
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	CleanupBuildArtifacts()
})

func RunCommand(args ...string) *Session {
	cmd := exec.Command(CommandPath, args...)
	if os.Getenv("CONFIG") != "" {
		existing := os.Environ()
		existing = append(existing, "CREDHUB_URL="+string(TestConfig.ApiUrl))
		existing = append(existing, "CREDHUB_CLIENT="+string(TestConfig.ClientName))
		existing = append(existing, "CREDHUB_SECRET="+string(TestConfig.ClientSecret))
		for _, cert := range TestConfig.CredhubCA {
			existing = append(existing, "CREDHUB_CA_CERT="+string(cert))
		}
		cmd.Env = existing
	}
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}

func LoadConfig() (Config, error) {
	configuration := Config{}
	configFilePath := os.Getenv("CONFIG")

	configurationYaml, err := ioutil.ReadFile(path.Join(os.Getenv("PWD"), configFilePath))
	if err != nil {
		return configuration, err
	}

	err = yaml.Unmarshal(configurationYaml, &configuration)
	if err != nil {
		return configuration, err
	}

	return configuration, nil
}
