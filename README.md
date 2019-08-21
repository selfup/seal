# Seal

Windows, Linux, and macOS/Unix file extension watcher and command executor.

Polls files by default every 1000ms (1sec) and if any file has been updated (saved/created) it will run the given command.

Perfect for spiking scripts!

### Expected Args

```bash
seal -ext=".go,.proto" -cmd="go run main.go"
```

### Overriding defaults

Watched directory: `-dir=`
Polling time in ms: `-poll=`

```bash
seal -dir="pkg" -ext=".go,.proto" -cmd="go test ./pkg/..." -poll="300"
```
