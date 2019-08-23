package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
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

	pollMs, err := time.ParseDuration(poll + "ms")
	if err != nil {
		panic(err)
	}

	seal := Seal{
		Directory:  dir,
		Extensions: extensions,
		Command:    cmd,
		PollMs:     pollMs,
	}

	seal.PollDir(currentTime)
}

// Seal holds cli args, process info, and a mutex
type Seal struct {
	sync.Mutex
	Directory  string
	Found      []string
	Extensions []string
	Command    string
	PollMs     time.Duration
	Process    *os.Process
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

// PollDir polls given directory and runs given command if files are changed
func (s *Seal) PollDir(currentTime time.Time) {
	err := s.Scan()
	if err != nil {
		panic(err)
	}

	for _, file := range s.Found {
		info, err := os.Stat(file)
		if err != nil {
			fmt.Println(err)
		}

		if info.ModTime().Unix() > currentTime.Unix() {
			args := strings.Split(s.Command, " ")

			root := args[0]
			rest := args[1:]

			if s.Process != nil {
				fmt.Println("seal: reloading..")

				_, perr := os.FindProcess(s.Process.Pid)
				if perr == nil {
					err := s.Process.Kill()
					if err != nil {
						panic(err)
					}
				}
			}

			go func() {
				s.Lock()

				cmd := exec.Command(root, rest...)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				cmd.Start()

				s.Process = cmd.Process

				s.Unlock()

				cmdErr := cmd.Wait()
				if cmdErr != nil {
					fmt.Println(cmdErr)
				}
			}()

			currentTime = time.Now()
		}
	}

	time.Sleep(s.PollMs)

	s.PollDir(currentTime)
}
