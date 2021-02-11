# basicweb

## Description

**basiweb** is a very light web server written in [Go](https://golang.org/).  
Here are the specifications:

- no external dependencies
- just a little mode than 100 lines of code
- very light cache management
- light [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) management
- light virtual host management
- upload files with **POST** or **PUT** HTTP verb
- remove files with **DELETE** HTTP verb
- protect against modifications with basic authentication
- force status code responses
- light dynamic scripts managment
- easy configuration with command-line parameters

## How to get it

```bash
git clone https://gitlab.techlabfdj.io/cyd/basicweb.git
```

## How to build it

Compilation example for `windows`, `linux` and `macos`.

```bash
BIN=basicweb
for OS in windows ; do
  for ARCH in 386 amd64 ; do
    echo "Building ${BIN}_${OS}_${ARCH}..."
    GOOS=$OS GOARCH=$ARCH go build -o ${BIN}_${OS}_${ARCH}.exe
  done
done
for OS in linux ; do
  for ARCH in 386 amd64 arm64 ; do
    echo "Building ${BIN}_${OS}_${ARCH}..."
    GOOS=$OS GOARCH=$ARCH go build -o ${BIN}_${OS}_${ARCH}
  done
done
for OS in darwin ; do
  for ARCH in 386 amd64 ; do
    echo "Building ${BIN}_${OS}_${ARCH}..."
    GOOS=$OS GOARCH=$ARCH go build -o ${BIN}_${OS}_${ARCH}
  done
done
```

## Usage

```bash
$ ./basicweb -h
Usage of ./basicweb:
  -cmd string
        external command
  -dir string
        root directory (default ".")
  -nocache
        force not to cache
  -pass string
        password for basic authentication (modification only)
  -port string
        port web server (default "80")
  -status int
        force return code
  -user string
        username for basic authentication (modification only)
```

## Lightest start command

```bash
$ ./basicweb
2020/12/03 18:04:59 Starting web server with port 80 on directory . with status response 0
2020/12/03 18:05:12 GET /
```

## Docker image

```bash
$ docker run --rm -p 8080:80 basicweb
2020/12/05 13:51:32 Starting web server with port 80 on directory . with status response 0

```

## Dynamic scripts example

```bash
$ ./basicweb -cmd "/bin/bash -c cmd.sh"
```
