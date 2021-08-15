// Package data
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-14
package data

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/teocci/go-chrome-cookies/filemgmt"
	"github.com/teocci/go-chrome-cookies/logger"
	"github.com/tidwall/gjson"
)

const (
	bookmarkID       = "id"
	bookmarkAdded    = "date_added"
	bookmarkUrl      = "url"
	bookmarkName     = "name"
	bookmarkType     = "type"
	bookmarkChildren = "children"
)

type bookmark struct {
	ID        int64
	Name      string
	Type      string
	URL       string
	DateAdded time.Time
}

type bookmarks struct {
	mainPath  string
	bookmarks []bookmark
}

func NewBookmarks(main, sub string) Item {
	return &bookmarks{mainPath: main}
}

func (b *bookmarks) ChromeParse(key []byte) error {
	bookmarks, err := filemgmt.ReadFile(ChromeBookmarkFile)
	if err != nil {
		return err
	}
	r := gjson.Parse(bookmarks)
	if r.Exists() {
		roots := r.Get("roots")
		roots.ForEach(func(key, value gjson.Result) bool {
			getBookmarkChildren(value, b)
			return true
		})
	}
	return nil
}

func (b *bookmarks) FirefoxParse() error {
	var (
		err          error
		keyDB        *sql.DB
		bookmarkRows *sql.Rows
		tempMap      map[int64]string
		bookmarkUrl  string
	)
	keyDB, err = sql.Open("sqlite3", FirefoxDataFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := keyDB.Close(); err != nil {
			logger.Error(err)
		}
	}()
	_, err = keyDB.Exec(CloseJournalMode)
	if err != nil {
		logger.Error(err)
	}
	bookmarkRows, err = keyDB.Query(QueryFirefoxBookMarks)
	if err != nil {
		return err
	}
	for bookmarkRows.Next() {
		var (
			id, bType, dateAdded int64
			title, url           string
		)
		err = bookmarkRows.Scan(&id, &url, &bType, &dateAdded, &title)
		if err != nil {
			logger.Warn(err)
		}
		if url, ok := tempMap[id]; ok {
			bookmarkUrl = url
		}
		b.bookmarks = append(b.bookmarks, bookmark{
			ID:        id,
			Name:      title,
			Type:      BookMarkType(bType),
			URL:       bookmarkUrl,
			DateAdded: filemgmt.TimeStampFormat(dateAdded / 1000000),
		})
	}
	return nil
}

func (b *bookmarks) CopyDB() error {
	return CopyToLocalPath(b.mainPath, filepath.Base(b.mainPath))
}

func (b *bookmarks) Release() error {
	return os.Remove(filepath.Base(b.mainPath))
}

func (b *bookmarks) OutPut(format OutputFormat, browser, dir string) error {
	sort.Slice(b.bookmarks, func(i, j int) bool {
		return b.bookmarks[i].ID < b.bookmarks[j].ID
	})
	switch format {
	case formatCSV:
		err := b.outPutCsv(browser, dir)
		return err
	case formatConsole:
		b.outPutConsole()
		return nil
	default:
		err := b.outPutJson(browser, dir)
		return err
	}
}

func (b *bookmarks) outPutJson(browser, dir string) error {
	filename := filemgmt.FormatFileName(dir, browser, ItemNameBookmark, GetFormatName(formatJson))
	sort.Slice(b.bookmarks, func(i, j int) bool {
		return b.bookmarks[i].ID < b.bookmarks[j].ID
	})
	err := WriteToJson(filename, b.bookmarks)
	if err != nil {
		return err
	}
	fmt.Printf("%s Get %d bookmarks, filename is %s \n", filemgmt.Prefix, len(b.bookmarks), filename)
	return nil
}

func (b *bookmarks) outPutCsv(browser, dir string) error {
	filename := filemgmt.FormatFileName(dir, browser, ItemNameBookmark, GetFormatName(formatCSV))
	if err := WriteToCsv(filename, b.bookmarks); err != nil {
		return err
	}
	fmt.Printf("%s Get %d bookmarks, filename is %s \n", filemgmt.Prefix, len(b.bookmarks), filename)
	return nil
}

func (b *bookmarks) outPutConsole() {
	for _, v := range b.bookmarks {
		fmt.Printf("%+v\n", v)
	}
}

func getBookmarkChildren(value gjson.Result, b *bookmarks) (children gjson.Result) {
	nodeType := value.Get(bookmarkType)
	bm := bookmark{
		ID:        value.Get(bookmarkID).Int(),
		Name:      value.Get(bookmarkName).String(),
		URL:       value.Get(bookmarkUrl).String(),
		DateAdded: filemgmt.TimeEpochFormat(value.Get(bookmarkAdded).Int()),
	}
	children = value.Get(bookmarkChildren)
	if nodeType.Exists() {
		bm.Type = nodeType.String()
		b.bookmarks = append(b.bookmarks, bm)
		if children.Exists() && children.IsArray() {
			for _, v := range children.Array() {
				children = getBookmarkChildren(v, b)
			}
		}
	}
	return children
}

func BookMarkType(a int64) string {
	switch a {
	case 1:
		return "url"
	default:
		return "folder"
	}
}
