package parser

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"encoding/pem"
	"crypto/x509"
	"errors"
	"reflect"
	"fmt"
)

func FindAndSetSigningCA(certificates []credentials.Certificate) ([]credentials.Certificate, error) {
	cas, certs, err := separateCAsFromCerts(certificates)
	if err != nil {
		return certificates, err
	}

	for i, cert := range certs {
		caName, err := findSigningCaName(cert, cas)
		if err != nil {
			return certificates, err
		}
		certs[i].Value.CaName = caName
		certs[i].Value.Ca = "" // CredHub server requires Ca value to be empty if ca name is set
	}

	result := []credentials.Certificate{}
	result = append(result, cas...)
	result = append(result, certs...)
	return result, nil
}

func separateCAsFromCerts(certificates []credentials.Certificate) ([]credentials.Certificate, []credentials.Certificate, error) {
	cas := []credentials.Certificate{}
	certs := []credentials.Certificate{}
	for _, cert := range certificates {
		parsedCert, err := parsePemCertificate(cert.Value.Certificate)
		if err != nil {
			return cas, certs, err
		}

		if parsedCert.IsCA {
			cas = append(cas, cert)
			if isIntermediate(parsedCert) {
				certs = append(certs, cert)
			}
		} else {
			certs = append(certs, cert)
		}
	}
	return cas, certs, nil
}

func isIntermediate(certificate *x509.Certificate) bool {
	return certificate.IsCA && !reflect.DeepEqual(certificate.Subject, certificate.Issuer)
}

func findSigningCaName(certificate credentials.Certificate, cas []credentials.Certificate) (string, error) {
	parsedCert, err := parsePemCertificate(certificate.Value.Certificate)
	if err != nil {
		return "", err
	}
	for _, ca := range cas {
		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM([]byte(ca.Value.Certificate)); ok {
			_, err := parsedCert.Verify(x509.VerifyOptions{Roots: certPool, KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny}})
			if err == nil {
				return ca.Name, nil
			}
		}
	}
	return "", errors.New(fmt.Sprintf("Could not find signing CA for '%s'", certificate.Name))
}

func parsePemCertificate(pemCert string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(pemCert))
	parsedCert, err := x509.ParseCertificate(block.Bytes)
	return parsedCert, err
}