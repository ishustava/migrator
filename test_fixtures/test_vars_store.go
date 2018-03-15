package test_fixtures

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

func GenerateTestVarsStore() (string, error) {
	tempVarsStoreFile, err := ioutil.TempFile("", "vars-store")
	if err != nil {
		return "", err
	}

	tempVarsStore := map[string]interface{}{
		"path1": PASSWORD1,
		"path2": PASSWORD2,
		"path3": map[string]string{
			"ca": SIGNED_BY_ROOT_LEAF1_CA,
			"certificate": SIGNED_BY_ROOT_LEAF1_CERT,
			"private_key": SIGNED_BY_ROOT_LEAF1_PRIV,
		},
		"path4": map[string]string{
			"ca": ROOT_CA_CERT,
			"certificate": ROOT_CA_CERT,
			"private_key": ROOT_CA_PRIV,
		},
		"path5": map[string]string{
			"public_key": SSH_PUB,
			"private_key": SSH_PRIV,
			"public_key_fingerprint": SSH_FINGERPRINT,
		},
		"path6": map[string]string{
			"public_key": RSA_PUB,
			"private_key": RSA_PRIV,
		},
	}

	marshaledVarsStore, err := yaml.Marshal(tempVarsStore)
	if err != nil {
		return "", err
	}

	tempVarsStoreFile.Write(marshaledVarsStore)

	return tempVarsStoreFile.Name(), nil
}