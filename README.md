# UBNT edgeos-pixelserv Transparent Pixel Server for IP Blackhole Redirects

[![License](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/britannic/pixelserv/blob/master/LICENSE.txt) [![Alpha  Version](https://img.shields.io/badge/version-v0.06-green.svg)](https://github.com/britannic/pixelserv) [![GoDoc](https://godoc.org/github.com/britannic/pixelserv?status.svg)](https://godoc.org/github.com/britannic/pixelserv) [![Build Status](https://travis-ci.org/britannic/pixelserv.svg?branch=master)](https://travis-ci.org/britannic/pixelserv) [![Coverage Status](https://coveralls.io/repos/github/britannic/pixelserv/badge.svg?branch=master)](https://coveralls.io/github/britannic/pixelserv?branch=master) [![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/github.com/britannic/pixelserv)

[community.ubnt.com](https://community.ubnt.com/t5/EdgeMAX/Self-Installer-to-configure-Ad-Server-and-pixelserv-Blocking/td-p/1337892)

NOTE: THIS IS NOT OFFICIAL UBIQUITI SOFTWARE AND THEREFORE NOT SUPPORTED OR ENDORSED BY Ubiquiti Networks®

## Copyright

* Copyright © 2018 Helm Rock Consulting

## Overview

pixelserv is a simple webserver that returns a single transparent pixel or content loaded from a file

## Licenses

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.
1. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

The views and conclusions contained in the software and documentation are those
of the authors and should not be interpreted as representing official policies,
either expressed or implied, of the FreeBSD Project.

## Features

* Prevents HTTP 404 page not found messages if used in conjunction with dnsmasq IP redirects

## Compatibility

* edgeos-pixelserv has been tested on the EdgeRouter ERLite-3, ERPoe-5, ER-X: EdgeOS versions v1.7.0-v1.9.7+hotfix.4
* Note: the debian package will not successfully install on a UniFi Gateway, since there is also a default HTTP port 80 listener configured all interfaces

## **Change Log**

* See [changelog](CHANGELOG.md) for details

## Installation

* edgeos-pixelserv installs itself as a service into /etc/init.d/pixelserv
* The installation will modify the router's configuration settings to move "service gui http-port 80" to "service gui http-port 8180" to prevent conflict with pixelserv on port 80

### EdgeRouter ERLite-3, ERPoe-5 and similar MIPS based Edgerouters

```bash
curl https://community.ubnt.com/ubnt/attachments/ubnt/EdgeMAX/195918/1/edgeos-pixelserv_1.0.1_mips.deb.tgz | tar -xvz
sudo dpkg -i edgeos-pixelserv_1.0.1_mips.deb
```

### EdgeRouter ER-X & ER-X-SFP

```bash
curl https://community.ubnt.com/ubnt/attachments/ubnt/EdgeMAX/195918/2/edgeos-pixelserv_1.0.1_mipsel.deb.tgz | tar -xvz
sudo dpkg -i edgeos-pixelserv_1.0.1_mipsel.deb
```

## Removal

* Removal will modify the router's configuration settings to move "service gui http-port 8180" back to the default "service gui http-port 80"

### EdgeMAX ERLite-x & EdgeMax ER-X

```bash
sudo apt-get remove edgeos-pixelserv
```

### Usage

* Standalone binary
 
```bash
/config/scripts/pixelserv -h
```

* pixelserv service

```bash
service pixelserv {start|stop|status|restart|force-reload|reload}
```