# flightaware-go 

## Description

This is the library implementing the [Flightaware](http://www.flightaware.com/)  API in Go.  It just exports the streaming data into JSON, using the FA API.

Used by [FA Export](https://github.com/keltia/fa-export) and [FA Tail](https://github.com/keltia/fa-tail) utilities.

## Build status

[![GitHub release](https://img.shields.io/github/release/keltia/flightaware-go.svg)](https://github.com/keltia/flightaware-go/releases)
[![GitHub issues](https://img.shields.io/github/issues/keltia/flightaware-go.svg)](https://github.com/keltia/flightaware-go/issues)
[![Go Version](https://img.shields.io/badge/go-1.10-blue.svg)](https://golang.org/dl/)
[![GoDoc](https://godoc.org/github.com/keltia/flightaware-go?status.svg)](http://godoc.org/github.com/keltia/flightaware-go)
[![SemVer](http://img.shields.io/SemVer/2.0.0.png)](https://semver.org/spec/v2.0.0.html)
[![License](https://img.shields.io/pypi/l/Django.svg)](https://opensource.org/licenses/BSD-2-Clause)
[![Go Report Card](https://goreportcard.com/badge/github.com/keltia/flightaware-go)](https://goreportcard.com/report/github.com/keltia/flightaware-go)

Branch: develop — [![develop|Build Status](https://travis-ci.org/keltia/flightaware-go.svg?branch=develop)](http://travis-ci.org/keltia/flightaware-go/tree/develop)
Branch: master — [![master|Build Status](https://travis-ci.org/keltia/flightaware-go.svg?branch=master)](http://travis-ci.org/keltia/flightaware-go)

## Requirements

* Go >= 1.10

`github.com/keltia/flightaware-go` is a Go module (you can use either Go 1.10 with `vgo` or 1.11+).  The API exposed follows the Semantic Versioning scheme to guarantee a consistent API compatibility.

## Installation

Like many Go utilities & libraries, it is very easy to install:

    go get ithub.com/keltia/flightaware-go/cmd/...

that way, you get both the library and its two bundled utilities.

You can also clone the repository and use `make install`

    git clone https://github.com/keltia/flightaware-go
    make install

## USAGE

There are two example programs included in `cmd/fa-export` and `cmd/fa-tail`.  The former is the main driver and the latter is a `tail(1)`-like utility.

### fa-export

```
fa-export -[AOpv] [options...]

Usage of fa-export:
  -A	Autorotate output file
  -B string
    	Begin time for -f pitr|range
  -D string
    	Default destination (NOT IMPL)
  -E string
    	End time for -f range
  -F string
    	Airline filter
  -I string
    	Aircraft Ident filter
  -L string
    	Lat/Long filter
  -O	Overwrite existing file?
  -P string
    	Airport filter
  -X string
    	Hexid output filter
  -d string
    	Stop after N s/mn/h/days
  -e string
    	Events to stream
  -f string
    	Specify which feed we want (default "live")
  -o string
    	Specify output FILE.
  -p	Enable profiling
  -u string
    	Username to connect with
  -v	Set verbose flag.
```

### fa-tail

```
fa-tail -[cv] file

Usage of fa-tail:
  -c	Count records.
  -v	Be verbose
```

The `file` parameter being the file specified by the `-o`option of `fa-export`.

XXX `fa-tail` does not implement most of `tail(1)`options, especially not `-f`.

## API Usage

 You start by creating a client instance with your credentials passed as Config
 struct. See `fa-export` for a configuration file loading and suff.

 	client := flightaware.NewClient(Config)

 Then you can configure the feed type with

 	client.SetFeed(string, []time.Time)

 You can also set a timeout time with a value in seconds

 	client.SetTimeout(int64)

 You can add one or more different input filters:

    client.AddInputFilter(<type>, <value>)

 where type can be one of

     FILTER_EVENT
     FILTER_AIRLINE
     FILTER_IDENT
     FILTER_AIRPORT
     FILTER_LATLONG

 The filters you specify will be checked remotely by FlightAware according to the
 documentation available at
 https://fr.flightaware.com/commercial/firehose/firehose_documentation.rvt

 You can specify output filters with using `client.AddOutputFilter(string)`

 The default handler is to display all packets.  You can change the default handler
 with

 	client.AddHandler(func([]byte)

 Last action is to start the consuming/producer loop with

 	client.Start()

 Reading will be closed either though getting an EOF from FA or being will killed either manually or through the timeout value.

 You can then use

 	client.Close()

 to properly close the reading channel.

## License

The [BSD 2-Clause license](https://github.com/keltia/flightaware-go/LICENSE.md).

# Contributing

This project is an open Open Source project, please read `CONTRIBUTING.md`.

# Feedback

We welcome pull requests, bug fixes and issue reports.

Before proposing a large change, first please discuss your change by raising an issue.
