# basicweb

## Description

**basiweb** is a very light web server written in [Go](https://golang.org/).  
Here are the specifications:

- no external dependencies
- less than 100 lines of code
- very light cache management
- light [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) management
- light virtual host management
- upload file with **POST** or **PUT** HTTP verb
- remove file with **DELETE** HTTP verb
- protect modifications with basic authentication
- easy configuration with command-line parameters

## How to get it

```bash
git clone https://github.com/cyd01/basicweb.git
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
Usage of C:\Users\cyril\scoop\home\src\basicweb\basicweb.exe:
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
