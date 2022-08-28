# go-svr-cli-sample
dependencies: Go, MySQL

If you want to run it on Windows, you will need to install some commands via MSYS2.

See how to install fyne.io for details.

## Run
### Client
```sh
git clone https://github.com/intrpt4004-matsugasaki/go-svr-cli-sample
cd go-svr-cli-sample/client
go run .
```

### Server
```sh
cd go-svr-cli-sample/server
go run .
```

## make executables (Windows)
```sh
cd go-svr-cli-sample/client
go install fyne.io/fyne/v2/cmd/fyne@latest
fyne package -os windows
client.exe
```

```sh
cd go-svr-cli-sample/server
go build .
server.exe
```
