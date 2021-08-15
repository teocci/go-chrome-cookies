// Package data
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-14
package data

import (
	"database/sql"
	"fmt"
	"github.com/teocci/go-chrome-cookies/filemgmt"
	"github.com/teocci/go-chrome-cookies/logger"
	"github.com/tidwall/gjson"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type download struct {
	TargetPath string
	Url        string
	TotalBytes int64
	StartTime  time.Time
	EndTime    time.Time
	MimeType   string
}

type downloads struct {
	mainPath  string
	downloads []download
}

func NewDownloads(main, sub string) Item {
	return &downloads{mainPath: main}
}

func (d *downloads) ChromeParse(key []byte) error {
	historyDB, err := sql.Open("sqlite3", ChromeDownloadFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := historyDB.Close(); err != nil {
			logger.Error(err)
		}
	}()
	rows, err := historyDB.Query(QueryChromiumDownload)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error(err)
		}
	}()
	for rows.Next() {
		var (
			targetPath, tabUrl, mimeType   string
			totalBytes, startTime, endTime int64
		)
		err := rows.Scan(&targetPath, &tabUrl, &totalBytes, &startTime, &endTime, &mimeType)
		data := download{
			TargetPath: targetPath,
			Url:        tabUrl,
			TotalBytes: totalBytes,
			StartTime:  filemgmt.TimeEpochFormat(startTime),
			EndTime:    filemgmt.TimeEpochFormat(endTime),
			MimeType:   mimeType,
		}
		if err != nil {
			logger.Error(err)
		}
		d.downloads = append(d.downloads, data)
	}
	return nil
}

func (d *downloads) FirefoxParse() error {
	var (
		err          error
		keyDB        *sql.DB
		downloadRows *sql.Rows
		tempMap      map[int64]string
	)
	tempMap = make(map[int64]string)
	keyDB, err = sql.Open("sqlite3", FirefoxDataFile)
	if err != nil {
		return err
	}
	_, err = keyDB.Exec(CloseJournalMode)
	if err != nil {
		logger.Error(err)
	}
	defer func() {
		if err := keyDB.Close(); err != nil {
			logger.Error(err)
		}
	}()
	downloadRows, err = keyDB.Query(QueryFirefoxDownload)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer func() {
		if err := downloadRows.Close(); err != nil {
			logger.Error(err)
		}
	}()
	for downloadRows.Next() {
		var (
			content, url       string
			placeID, dateAdded int64
		)
		err = downloadRows.Scan(&placeID, &content, &url, &dateAdded)
		if err != nil {
			logger.Warn(err)
		}
		contentList := strings.Split(content, ",{")
		if len(contentList) > 1 {
			path := contentList[0]
			json := "{" + contentList[1]
			endTime := gjson.Get(json, "endTime")
			fileSize := gjson.Get(json, "fileSize")
			d.downloads = append(d.downloads, download{
				TargetPath: path,
				Url:        url,
				TotalBytes: fileSize.Int(),
				StartTime:  filemgmt.TimeStampFormat(dateAdded / 1000000),
				EndTime:    filemgmt.TimeStampFormat(endTime.Int() / 1000),
			})
		}
		tempMap[placeID] = url
	}
	return nil
}

func (d *downloads) CopyDB() error {
	return CopyToLocalPath(d.mainPath, filepath.Base(d.mainPath))
}

func (d *downloads) Release() error {
	return os.Remove(filepath.Base(d.mainPath))
}

func (d *downloads) OutPut(format OutputFormat, browser, dir string) error {
	switch format {
	case formatCSV:
		err := d.outPutCsv(browser, dir)
		return err
	case formatConsole:
		d.outPutConsole()
		return nil
	default:
		err := d.outPutJson(browser, dir)
		return err
	}
}

func (d *downloads) outPutJson(browser, dir string) error {
	filename := filemgmt.FormatFileName(dir, browser, ItemNameDownload, GetFormatName(formatJson))
	err := WriteToJson(filename, d.downloads)
	if err != nil {
		return err
	}
	fmt.Printf("%s Get %d history, filename is %s \n", filemgmt.Prefix, len(d.downloads), filename)
	return nil
}

func (d *downloads) outPutCsv(browser, dir string) error {
	filename := filemgmt.FormatFileName(dir, browser, ItemNameDownload, GetFormatName(formatCSV))
	if err := WriteToCsv(filename, d.downloads); err != nil {
		return err
	}
	fmt.Printf("%s Get %d downloads history, filename is %s \n", filemgmt.Prefix, len(d.downloads), filename)
	return nil
}

func (d *downloads) outPutConsole() {
	for _, v := range d.downloads {
		fmt.Printf("%+v\n", v)
	}
}

func (d downloads) Len() int {
	return len(d.downloads)
}

func (d downloads) Less(i, j int) bool {
	return d.downloads[i].StartTime.After(d.downloads[j].StartTime)
}

func (d downloads) Swap(i, j int) {
	d.downloads[i], d.downloads[j] = d.downloads[j], d.downloads[i]
}
