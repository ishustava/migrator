package parser

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"encoding/pem"
	"crypto/x509"
	"reflect"
)

func FindAndSetSigningCA(certificates []credentials.Certificate) ([]credentials.Certificate, error) {
	result := make([]credentials.Certificate, len(certificates))
	copy(result, certificates)

	roots, ints, leafs, err := separateCerts(certificates)
	if err != nil {
		return nil, err
	}

	allCas := append(roots, ints...)

	err = findAndSetSigningCaNames(ints, allCas)
	if err != nil {
		return nil, err
	}

	err = findAndSetSigningCaNames(leafs, allCas)
	if err != nil {
		return nil, err
	}

	result = append(roots, ints...)

	return append(result, leafs...), nil
}

func isRootCa(cert *x509.Certificate) bool {
	return cert.IsCA && reflect.DeepEqual(cert.Subject, cert.Issuer)
}

func isLeaf(cert *x509.Certificate) bool {
	return !cert.IsCA
}

func separateCerts(certs []credentials.Certificate) ([]credentials.Certificate, []credentials.Certificate, []credentials.Certificate, error) {
	roots := make([]credentials.Certificate, 0)
	ints := make([]credentials.Certificate, 0)
	leafs := make([]credentials.Certificate, 0)

	for _, cert := range certs {
		parsedCert, err := parsePemCertificate(cert.Value.Certificate)
		if err != nil {
			return nil, nil, nil, err
		}

		if isRootCa(parsedCert) {
			roots = append(roots, cert)
			continue
		}

		if isLeaf(parsedCert) {
			leafs = append(leafs, cert)
			continue
		}

		ints = append(ints, cert)
	}
	return roots, ints, leafs, nil
}

func findAndSetSigningCaNames(certificates []credentials.Certificate, cas []credentials.Certificate) error {
	for i, int := range certificates {
		caName, err := findSigningCaName(int, cas)
		if err != nil {
			return err
		}

		certificates[i].Value.CaName = caName
		certificates[i].Value.Ca = "" // CredHub server requires Ca value to be empty if ca name is set
	}

	return nil
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
	return certificate.Name, nil
}

func parsePemCertificate(pemCert string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(pemCert))
	parsedCert, err := x509.ParseCertificate(block.Bytes)
	return parsedCert, err
}