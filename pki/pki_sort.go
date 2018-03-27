package pki

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"crypto/x509"
	"encoding/pem"
	"sort"
)

type pki []*pkiElement

type pkiElement struct {
	certificate credentials.Certificate
	signer      *pkiElement
}

func Sort(certificates []credentials.Certificate) {
	pkiElements := make(pki, len(certificates))

	for i, cert := range certificates {
		pkiElements[i] = &pkiElement{
			certificate: cert,
			signer:      findSigner(cert, remove(certificates, i)),
		}
	}

	sort.Sort(pkiElements)

	for i, pe := range pkiElements {
		certificates[i] = pe.certificate
		if pe.signer != nil {
			certificates[i].Value.CaName = pe.signer.certificate.Name
			certificates[i].Value.Ca = ""
		}
	}
}

func (pki pki) Len() int           { return len(pki) }
func (pki pki) Swap(i, j int)      { pki[i], pki[j] = pki[j], pki[i] }
func (pki pki) Less(i, j int) bool { return depth(pki[i]) < depth(pki[j]) }

func depth(pe *pkiElement) int {
	if pe.signer == nil {
		return 0
	} else {
		return 1 + depth(pe.signer)
	}
}

func remove(certificates []credentials.Certificate, i int) []credentials.Certificate {
	result := make([]credentials.Certificate, i)
	copy(result, certificates[:i])

	return append(result, certificates[i+1:]...)
}

func findSigner(certificate credentials.Certificate, possibleCAs []credentials.Certificate) *pkiElement {
	parsedCertificate := parsePemCertificate(certificate.Value.Certificate)
	for i, possibleCA := range possibleCAs {
		parsedPossibleCA := parsePemCertificate(possibleCA.Value.Certificate)
		roots := x509.NewCertPool()
		roots.AddCert(parsedPossibleCA)
		if _, err := parsedCertificate.Verify(x509.VerifyOptions{Roots: roots, KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny}}); err == nil {
			return &pkiElement{
				certificate: possibleCA,
				signer:      findSigner(possibleCA, remove(possibleCAs, i)),
			}
		}
	}

	return nil
}

func parsePemCertificate(pemCert string) (*x509.Certificate) {
	block, _ := pem.Decode([]byte(pemCert))
	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("could not parse pem certificate")
	}
	return parsedCert
}
