# tinyimg

[![Go Report Card](https://goreportcard.com/badge/github.com/zcong1993/tinyimg)](https://goreportcard.com/report/github.com/zcong1993/tinyimg)
[![CircleCI](https://circleci.com/gh/zcong1993/tinyimg/tree/master.svg?style=shield)](https://circleci.com/gh/zcong1993/tinyimg/tree/master)

> image compress cli

## Install

### install with go

```sh
$ go get -u -v github.com/zcong1993/tinyimg
# will spend several minutes cause lib github.com/discordapp/lilliput is a bit large
```

### build yourself

```sh
$ git clone https://github.com/zcong1993/tinyimg.git
$ cd tinyimg
# install deps, use go dep
$ dep ensure
# or normal
$ go get -v ./...
# build
$ chmod +x ./build.sh
$ make build
# copy to any folder in your $PATH
$ cp ./bin/tinyimg $GOPATH/bin/
```

## Usage

```sh
$ tinyimg [options] source files...
# example
$ tinyimg -q=80 -o=outdir ./*.jpg
# show help
$ tinyimg -h
# show version
$ tinyimg -v
```

## License

MIT &copy; zcong1993
