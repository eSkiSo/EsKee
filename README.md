## EsKee

### Simple terminal keepass client in go


### Install

> Move **kee** to a folder that is on your $PATH (/usr/bin/local in linux or mac for example*) or add its location to the $PATH variable.

### Usage

> kee -v # See version
> kee **<kdbx password>** *<kdbx file|optional>* # Open database

### Database

> By default it will load **Database.kdbx** file on your home folder
> You can change that behavior by adding an enviroment variable called **KEE_DB="your_db_file_location"** or just pass it as 3rd parameter

### Compile

```bash 

GOOS=windows GOARCH=amd64 go build -o bin/windows/kee.exe kee.go
GOOS=linux GOARCH=amd64 go build -o bin/linux/kee kee.go
GOOS=darwin GOARCH=amd64 go build -o bin/macos/kee kee.go

```
or run bash build.sh*

## Started

> 2026/04/04