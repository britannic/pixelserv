# Changelog

## Releases

### Documentation update (Feb 2, 2020)

* Blacklist repository has been updated
* Requires apt sources update on OS

### Release v1.0.8 (Apr 19, 2019)

* Moved /etc/systemd/system/pixelserv.service to /lib/systemd/system/pixelserv.service
* Fix auto start of pixelserv following reboot by updating post install script with:

```bash
    systemctl enable pixelserv
```

### Release v1.0.7 (Apr 7, 2019)

* Updated to support EdgeOS 2.0.1

### Patch v1.0.4 (Jun 4, 2018)

* Changed detection logic for existing pseudo-ethernet interface from IP to actual device
* Added debian package to online repository

### Patch v1.0.3 (Jun 4, 2018)

* Documentation changes for EdgeOS v1.10.3

### Patch v1.0.2 (Mar 6, 2018)

* Added logic to use pppoe parent link interface for pseudo-ethernet device

### Patch v1.0.1 (Feb 4, 2018)

* Added EdgeRouter configuration settings to move "service gui http-port 80" to "service gui http-port 8180" to prevent conflict with pixelserv on port 80

### Initial Release (Dec 17, 2017)

* Initial release v1.0.0
