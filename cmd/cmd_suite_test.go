package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"

	"testing"
	"os/exec"
	"crypto/tls"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

var (
	commandPath string
)

var _ = SynchronizedBeforeSuite(func() []byte {
	executable_path, err := Build("github.com/ishustava/migrator")
	Expect(err).NotTo(HaveOccurred())
	return []byte(executable_path)
}, func(data []byte) {
	commandPath = string(data)
})

func runCommand(args ...string) *Session {
	cmd := exec.Command(commandPath, args...)
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}

func newTLSServer(certPath, keyPath string) *Server {
	server := NewUnstartedServer()

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	Expect(err).ToNot(HaveOccurred())

	server.HTTPTestServer.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	server.HTTPTestServer.StartTLS()

	return server
}
