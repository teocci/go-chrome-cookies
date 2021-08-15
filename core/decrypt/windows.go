// Package decrypt
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-12
// +build !linux,!plan9,!nacl,!darwin

package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/teocci/go-chrome-cookies/core/throw"
	"syscall"
	"unsafe"
)

func ChromePass(key, encryptPass []byte) ([]byte, error) {
	if len(encryptPass) > 15 {
		// remove Prefix 'v10'
		return aesGCMDecrypt(encryptPass[15:], key, encryptPass[3:15])
	} else {
		return nil, throw.ErrorPasswordIsEmpty()
	}
}

// aesGCMDecrypt
// chromium > 80
// more info here: https://tinyurl.com/nax82cpn
func aesGCMDecrypt(encrypted, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	origData, err := blockMode.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, err
	}
	return origData, nil
}

type dataBlob struct {
	cbData uint32
	pbData *byte
}

func NewBlob(d []byte) *dataBlob {
	if len(d) == 0 {
		return &dataBlob{}
	}
	return &dataBlob{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

func (b *dataBlob) ToByteArray() []byte {
	d := make([]byte, b.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])
	return d
}

// DPApi
// chrome < 80
// more info here https://tinyurl.com/4zymkmpw
func DPApi(data []byte) ([]byte, error) {
	dllCrypt := syscall.NewLazyDLL("Crypt32.dll")
	dllKernel := syscall.NewLazyDLL("Kernel32.dll")
	procDecryptData := dllCrypt.NewProc("CryptUnprotectData")
	procLocalFree := dllKernel.NewProc("LocalFree")
	var outBlob dataBlob
	r, _, err := procDecryptData.Call(uintptr(unsafe.Pointer(NewBlob(data))), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&outBlob)))
	if r == 0 {
		return nil, err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(outBlob.pbData)))
	return outBlob.ToByteArray(), nil
}