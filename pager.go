package main

import (
	"io"
	"os"
	"fmt"
	"syscall"
)

type Pager struct {
    File *os.File
    FileLength int64
	NumPages uint32
    Pages [TABLE_MAX_PAGES][]byte
}

func PagerOpen(fileName string) *Pager {
    
    file, err := os.OpenFile(
        fileName, 
        os.O_RDWR|os.O_CREATE|syscall.S_IWUSR|syscall.S_IRUSR, 
        0644,
    )

    if err != nil {
        fmt.Println(err)
    }

    stat, err := file.Stat()
    if err != nil {
        fmt.Println(err)
    }

    p := Pager{
        File: file,
        FileLength: stat.Size(),
        NumPages: uint32(stat.Size() / PAGE_SIZE),
    }

    if stat.Size() % PAGE_SIZE != 0 {
        fmt.Printf("Db file is not a whole number of pages. Corrupt file.\n");
        os.Exit(1);
    }

    return &p
}

func PagerFlush(pager *Pager, pageNum uint32) {

    if pager.Pages[pageNum] == nil {
        fmt.Printf("tried to flush nil page\n")
    }

    _, err := pager.File.Seek(
        int64(pageNum)*PAGE_SIZE,
        0,
    )

    if err != nil {
        fmt.Println(err)
    }

    _, err = pager.File.Write(
        pager.Pages[pageNum][:PAGE_SIZE],
    )

    if err != nil {
        fmt.Println(err)
    }
}

func GetPage(pager *Pager, pageNum uint32) []byte {

    if pageNum >= TABLE_MAX_PAGES {
        fmt.Printf("Tried to fetch page number out of bounds. %d > %d\n", 
            pageNum,
            TABLE_MAX_PAGES,
        );
        os.Exit(1)
    }

    if pager.Pages[pageNum] == nil {
        page := make([]byte, PAGE_SIZE)
        numPages := pager.FileLength / PAGE_SIZE

        if pager.FileLength % PAGE_SIZE != 0 {
            numPages++
        }

        if int64(pageNum) < numPages {
            pager.File.Seek(int64(pageNum)*PAGE_SIZE, 0)

            _, err := io.ReadFull(pager.File, page)
            if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
                fmt.Println(err)
                os.Exit(1)
            }
        }

        pager.Pages[pageNum] = page

        if pageNum >= pager.NumPages {
            pager.NumPages = pageNum + 1
        }
    }

    return pager.Pages[pageNum]
}
