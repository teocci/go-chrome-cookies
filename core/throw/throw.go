// Package throw
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-15
package throw

import "errors"

const (
	errItemNotSupported    = `item not supported, default is "all", choose from history|downloads|password|bookmark|cookie`
	errBrowserNotSupported = "browser not supported"
	errChromeSecretIsEmpty = "chrome secret is empty"
	errDbusSecretIsEmpty   = "dbus secret key is empty"

	errSecurityKeyIsEmpty = "input [security find-generic-password -wa 'Chrome'] in terminal"
	errPasswordIsEmpty    = "password is empty"
	errDecryptFailed      = "decrypt failed, password is empty"
	errDecodeASN1Failed   = "decode ASN1 data failed"
)

func ErrorItemNotSupported() error {
	return errors.New(errItemNotSupported)
}

func ErrorBrowserNotSupported() error {
	return errors.New(errBrowserNotSupported)
}

func ErrorChromeSecretIsEmpty() error {
	return errors.New(errChromeSecretIsEmpty)
}

func ErrorDbusSecretIsEmpty() error {
	return errors.New(errDbusSecretIsEmpty)
}

func ErrorSecurityKeyIsEmpty() error {
	return errors.New(errSecurityKeyIsEmpty)
}

func ErrorPasswordIsEmpty() error {
	return errors.New(errPasswordIsEmpty)
}

func ErrorDecryptFailed() error {
	return errors.New(errDecryptFailed)
}

func ErrorDecodeASN1Failed() error {
	return errors.New(errDecodeASN1Failed)
}