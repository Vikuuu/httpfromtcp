package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const inputFile = "message.txt"

func main() {
	f, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open %s: %s\n", inputFile, err)
	}

	fmt.Printf("Reading Data from: %s\n", inputFile)
	fmt.Println("===================================")

	lineChan := getLinesChannel(f)
	for line := range lineChan {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	c := make(chan string)
	go func() {
		defer f.Close()
		defer close(c)
		line := ""
		for {
			b := make([]byte, 8, 8)
			n, err := f.Read(b)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
				break
			}
			str := string(b[:n])
			parts := strings.Split(str, "\n")
			line += parts[0]
			if len(parts) > 1 {
				c <- line
				line = parts[1]
			}
		}
	}()
	return c
}
