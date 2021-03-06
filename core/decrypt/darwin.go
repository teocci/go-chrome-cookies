// Package decrypt
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-12
// +build !windows,!plan9,!nacl,!linux

package decrypt

func ChromePass(key, encryptPass []byte) ([]byte, error) {
	if len(encryptPass) > 3 {
		if len(key) == 0 {
			return nil, errSecurityKeyIsEmpty
		}
		var chromeIV = []byte{32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32}
		return aes128CBCDecrypt(key, chromeIV, encryptPass[3:])
	} else {
		return nil, throw.ErrorDecryptFailed()
	}
}

func DPApi(data []byte) ([]byte, error) {
	return nil, nil
}