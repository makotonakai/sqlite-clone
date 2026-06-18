package main

import (
	"os"
)

const (
	PageSize      = 4096
	TableMaxPages = 100
)

type Pager struct {
	File     *os.File
	FileSize int64

	NumPages uint32

	Pages [TableMaxPages][]byte
}

func OpenPager(filename string) (*Pager, error) {

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
		return nil, err
	}

	numPages := uint32(info.Size() / PageSize)

	if info.Size()%PageSize != 0 {
		numPages++
	}

	return &Pager{
		File:     file,
		FileSize: info.Size(),
		NumPages: numPages,
	}, nil
}

func (p *Pager) GetPage(pageNum uint32) []byte {

	if p.Pages[pageNum] == nil {

		page := make([]byte, PageSize)

		offset := int64(pageNum) * PageSize

		if offset < p.FileSize {
			_, _ = p.File.ReadAt(page, offset)
		}

		p.Pages[pageNum] = page

		if pageNum >= p.NumPages {
			p.NumPages = pageNum + 1
		}
	}

	return p.Pages[pageNum]
}

func (p *Pager) Flush(pageNum uint32) error {

	page := p.Pages[pageNum]

	if page == nil {
		return nil
	}

	offset := int64(pageNum) * PageSize

	_, err := p.File.WriteAt(
		page,
		offset,
	)

	return err
}

func (p *Pager) Close() error {

	for i := uint32(0); i < p.NumPages; i++ {

		if p.Pages[i] == nil {
			continue
		}

		if err := p.Flush(i); err != nil {
			return err
		}
	}

	return p.File.Close()
}
