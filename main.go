package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var count = 0
var fileCount = 0
var wg sync.WaitGroup

func main() {
	paths, err := getPaths(os.Args[1])
	if err != nil {
		log.Println(err)
	}

	out := make(chan string)

	wg.Add(len(paths))

	for _, path := range paths {
		go readFile(path, out)
	}

	go writeToFile(out)

	wg.Wait()
}

func writeToFile(out chan string) {
	maxLines, _ := strconv.Atoi(os.Args[4])
	name := os.Args[2] + "/" + os.Args[3] + strconv.Itoa(fileCount)

	f, err := os.Create(name)
	if err != nil {
		os.Mkdir(os.Args[2], os.ModePerm)
	}

	for {
		if count%maxLines == 0 {
			name = os.Args[2] + "/" + os.Args[3] + strconv.Itoa(fileCount)
			f, err = os.Create(name)

			if err != nil {
				os.Mkdir(os.Args[2], os.ModePerm)
			}

			fileCount++

			defer f.Close()

		}

		count++

		_, err2 := f.WriteString(<-out + "\n")
		if err2 != nil {
			log.Println(err2)
		}
	}
	defer wg.Done()
}

func readFile(path string, out chan string) {
	file, err := os.Open(path)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		out <- scanner.Text()
	}

	file.Close()
	defer wg.Done()
}

func getPaths(root string) ([]string, error) {
	var paths []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() == false {
			paths = append(paths, path)
		}
		return nil
	})

	return paths, err
}
