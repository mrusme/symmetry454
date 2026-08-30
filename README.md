Symmetry454
-----------

[![SEGV
LICENSE](https://img.shields.io/static/v1?label=SEGV%20LICENSE&message=1.1&labelColor=0060A8&color=ffffff)](https://xn--gckvb8fzb.com/segv/)

[<img src="https://xn--gckvb8fzb.com/images/chatroom.png" width="275">](https://xn--gckvb8fzb.com/contact/)

[![GoDoc](https://godoc.org/github.com/mrusme/symmetry454?status.svg)](https://godoc.org/github.com/mrusme/symmetry454)

A Go implementation of the
[Symmetry454](http://individual.utoronto.ca/kalendis/symmetry.htm) calendar.

## Build

Building from a checkout writes a `sym454` binary into the current directory:

```sh
go build ./cmd/sym454
```

To install it into `$(go env GOPATH)/bin` instead, which needs Go 1.16 or
newer:

```sh
go install github.com/mrusme/symmetry454/cmd/sym454@latest
```

## Usage

Run `sym454 -h` for the calendar and converter flags.

```sh
 sym454 2026
                 2026, leap year of 371 days in 53 weeks

┌──────────────────────┐ ┌──────────────────────┐ ┌──────────────────────┐
│       January        │ │       February       │ │        March         │
├──────────────────────┤ ├──────────────────────┤ ├──────────────────────┤
│ Mo Tu We Th Fr Sa Su │ │ Mo Tu We Th Fr Sa Su │ │ Mo Tu We Th Fr Sa Su │
│  1  2  3  4  5  6  7 │ │  1  2  3  4  5  6  7 │ │  1  2  3  4  5  6  7 │
│  8  9 10 11 12 13 14 │ │  8  9 10 11 12 13 14 │ │  8  9 10 11 12 13 14 │
│ 15 16 17 18 19 20 21 │ │ 15 16 17 18 19 20 21 │ │ 15 16 17 18 19 20 21 │
│ 22 23 24 25 26 27 28 │ │ 22 23 24 25 26 27 28 │ │ 22 23 24 25 26 27 28 │
│                      │ │ 29 30 31 32 33 34 35 │ │                      │
└──────────────────────┘ └──────────────────────┘ └──────────────────────┘

┌──────────────────────┐ ┌──────────────────────┐ ┌──────────────────────┐
│        April         │ │         May          │ │         June         │
├──────────────────────┤ ├──────────────────────┤ ├──────────────────────┤
│ Mo Tu We Th Fr Sa Su │ │ Mo Tu We Th Fr Sa Su │ │ Mo Tu We Th Fr Sa Su │
│  1  2  3  4  5  6  7 │ │  1  2  3  4  5  6  7 │ │  1  2  3  4  5  6  7 │
│  8  9 10 11 12 13 14 │ │  8  9 10 11 12 13 14 │ │  8  9 10 11 12 13 14 │
│ 15 16 17 18 19 20 21 │ │ 15 16 17 18 19 20 21 │ │ 15 16 17 18 19 20 21 │
│ 22 23 24 25 26 27 28 │ │ 22 23 24 25 26 27 28 │ │ 22 23 24 25 26 27 28 │
│                      │ │ 29 30 31 32 33 34 35 │ │                      │
└──────────────────────┘ └──────────────────────┘ └──────────────────────┘

┌──────────────────────┐ ┌──────────────────────┐ ┌──────────────────────┐
│         July         │ │        August        │ │      September       │
├──────────────────────┤ ├──────────────────────┤ ├──────────────────────┤
│ Mo Tu We Th Fr Sa Su │ │ Mo Tu We Th Fr Sa Su │ │ Mo Tu We Th Fr Sa Su │
│  1  2  3  4  5  6  7 │ │  1  2  3  4  5  6  7 │ │  1  2  3  4  5  6  7 │
│  8  9 10 11 12 13 14 │ │  8  9 10 11 12 13 14 │ │  8  9 10 11 12 13 14 │
│ 15 16 17 18 19 20 21 │ │ 15 16 17 18 19 20 21 │ │ 15 16 17 18 19 20 21 │
│ 22 23 24 25 26 27 28 │ │ 22 23 24 25 26 27 28 │ │ 22 23 24 25 26 27 28 │
│                      │ │ 29 30 31 32 33 34 35 │ │                      │
└──────────────────────┘ └──────────────────────┘ └──────────────────────┘

┌──────────────────────┐ ┌──────────────────────┐ ┌──────────────────────┐
│       October        │ │       November       │ │       December       │
├──────────────────────┤ ├──────────────────────┤ ├──────────────────────┤
│ Mo Tu We Th Fr Sa Su │ │ Mo Tu We Th Fr Sa Su │ │ Mo Tu We Th Fr Sa Su │
│  1  2  3  4  5  6  7 │ │  1  2  3  4  5  6  7 │ │  1  2  3  4  5  6  7 │
│  8  9 10 11 12 13 14 │ │  8  9 10 11 12 13 14 │ │  8  9 10 11 12 13 14 │
│ 15 16 17 18 19 20 21 │ │ 15 16 17 18 19 20 21 │ │ 15 16 17 18 19 20 21 │
│ 22 23 24 25 26 27 28 │ │ 22 23 24 25 26 27 28 │ │ 22 23 24 25 26 27 28 │
│                      │ │ 29 30 31 32 33 34 35 │ │ 29 30 31 32 33 34 35 │
└──────────────────────┘ └──────────────────────┘ └──────────────────────┘
```
