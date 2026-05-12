package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sync"
)

const FilePath string = "../files/"

func worker(id int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for fileName := range jobs {

		fileDirectory := fmt.Sprintf("%s%s", FilePath, fileName)
		text, err := os.ReadFile(fileDirectory)
		if err != nil {
			fmt.Printf("file %s is not valid error: %v", fileDirectory, err)
		}

		countWordsInFiles := len(bytes.Fields(text))

		results <- fmt.Sprintf("File name: %s Count worlds in file: %d", fileName, countWordsInFiles)
	}
}

func main() {
	files, err := os.ReadDir(FilePath)
	if err != nil {
		fmt.Printf("directory is not valid: %v", err)
	}
	workers := runtime.NumCPU()

	wg := sync.WaitGroup{}

	results := make(chan string)
	jobs := make(chan string)

	go func() {
		defer close(jobs)
		for _, f := range files {

			jobs <- f.Name()
		}
	}()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for n := range results {
		fmt.Println(n)
	}

}
