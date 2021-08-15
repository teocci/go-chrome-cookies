// Package data
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-14
package data

import (
	"database/sql"
	"fmt"
	"github.com/teocci/go-chrome-cookies/core/decrypt"
	"github.com/teocci/go-chrome-cookies/filemgmt"
	"github.com/teocci/go-chrome-cookies/logger"
	"os"
	"path/filepath"
)

type card struct {
	GUID            string
	Name            string
	ExpirationYear  string
	ExpirationMonth string
	CardNumber      string
}

type creditCards struct {
	mainPath string
	cards    map[string][]card
}

func NewCCards(main string, sub string) Item {
	return &creditCards{mainPath: main}
}

func (c *creditCards) FirefoxParse() error {
	return nil // FireFox does not have a credit card saving feature
}

func (c *creditCards) ChromeParse(secretKey []byte) error {
	c.cards = make(map[string][]card)
	creditDB, err := sql.Open("sqlite3", ChromeCreditFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := creditDB.Close(); err != nil {
			logger.Debug(err)
		}
	}()
	rows, err := creditDB.Query(QueryChromiumCredit)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Debug(err)
		}
	}()
	for rows.Next() {
		var (
			name, month, year, guid string
			value, encryptValue     []byte
		)
		err := rows.Scan(&guid, &name, &month, &year, &encryptValue)
		if err != nil {
			logger.Error(err)
		}
		creditCardInfo := card{
			GUID:            guid,
			Name:            name,
			ExpirationMonth: month,
			ExpirationYear:  year,
		}
		if secretKey == nil {
			value, err = decrypt.DPApi(encryptValue)
		} else {
			value, err = decrypt.ChromePass(secretKey, encryptValue)
		}
		if err != nil {
			logger.Debug(err)
		}
		creditCardInfo.CardNumber = string(value)
		c.cards[guid] = append(c.cards[guid], creditCardInfo)
	}
	return nil
}

func (c *creditCards) CopyDB() error {
	return CopyToLocalPath(c.mainPath, filepath.Base(c.mainPath))
}

func (c *creditCards) Release() error {
	return os.Remove(filepath.Base(c.mainPath))
}

func (c *creditCards) OutPut(format OutputFormat, browser, dir string) error {
	switch format {
	case formatCSV:
		err := c.outPutCsv(browser, dir)
		return err
	case formatConsole:
		c.outPutConsole()
		return nil
	default:
		err := c.outPutJson(browser, dir)
		return err
	}
}

func (c *creditCards) outPutJson(browser, dir string) error {
	filename := filemgmt.FormatFileName(dir, browser, ItemNameCreditCard, GetFormatName(formatJson))
	err := WriteToJson(filename, c.cards)
	if err != nil {
		return err
	}
	fmt.Printf("%s Get %d credit cards, filename is %s \n", filemgmt.Prefix, len(c.cards), filename)
	return nil
}

func (c *creditCards) outPutCsv(browser, dir string) error {
	filename := filemgmt.FormatFileName(dir, browser, ItemNameCreditCard, GetFormatName(formatCSV))
	var tempSlice []card
	for _, v := range c.cards {
		tempSlice = append(tempSlice, v...)
	}
	if err := WriteToCsv(filename, tempSlice); err != nil {
		return err
	}
	fmt.Printf("%s Get %d credit cards, filename is %s \n", filemgmt.Prefix, len(c.cards), filename)
	return nil
}

func (c *creditCards) outPutConsole() {
	for _, v := range c.cards {
		fmt.Printf("%+v\n", v)
	}
}
