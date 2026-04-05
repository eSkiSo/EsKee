# EsKee

_Simple terminal keepass client in go_


![Screenshot](/screenshot.png)

## Features

- Quickly search from all entries

- Password is hidden by default (press **v** to see it)

- Copy password directly to **clipboard**

- See all details by pressing **Enter**


## Install

Move **kee** to a folder that is on your $PATH (/usr/bin/local in linux or mac for example*).

## Usage

### See version

```bash
kee -v
```

### Open database

```bash 
kee <kdbx password> <kdbx file|optional>
```

## Database

By default it will load **Database.kdbx** file on your home folder

You can change that behavior by adding an enviroment variable called **KEE_DB** or just pass it as 3rd parameter:

```bash 

export KEE_DB=~/MyOtherDatabase.kdbx
kee <kdbx password> 

-- or --

kee <kdbx password> ~/MyOtherDatabase.kdbx


```

## Compile

```bash 

GOOS=windows GOARCH=amd64 go build -o bin/windows/kee.exe kee.go
GOOS=linux GOARCH=amd64 go build -o bin/linux/kee kee.go
GOOS=darwin GOARCH=amd64 go build -o bin/macos/kee kee.go

```
or run

```bash 
bash build.sh
```

## Started

> 2026/04/04