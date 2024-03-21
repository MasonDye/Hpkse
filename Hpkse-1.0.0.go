package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

func searchInFile(filepath string, query string, wg *sync.WaitGroup, ch chan<- bool) {
	defer wg.Done()
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filepath, err)
		return
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 1024*1024) // 1MB 缓冲区大小

	found := false
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading file %s: %v\n", filepath, err)
			}
			break
		}
		if regexp.MustCompile(regexp.QuoteMeta(query)).MatchString(line) {
			fmt.Printf("%s: %s\n", filepath, line)
			found = true
		}
	}

	if found {
		ch <- true
	}
}

func searchInDirectory(directory string, query string) {
	var wg sync.WaitGroup
	ch := make(chan bool)

	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			wg.Add(1)
			go searchInFile(path, query, &wg, ch)
		}
		return nil
	})

	go func() {
		wg.Wait()
		close(ch)
	}()

	found := false
	for range ch {
		found = true
	}

	if !found {
		fmt.Println("Not found keyword")
	}
}

func main() {
	fmt.Println("Hpkse High performance keyword search engine")
	fmt.Println("Version : 1.0.0")
	fmt.Println("Author by MasonDye\n")

	var queryDirectory string

	if len(os.Args) > 1 && os.Args[1] == "-data" {
		if len(os.Args) > 2 {
			queryDirectory = os.Args[2]
		} else {
			fmt.Println("Folder not provided")
			return
		}
	} else {
		fmt.Print("Search folder (relative path): ")
		fmt.Scanln(&queryDirectory)
	}

	var query string
	fmt.Print("Search keyword: ")
	fmt.Scanln(&query)

	if fileInfo, err := os.Stat(queryDirectory); err == nil && fileInfo.IsDir() {
		searchInDirectory(queryDirectory, query)
	} else {
		fmt.Println("input path is not a folder.")
	}
}
