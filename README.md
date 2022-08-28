# go-svr-cli-sample

## Windows
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

### make .exe
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
