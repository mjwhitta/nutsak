# Network Utility Swiss-Army Knife (NUtSAK)

<a href="https://www.buymeacoffee.com/mjwhitta">🍪 Buy me a cookie</a>

[![Go Report Card](https://goreportcard.com/badge/gitlab.com/mjwhitta/nutsak)](https://goreportcard.com/report/gitlab.com/mjwhitta/nutsak)

## What is this?

NUtSAK is a network utility library. Included is a client offering
similar features to socat.

## How to install

Open a terminal and run the following:

```
$ # For library usage
$ go get --ldflags "-s -w" --trimpath -u gitlab.com/mjwhitta/nutsak
$ # For cli usage
$ go install --ldflags "-s -w" --trimpath \
    gitlab.com/mjwhitta/nutsak/cmd/sak@latest
```

Or compile from source:

```
$ git clone https://gitlab.com/mjwhitta/nutsak.git
$ cd nutsak
$ git submodule update --init
$ make
```

## Usage

### CLI

To create a simple TCP listner:

```
$ sak tcp-listen:4444,fork stdout
```

To create a simple TCP client:

```
$ sak stdin tcp:4444
```

### Library

```
package main

import (
    "time"

    sak "gitlab.com/mjwhitta/nutsak"
)

func main() {
    var a sak.NUt
    var b sak.NUt
    var e error

    // Create first NUt
    if a, e = sak.NewNUt("tcp-listen:4444,fork"); e != nil {
        panic(e)
    }

    // Create second NUt
    if b, e = sak.NewNUt("stdout"); e != nil {
        panic(e)
    }

    // Shutdown in 10 seconds
    go func() {
        time.Sleep(10 * time.Second)
        a.Down()
        b.Down()
    }()

    // Pair NUts to create two-way tunnel
    if e = sak.Pair(a, b); e != nil {
        panic(e)
    }
}
```

## Links

- [Source](https://gitlab.com/mjwhitta/nutsak)

## TODO

- Improve TCP/TLS fork option
- EXEC process support
- HTTP proxy support
