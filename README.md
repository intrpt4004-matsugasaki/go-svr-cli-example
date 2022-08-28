# go-svr-cli-sample

## Windows
### Client
```sh
cd C:\Users\{user}
git clone https://github.com/intrpt4004-matsugasaki/go-svr-cli-sample
cd go-svr-cli-sample/client
go install fyne.io/fyne/v2/cmd/fyne@latest
C:\Users\{user}\go\bin\fyne package -os windows
.\client.exe
```

### Server
```sh
cd ../server
go run .
```
