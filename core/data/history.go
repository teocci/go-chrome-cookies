// Package data
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-14
package data

import (
	"database/sql"
	"fmt"
	"github.com/teocci/go-chrome-cookies/filemgmt"
	"github.com/teocci/go-chrome-cookies/logger"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type history struct {
	Title         string
	Url           string
	VisitCount    int
	LastVisitTime time.Time
}

type historyData struct {
	mainPath string
	history  []history
}

func NewHistoryData(main, sub string) Item {
	return &historyData{mainPath: main}
}

func (h *historyData) ChromeParse(key []byte) error {
	historyDB, err := sql.Open("sqlite3", ChromeHistoryFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := historyDB.Close(); err != nil {
			logger.Error(err)
		}
	}()
	rows, err := historyDB.Query(QueryChromiumHistory)
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
			url, title    string
			visitCount    int
			lastVisitTime int64
		)
		err := rows.Scan(&url, &title, &visitCount, &lastVisitTime)
		hData := history{
			Url:           url,
			Title:         title,
			VisitCount:    visitCount,
			LastVisitTime: filemgmt.TimeEpochFormat(lastVisitTime),
		}
		if err != nil {
			logger.Error(err)
		}
		h.history = append(h.history, hData)
	}
	return nil
}

func (h *historyData) FirefoxParse() error {
	var (
		err         error
		keyDB       *sql.DB
		historyRows *sql.Rows
		tempMap     map[int64]string
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
	historyRows, err = keyDB.Query(QueryFirefoxHistory)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer func() {
		if err := historyRows.Close(); err != nil {
			logger.Error(err)
		}
	}()
	for historyRows.Next() {
		var (
			id, visitDate int64
			url, title    string
			visitCount    int
		)
		err = historyRows.Scan(&id, &url, &visitDate, &title, &visitCount)
		if err != nil {
			logger.Warn(err)
		}
		h.history = append(h.history, history{
			Title:         title,
			Url:           url,
			VisitCount:    visitCount,
			LastVisitTime: filemgmt.TimeStampFormat(visitDate / 1000000),
		})
		tempMap[id] = url
	}
	return nil
}

func (h *historyData) CopyDB() error {
	return CopyToLocalPath(h.mainPath, filepath.Base(h.mainPath))
}

func (h *historyData) Release() error {
	return os.Remove(filepath.Base(h.mainPath))
}

func (h *historyData) OutPut(format OutputFormat, browser, dir string) error {
	sort.Slice(h.history, func(i, j int) bool {
		return h.history[i].VisitCount > h.history[j].VisitCount
	})
	switch format {
	case formatCSV:
		err := h.outPutCsv(browser, dir)
		return err
	case formatConsole:
		h.outPutConsole()
		return nil
	default:
		err := h.outPutJson(browser, dir)
		return err
	}
}

func (h *historyData) outPutJson(browser, dir string) error {
	filename := filemgmt.FormatFileName(dir, browser, ItemNameHistory, GetFormatName(formatJson))
	sort.Slice(h.history, func(i, j int) bool {
		return h.history[i].VisitCount > h.history[j].VisitCount
	})
	err := WriteToJson(filename, h.history)
	if err != nil {
		return err
	}
	fmt.Printf("%s Get %d history, filename is %s \n", filemgmt.Prefix, len(h.history), filename)
	return nil
}

func (h *historyData) outPutCsv(browser, dir string) error {
	filename := filemgmt.FormatFileName(dir, browser, ItemNameHistory, GetFormatName(formatCSV))
	if err := WriteToCsv(filename, h.history); err != nil {
		return err
	}
	fmt.Printf("%s Get %d history, filename is %s \n", filemgmt.Prefix, len(h.history), filename)
	return nil
}

func (h *historyData) outPutConsole() {
	for _, v := range h.history {
		fmt.Printf("%+v\n", v)
	}
}
