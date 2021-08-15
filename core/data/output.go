// Package data
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-12
package data

import (
	"bytes"
	"encoding/json"
	"github.com/teocci/go-chrome-cookies/filemgmt"
	"os"

	"github.com/jszwec/csvutil"
)

type OutputFormat int

const (
	formatJson OutputFormat = iota
	formatCSV
	formatConsole
)

const (
	FormatNameJson    = "json"
	FormatNameCSV     = "csv"
	FormatNameConsole = "console"
)

var (
	utf8Bom = []byte{239, 187, 191}
	formats = map[string]OutputFormat{
		FormatNameJson:    formatJson,
		FormatNameCSV:     formatCSV,
		FormatNameConsole: formatConsole,
	}
)

func GetFormat(formatName string) OutputFormat {
	return formats[formatName]
}

func GetFormatName(format OutputFormat) string {
	return formatNames()[format]
}

func formatNames() []string {
	return []string{FormatNameJson, FormatNameCSV, FormatNameConsole}
}

func WriteToJson(filename string, data interface{}) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	fnc := filemgmt.CloseFile()
	defer fnc(f)
	w := new(bytes.Buffer)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	err = enc.Encode(data)
	if err != nil {
		return err
	}
	_, err = f.Write(w.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func WriteToCsv(filename string, data interface{}) error {
	var d []byte
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	fnc := filemgmt.CloseFile()
	defer fnc(f)
	_, err = f.Write(utf8Bom)
	if err != nil {
		return err
	}
	d, err = csvutil.Marshal(data)
	if err != nil {
		return err
	}
	_, err = f.Write(d)
	if err != nil {
		return err
	}
	return nil
}
