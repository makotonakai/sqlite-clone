package main

import (
	"fmt"
	"os"
)

const (
	TABLE_MAX_PAGES = 100
)

type Pager struct {
	file       *os.File
	fileLength uint32
	numPages   uint32
	pages      [TABLE_MAX_PAGES][]byte
}

func pagerOpen(filename string) (*Pager, error) {

	file, err := os.OpenFile(
		filename,
		os.O_RDWR|os.O_CREATE,
		0644,
	)

	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	fileLength := uint32(info.Size())

	if fileLength%PAGE_SIZE != 0 {
		return nil,
			fmt.Errorf(
				"db file is not a whole number of pages. corrupt file",
			)
	}

	pager := &Pager{
		file:       file,
		fileLength: fileLength,
		numPages:   fileLength / PAGE_SIZE,
	}

	return pager, nil
}

func getPage(
	pager *Pager,
	pageNum uint32,
) []byte {

	if pageNum >= TABLE_MAX_PAGES {
		panic(
			fmt.Sprintf(
				"Tried to fetch page number out of bounds. %d >= %d",
				pageNum,
				TABLE_MAX_PAGES,
			),
		)
	}

	if pager.pages[pageNum] == nil {

		page := make([]byte, PAGE_SIZE)

		numPages := pager.fileLength / PAGE_SIZE

		if pager.fileLength%PAGE_SIZE != 0 {
			numPages++
		}

		if pageNum < numPages {

			_, err :=
				pager.file.ReadAt(
					page,
					int64(pageNum*PAGE_SIZE),
				)

			if err != nil &&
				err.Error() != "EOF" {
				panic(err)
			}
		}

		pager.pages[pageNum] = page

		if pageNum >= pager.numPages {
			pager.numPages = pageNum + 1
		}
	}

	return pager.pages[pageNum]
}

func pagerFlush(
	pager *Pager,
	pageNum uint32,
) error {

	page := pager.pages[pageNum]

	if page == nil {
		return fmt.Errorf(
			"tried to flush null page",
		)
	}

	_, err :=
		pager.file.WriteAt(
			page,
			int64(pageNum*PAGE_SIZE),
		)

	return err
}
