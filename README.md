# Seal

Windows, Linux, and macOS/Unix file extension watcher and command executor.

Polls files by default every 1000ms (1sec) and if any file has been updated (saved/created) it will run the given command.

Perfect for spiking scripts!

### Install

Currently just supports using go as a package manager/installer:

```
go get github.com/selfup/seal
go install github.com/selfup/seal
```

### Expected Args

```bash
seal -ext=".go,.proto" -cmd="go run main.go"
```

### Chained Commands

If you want to use a command like: `echo wow | grep ow && ls -lagh` you will need to write the script in a file and tell seal to execute said file.

Example: `seal -ext=".sh" -cmd="./scripts/run.sh"`

### Overriding defaults

Watched directory: `-dir=`

Polling time in ms: `-poll=`

```bash
seal -dir="pkg" -ext=".go,.proto" -cmd="go test ./pkg/..." -poll="300"
```

### Help

```
$ seal -h
Usage of C:\Users\selfup\go\bin\seal.exe:
  -cmd string
        REQUIRED
                command written as it were to be written in the terminal surrounded in double quotes
  -dir string
        OPTIONAL
                directory where seal will poll (default ".")
  -ext string
        REQUIRED
                a comma delimted list of file extensions to scan
                if none are given all files will be searched
  -poll string
        OPTIONAL
                time spent between directory scans (default "1000")
```
