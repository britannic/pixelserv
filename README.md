
# pixelserv

© 2017 Helmrock Consulting. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE.txt file.

[![License](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/britannic/pixelserv/blob/master/LICENSE.txt) [![Alpha  Version](https://img.shields.io/badge/version-v1.0.0-red.svg)](https://github.com/britannic/pixelserv) [![Build Status](https://travis-ci.org/britannic/pixelserv.svg?branch=master)](https://travis-ci.org/britannic/pixelserv) [![Coverage Status](https://coveralls.io/repos/github/britannic/pixelserv/badge.svg?branch=master)](https://coveralls.io/github/britannic/pixelserv?branch=master) [![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/github.com/britannic/pixelserv)

### Overview
pixelserv is a webserver that simply returns a single transparent pixel or content loaded from file

### Features
Prevents HTTP 404 page not found messages if used in conjunction with dnsmasq blacklisted IP redirects


### Compatibility
If [Go](https://golang.org) supports your platform's CPU architecture and OS, the source should compile and work with little to no modification in most cases

### Installation Notes:

Install it as a unix service, or run it manually in the background and use it as your blackhole/dns redirect address for any blacklisted advert servers

### Usage

	Usage: pixelserv [options]

	  -f <file>
	    	load pixel or other content from <file> source
	  -h	Display help
	  -ip string
	    	IP address for pixelserv to bind to (default "127.0.0.1")
	  -path string
	    	Set HTTP root path (default "/")
	  -port string
	    	Port number for pixelserv to listen on (default "80")
	  -version
	    	Show version


- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)