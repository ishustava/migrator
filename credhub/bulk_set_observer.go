package credhub

import (
	"io"
	"errors"
	"fmt"
)

//go:generate counterfeiter . BulkSetObserver
type BulkSetObserver interface {
	BeginBulkSet(numPasswords, numCertificates, numRsaKeys, numSshKeys int)
	FailPasswordSet(name string, err error)
	EndPasswordsSet()
	FailCertificateSet(name string, err error)
	EndCertificatesSet()
	FailRsaKeySet(name string, err error)
	EndRsaKeysSet()
	FailSshKeySet(name string, err error)
	EndSshKeysSet()
	EndBulkSet() error
}

const (
	PASSWORD = "password"
	CERTIFICATE = "certificate"
	RSA_KEY = "RSA key"
	SSH_KEY = "SSH key"
)
type bulkSetObserver struct {
	writer io.Writer

	totalPasswords int
	failedPasswords int

	totalCertificates int
	failedCertificates int

	totalRsaKeys int
	failedRsaKeys int

	totalSshKeys int
	failedSshKeys int
}

func NewBulkSetObserver(writer io.Writer) BulkSetObserver {
	return &bulkSetObserver{
		writer: writer,
	}
}
func (bso *bulkSetObserver) BeginBulkSet(numPasswords, numCertificates, numRsaKeys, numSshKeys int) {
	bso.totalPasswords = numPasswords
	bso.totalCertificates = numCertificates
	bso.totalRsaKeys = numRsaKeys
	bso.totalSshKeys = numSshKeys

	total := numPasswords + numCertificates + numRsaKeys + numSshKeys

	fmt.Fprintf(
		bso.writer,
		"Planning to migrate %d %s, %d %s, %d %s, and %d %s (%d %s total).\n",
		numPasswords, pluralizeIfNecessary(numPasswords, PASSWORD),
		numCertificates, pluralizeIfNecessary(numCertificates, CERTIFICATE),
		numRsaKeys, pluralizeIfNecessary(numRsaKeys, RSA_KEY),
		numSshKeys, pluralizeIfNecessary(numSshKeys, SSH_KEY),
		total, pluralizeIfNecessary(total, "credential"),
	)
}

func (bso *bulkSetObserver) FailPasswordSet(name string, err error) {
	bso.failSet(PASSWORD, name, err)
}
func (bso bulkSetObserver) EndPasswordsSet() {
	bso.endSet(PASSWORD)
}
func (bso *bulkSetObserver) FailCertificateSet(name string, err error) {
	bso.failSet(CERTIFICATE, name, err)
}
func (bso bulkSetObserver) EndCertificatesSet() {
	bso.endSet(CERTIFICATE)
}
func (bso *bulkSetObserver) FailRsaKeySet(name string, err error) {
	bso.failSet(RSA_KEY, name, err)
}
func (bso bulkSetObserver) EndRsaKeysSet() {
	bso.endSet(RSA_KEY)
}
func (bso *bulkSetObserver) FailSshKeySet(name string, err error) {
	bso.failSet(SSH_KEY, name, err)
}
func (bso bulkSetObserver) EndSshKeysSet() {
	bso.endSet(SSH_KEY)
}
func (bso bulkSetObserver) EndBulkSet() error {
	total := bso.totalPasswords + bso.totalCertificates + bso.totalRsaKeys + bso.totalSshKeys
	failed := bso.failedPasswords + bso.failedCertificates + bso.failedRsaKeys + bso.failedSshKeys

	fmt.Fprintf(
		bso.writer,
		"Finished migrating credentials: %d succeeded, %d failed.\n",
		total - failed,
		failed,
	)

	if failed != 0 {
		return errors.New("Failed migrating credentials.")
	} else {
		return nil
	}
}

func (bso *bulkSetObserver) failSet(credentialType, name string, err error) {
	switch credentialType {
	case PASSWORD:
		bso.failedPasswords++
	case CERTIFICATE:
		bso.failedCertificates++
	case RSA_KEY:
		bso.failedRsaKeys++
	case SSH_KEY:
		bso.failedSshKeys++
	default:
		panic("should not be called with invalid credential type")
	}


	fmt.Fprintf(
		bso.writer,
		"Failed migrating %s \"%s\".\nError: %s\n",
		credentialType,
		name,
		err.Error(),
	)
}

func (bso bulkSetObserver) endSet(credentialType string) {
	var total, failed int
	switch credentialType {
	case PASSWORD:
		total = bso.totalPasswords
		failed = bso.failedPasswords
	case CERTIFICATE:
		total = bso.totalCertificates
		failed = bso.failedCertificates
	case RSA_KEY:
		total = bso.totalRsaKeys
		failed = bso.failedRsaKeys
	case SSH_KEY:
		total = bso.totalSshKeys
		failed = bso.failedSshKeys
	default:
		panic("should not be called with invalid credential type")
	}

	fmt.Fprintf(
		bso.writer,
		"Finished migrating %ss: %d succeeded, %d failed.\n",
		credentialType,
		total - failed,
		failed,
	)
}

func pluralizeIfNecessary(count int, word string) string {
	if count != 1 {
		return word + "s"
	}
	return word
}

