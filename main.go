package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	currentTime := time.Now()

	var ext string
	flag.StringVar(&ext, "ext", "", `REQUIRED
	a comma delimted list of file extensions to scan
	if none are given all files will be searched`)

	var cmd string
	flag.StringVar(&cmd, "cmd", "", `REQUIRED
	command written as it were to be written in the terminal surrounded in double quotes`)

	var dir string
	flag.StringVar(&dir, "dir", ".", `OPTIONAL
	directory where seal will poll`)

	var poll string
	flag.StringVar(&poll, "poll", "1000", `OPTIONAL
	time spent between directory scans`)

	flag.Parse()

	extensions := strings.Split(ext, ",")

	pollMs, err := strconv.Atoi(poll)
	if err != nil {
		panic(err)
	}

	seal := Seal{
		Directory:  dir,
		Extensions: extensions,
		Command:    cmd,
		PollMs:     pollMs,
	}

	seal.pollDir(seal, currentTime)
}

func (s *Seal) pollDir(scnnr Seal, currentTime time.Time) {
	err := scnnr.Scan()
	if err != nil {
		panic(err)
	}

	for _, file := range scnnr.Found {
		info, err := os.Stat(file)
		if err != nil {
			fmt.Println(err)
		}

		if info.ModTime().Unix() > currentTime.Unix() {
			args := strings.Split(s.Command, " ")

			root := args[0]
			rest := args[1:]

			cmd := exec.Command(root, rest...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			cmdErr := cmd.Run()
			if cmdErr != nil {
				log.Fatal(cmdErr)
			}

			currentTime = time.Now()
		}
	}

	time.Sleep(1 * time.Second)

	s.pollDir(scnnr, currentTime)
}

// Seal holds cli args
type Seal struct {
	Directory  string
	Found      []string
	Extensions []string
	Command    string
	PollMs     int
}

// Scan walks the given directory tree
func (s *Seal) Scan() error {
	err := filepath.Walk(s.Directory, s.scan)
	if err != nil {
		return err
	}

	return nil
}

func (s *Seal) scan(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.IsDir() {
		for _, extension := range s.Extensions {
			if filepath.Ext(path) == extension {
				s.Found = append(s.Found, path)
			}
		}
	}

	return nil
}
