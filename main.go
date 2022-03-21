package main

import (
	"bufio"
	"fmt"
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
	go func() {
		for _, path := range paths {
			readFile(path, out)
		}
	}()

	go func() {
		maxLines, _ := strconv.Atoi(os.Args[4])
		name := os.Args[2] + "/" + os.Args[3] + strconv.Itoa(fileCount)
		f, _ := os.Create(name)
		for {
			if count%maxLines == 0 {
				name = os.Args[2] + "/" + os.Args[3] + strconv.Itoa(fileCount)
				f, err = os.Create(name)

				if err != nil {
					log.Println(err)
				}
				fileCount++

				defer f.Close()

			}
			res := <-out
			_, err2 := f.WriteString(res + "\n")
			fmt.Println("writing", res)

			if err2 != nil {
				log.Println(err2)
			}
		}
	}()
	fmt.Println(os.Args[1], os.Args[2])
	wg.Wait()
}
func readFile(path string, out chan string) {
	file, err := os.Open(path)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		fmt.Println("reading", scanner.Text())
		out <- scanner.Text()
		count++
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
