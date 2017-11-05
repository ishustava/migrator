package integration

import (
	"io/ioutil"
	"os"
	"testing"
	"github.com/ishustava/migrator/test_fixtures"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var (
	homeDir string
	cfg     test_fixtures.Config
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeEach(func() {
	var err error
	homeDir, err = ioutil.TempDir("", "migrator-test")
	Expect(err).NotTo(HaveOccurred())

	os.Setenv("HOME", homeDir)

	cfg, err = test_fixtures.LoadConfig()
	Expect(err).NotTo(HaveOccurred())

	// These happen before each test due to the lack of a BeforeAll
	// (https://github.com/onsi/ginkgo/issues/70) :(
	// If the tests are slow, they should be runnable in parallel with the -p option.
	test_fixtures.TargetAndLogin(cfg)
})

var _ = AfterEach(func() {
	os.RemoveAll(homeDir)
})

var _ = SynchronizedBeforeSuite(func() []byte {
	path, err := Build("github.com/ishustava/migrator")
	Expect(err).NotTo(HaveOccurred())

	return []byte(path)
}, func(data []byte) {
	test_fixtures.CommandPath = string(data)
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	CleanupBuildArtifacts()
})
