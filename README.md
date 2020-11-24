# UBNT edgeos-pixelserv Transparent Pixel Server for IP Blackhole Redirects

[![License](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/britannic/pixelserv/blob/master/LICENSE.txt) [![Alpha Version](https://img.shields.io/badge/version-v0.06-green.svg)](https://github.com/britannic/pixelserv) [![GoDoc](https://godoc.org/github.com/britannic/pixelserv?status.svg)](https://godoc.org/github.com/britannic/pixelserv) [![Build Status](https://travis-ci.org/britannic/pixelserv.svg?branch=master)](https://travis-ci.org/britannic/pixelserv) [![Coverage Status](https://coveralls.io/repos/github/britannic/pixelserv/badge.svg?branch=master)](https://coveralls.io/github/britannic/pixelserv?branch=master) [![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/github.com/britannic/pixelserv)

[community.ubnt.com](https://community.ubnt.com/t5/EdgeMAX/Self-Installer-to-configure-Ad-Server-and-pixelserv-Blocking/td-p/1337892)

NOTE: THIS IS NOT OFFICIAL UBIQUITI SOFTWARE AND THEREFORE NOT SUPPORTED OR ENDORSED BY Ubiquiti NetworksÂ®

## Donations and Sponsorship

Please show your thanks by donating to the project using [Square Cash](https://cash.me/$HelmRockSecurity/ "Securely send and receive cash without fees using Square Cash") or [PayPal](https://www.paypal.me/helmrocksecurity/)

[![Donate](https://img.shields.io/badge/Donate-%245-orange.svg?style=plastic)](https://cash.me/$HelmRockSecurity/5 "Give $5 using Square Cash (free money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2410-red.svg?style=plastic)](https://cash.me/$HelmRockSecurity/10 "Give $10 using Square Cash (free money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2415-yellow.svg?style=plastic)](https://cash.me/$HelmRockSecurity/15 "Give $15 using Square Cash (free money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2420-yellowgreen.svg?style=plastic)](https://cash.me/$HelmRockSecurity/20 "Give $20 using Square Cash (free money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2450-ff69b4.svg?style=plastic)](https://cash.me/$HelmRockSecurity/50 "Give $50 using Square Cash (free money transfer)")
[![Donate](https://img.shields.io/badge/Donate-Custom%20Amount-4B0082.svg?style=plastic)](https://cash.me/$HelmRockSecurity/ "Choose your own donation amount using Square Cash (free money transfer)")

[![Donate](https://img.shields.io/badge/Donate-%245-orange.svg?style=plastic)](https://paypal.me/helmrocksecurity/5 "Give $5 using PayPal (PayPal money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2410-red.svg?style=plastic)](https://paypal.me/helmrocksecurity/10 "Give $10 using PayPal (PayPal money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2415-yellow.svg?style=plastic)](https://paypal.me/helmrocksecurity/15 "Give $15 using PayPal (PayPal money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2420-yellowgreen.svg?style=plastic)](https://paypal.me/helmrocksecurity/20 "Give $20 using PayPal (PayPal money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2450-ff69b4.svg?style=plastic)](https://paypal.me/helmrocksecurity/50 "Give $50 using PayPal (PayPal money transfer)")
[![Donate](https://img.shields.io/badge/Donate-Custom%20Amount-4B0082.svg?style=plastic)](https://paypal.me/helmrocksecurity/ "Choose your own donation amount using PayPal (PayPal money transfer)")

We greatly appreciate any and all donations - Thank you! Funds go to maintaining development servers and networks.

## Copyright

* Copyright © [2020 Helm Rock Consulting](https://www.helmrock.com/ "Visit Helm Rock Consulting at https://www.helmrock.com/")

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

* Prevents HTTP 404 page not found messages if used in conjunction with dnsmasq IP redirects ([edgeos-dnsmasq-blacklist]( https://britannic.github.io/blacklist/) provides dnsmasq redirection and blacklisting functionality)

## Compatibility

* edgeos-pixelserv has been tested on the EdgeRouter **ERLite-3**, **ERPoe-5**, **ER-X**: EdgeOS versions **v1.7.0-v2.0.9**
* Note: the debian package will not successfully install on a UniFi Gateway, since there is also a default HTTP port 80 listener configured all interfaces

## **Change Log**

* See [changelog](CHANGELOG.md) for details

## Installation

* edgeos-pixelserv installs itself as a service into /etc/init.d/pixelserv
* The installation will modify the router's configuration settings to move "service gui http-port 80" to "service gui http-port 8180" to prevent conflict with pixelserv on port 80

* [Using apt-get](https://github.com/britannic/pixelserv#apt-get-installation---erlite-3-erpoe-5-er-x--er-x-sfp) - works for all routers
* [Using dpkg](#dpkg-installation---best-for-disk-space-constrained-routers) - best for disk space constrained routers

## apt-get Installation - ERLite-3, ERPoe-5, ER-X & ER-X-SFP

* Add the blacklist debian package repository using the router's CLI shell

```bash
configure
set system package repository blacklist components main
set system package repository blacklist description 'Britannic blacklist debian stretch repository'
set system package repository blacklist distribution stretch
set system package repository blacklist url 'https://raw.githubusercontent.com/britannic/debian-repo/master/blacklist/public/'
commit;save;exit
```

* Add the GPG signing key

```bash
sudo curl -L https://raw.githubusercontent.com/britannic/debian-repo/master/blacklist/public.key | sudo apt-key add -
```

* Update the system repositorities and install edgeos-pixelserv

```bash
sudo apt-get update && sudo apt-get install edgeos-pixelserv
```

## dpkg installation - Best for disk space constrained routers

### EdgeRouter ERLite-3, ERPoe-5 and similar MIPS based Edgerouters

```bash
curl -L -O  https://raw.githubusercontent.com/britannic/pixelserv/master/edgeos-pixelserv_1.0.8_mips.deb
sudo dpkg -i edgeos-pixelserv_1.0.8_mips.deb
```

### EdgeRouter ER-X & ER-X-SFP

```bash
curl -L -O  https://raw.githubusercontent.com/britannic/pixelserv/master/edgeos-pixelserv_1.0.8_mipsel.deb
sudo dpkg -i edgeos-pixelserv_1.0.8_mipsel.deb
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
Usage: pixelserv [options]

  -f <file>
        load pixel or other content from <file> source
  -h    Display help
  -ip string
        IP address for pixelserv.amd64 to bind to (default "127.0.0.1")
  -path string
        Set HTTP root path (default "/")
  -port string
        Port number for pixelserv.amd64 to listen on (default "80")
  -version
        Show version
```

* pixelserv.sysv service (EdgeOS < v1.10.x)

```bash
service pixelserv {start|stop|status|restart|force-reload|reload}
```

* pixelserv service (EdgeOS > v2.0.1)

```bash
systemctl start pixelserv
```# pixelserv
--
