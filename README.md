# basicweb

## Description

**basicweb** is a very light web server written in [Go](https://golang.org/).  
Here are the specifications:

- very few external dependencies
- just a little more than 600 lines of code  ;-)
- very light cache management
- light [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) management
- light virtual host management
- HTTP/2 and HTTP/3 compatibility
- TLS with ACME management
- upload files with **POST** or **PUT** HTTP verb
- remove files with **DELETE** HTTP verb
- protect against modifications with basic authentication
- force status code responses
- light dynamic scripts management
- easy configuration with command-line parameters
- simple echo server (send the request content into a JSON structure)

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
  for ARCH in amd64 arm64 ; do
    echo "Building ${BIN}_${OS}_${ARCH}..."
    GOOS=$OS GOARCH=$ARCH go build -o ${BIN}_${OS}_${ARCH}
  done
done
```

```log
Building basicweb_windows_386...
Building basicweb_windows_amd64...
Building basicweb_linux_386...
Building basicweb_linux_amd64...
Building basicweb_linux_arm64...
Building basicweb_linux_386...
Building basicweb_linux_amd64...
Building basicweb_linux_arm64...
```

## Usage

```bash
$ ./basicweb -h
Usage of ./basicweb:
  -acme string
    	directory URL of ACME server
  -cmd string
    	external command (/path1/=cmd1,...)
  -delay int
    	delay (in seconds) before response
  -dir string
    	root directory (default ".")
  -echo
    	start echo web server
  -follow
    	add a follow redirect (302) from /follow to /
  -headers string
    	add specific headers (header1=value1[,...])
  -http3
    	active HTTP/3 mode (over TCP)
  -mime string
    	add new type mime (coma separated extention:value list)
  -multi
    	start HTTP and HTTPS on same port
  -nocache
    	force not to cache
  -pass string
    	password for basic authentication (modification only)
  -port string
    	port web server (default "80")
  -ssl
    	active SSL with key.pem and cert.pem files
  -sslcert string
    	SSL certificate (default "cert.pem")
  -sslkey string
    	SSL private key (default "key.pem")
  -status int
    	force return code
  -timeout int
    	timeout for external command (default 30)
  -tls13
    	force TLS 1.3
  -udp
    	change UDP mode for HTTP/3
  -user string
    	username for basic authentication (modification only)
```

## Lightest start command

```bash
./basicweb -port 8080
2024/05/29 11:01:04 ☢ Starting web server
2024/05/29 11:01:04 on :8080 with directory . with status response 0
2024/05/29 11:01:04 Start HTTP/1.1 (+ HTTP/2 with h2c) server
```

## Docker image

```bash
$ docker build . -t basicweb && docker run --rm -p 8080:80 basicweb
2020/12/05 13:51:32 Starting web server with port 80 on directory . with status response 0

```

## Dynamic scripts example

```bash
$ ./basicweb -cmd "/cmd/=/bin/bash -c cmd.sh"
```

## Start echo werver

```bash
$ ./basicweb -echo
```

## Start SSL mode

```bash
$ ./basicweb -ssl
```

> Private key in PEM format must be provided in `key.pem` file, and Certificate in PEM format must be provided in `cert.pem` file. It is also possible to use `-sslkey` and `-sslcert` options.

## Start HTTP/3 mode

```bash
$ ./basicweb -http3
```
