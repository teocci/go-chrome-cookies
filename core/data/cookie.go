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
	"time"
)

type cookie struct {
	Host         string
	Path         string
	KeyName      string
	encryptValue []byte
	Value        string
	IsSecure     bool
	IsHTTPOnly   bool
	HasExpire    bool
	IsPersistent bool
	CreateDate   time.Time
	ExpireDate   time.Time
}

type cookies struct {
	mainPath string
	cookies  map[string][]cookie
}

func NewCookies(main, sub string) Item {
	return &cookies{mainPath: main}
}

func (c *cookies) ChromeParse(secretKey []byte) error {
	c.cookies = make(map[string][]cookie)
	cookieDB, err := sql.Open("sqlite3", ChromeCookieFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := cookieDB.Close(); err != nil {
			logger.Debug(err)
		}
	}()
	rows, err := cookieDB.Query(QueryChromiumCookie)
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
			key, host, path                               string
			isSecure, isHTTPOnly, hasExpire, isPersistent int
			createDate, expireDate                        int64
			value, encryptValue                           []byte
		)
		err = rows.Scan(&key, &encryptValue, &host, &path, &createDate, &expireDate, &isSecure, &isHTTPOnly, &hasExpire, &isPersistent)
		if err != nil {
			logger.Error(err)
		}
		cookie := cookie{
			KeyName:      key,
			Host:         host,
			Path:         path,
			encryptValue: encryptValue,
			IsSecure:     filemgmt.IntToBool(isSecure),
			IsHTTPOnly:   filemgmt.IntToBool(isHTTPOnly),
			HasExpire:    filemgmt.IntToBool(hasExpire),
			IsPersistent: filemgmt.IntToBool(isPersistent),
			CreateDate:   filemgmt.TimeEpochFormat(createDate),
			ExpireDate:   filemgmt.TimeEpochFormat(expireDate),
		}
		// remove 'v10'
		if secretKey == nil {
			value, err = decrypt.DPApi(encryptValue)
		} else {
			value, err = decrypt.ChromePass(secretKey, encryptValue)
		}
		if err != nil {
			logger.Debug(err)
		}
		cookie.Value = string(value)
		c.cookies[host] = append(c.cookies[host], cookie)
	}
	return nil
}

func (c *cookies) FirefoxParse() error {
	c.cookies = make(map[string][]cookie)
	cookieDB, err := sql.Open("sqlite3", FirefoxCookieFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := cookieDB.Close(); err != nil {
			logger.Debug(err)
		}
	}()
	rows, err := cookieDB.Query(QueryFirefoxCookie)
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
			name, value, host, path string
			isSecure, isHttpOnly    int
			creationTime, expiry    int64
		)
		err = rows.Scan(&name, &value, &host, &path, &creationTime, &expiry, &isSecure, &isHttpOnly)
		if err != nil {
			logger.Error(err)
		}
		c.cookies[host] = append(c.cookies[host], cookie{
			KeyName:    name,
			Host:       host,
			Path:       path,
			IsSecure:   filemgmt.IntToBool(isSecure),
			IsHTTPOnly: filemgmt.IntToBool(isHttpOnly),
			CreateDate: filemgmt.TimeStampFormat(creationTime / 1000000),
			ExpireDate: filemgmt.TimeStampFormat(expiry),
			Value:      value,
		})
	}
	return nil
}

func (c *cookies) CopyDB() error {
	return CopyToLocalPath(c.mainPath, filepath.Base(c.mainPath))
}

func (c *cookies) Release() error {
	return os.Remove(filepath.Base(c.mainPath))
}

func (c *cookies) OutPut(format OutputFormat, browser, dir string) error {
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

func (c *cookies) outPutJson(browser, dir string) error {
	filename := filemgmt.FormatFileName(dir, browser, ItemNameCookie, GetFormatName(formatJson))
	err := WriteToJson(filename, c.cookies)
	if err != nil {
		return err
	}
	fmt.Printf("%s Get %d cookies, filename is %s \n", filemgmt.Prefix, len(c.cookies), filename)
	return nil
}

func (c *cookies) outPutCsv(browser, dir string) error {
	filename := filemgmt.FormatFileName(dir, browser, ItemNameCookie, GetFormatName(formatCSV))
	var tempSlice []cookie
	for _, v := range c.cookies {
		tempSlice = append(tempSlice, v...)
	}
	if err := WriteToCsv(filename, tempSlice); err != nil {
		return err
	}
	fmt.Printf("%s Get %d cookies, filename is %s \n", filemgmt.Prefix, len(c.cookies), filename)
	return nil
}

func (c *cookies) outPutConsole() {
	for host, value := range c.cookies {
		fmt.Printf("%s\n%+v\n", host, value)
	}
}
