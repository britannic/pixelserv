 # Makefile to build dnsmasq blacklist
 SHELL=/usr/bin/env bash

 # Go parameters
	EXECUTABLE=pixelserv
	GOBUILD=$(GOCMD) build
	GOCLEAN=$(GOCMD) clean
	GOCMD=go
	GOGET=$(GOCMD) get
	GOTEST=$(GOCMD) test

# Executables
	GSED=$(shell which gsed || which sed) -i.bak -e
	PKG=edgeos-$(EXECUTABLE)

# Environment variables
	AWS=aws
	COPYRIGHT=s/Copyright © 20../Copyright © $(shell date +"%Y")/g
	COVERALLS_TOKEN=W6VHc8ZFpwbfTzT3xoluEWbKkrsKT1w25
	DATE=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
	GIT=$(shell git rev-parse --short HEAD)
	LIC=LICENSE
	PAYLOAD=./.payload
	README=README.md
	READMEHDR=README.header
	SCRIPTS=/config/scripts
	OLDVER=$(shell cat ./OLDVERSION)
	VER=$(shell cat ./VERSION)
	VERSIONS=s/edgeos-$(EXECUTABLE)_$(OLDVER)_/edgeos-$(EXECUTABLE)_$(VER)_/g
	BADGE=s/version-v$(OLDVER)-green.svg/version-v$(VER)-green.svg/g
	TAG="v$(VER)"

.PHONY: all clean deps mips coverage copyright docs readme
all: clean deps mips coverage copyright docs readme

.PHONY: amd64 
amd64:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(EXECUTABLE).amd64 \
	-ldflags \
	"-X main.build=$(DATE) \
	-X main.githash=$(GIT) \
	-X main.version=$(VER) \
	-s -w" -v

.PHONY: build
build: clean amd64 mips copyright docs readme 

.PHONY: cdeps 
cdeps: 
	dep status -dot | dot -T png | open -f -a /Applications/Preview.app

.PHONY: clean
clean:
	$(GOCLEAN)
	find . -name "$(EXECUTABLE).*" -type f \
	-o -name debug -type f \
	-o -name "*.deb" -type f \
	-o -name debug.test -type f \
	-o -name "*.tgz" -type f \
	| xargs rm 

.PHONY: copyright
copyright:
	$(GSED) '$(COPYRIGHT)' $(README)
	$(GSED) '$(COPYRIGHT)' $(LIC)

.PHONY: coverage 
coverage: 
	./testcoverage

.PHONY: dep-stat 
dep-stat: 
	dep status

.PHONY: deps
deps: 
	dep ensure -update

.PHONY: docs
docs: version readme 

.PHONY: mips
mips: mips64 mipsle

.PHONY: mips64
mips64:
	GOOS=linux GOARCH=mips64 $(GOBUILD) -o $(EXECUTABLE).mips \
	-ldflags \
	"-X main.build=$(DATE) \
	-X main.githash=$(GIT) \
	-X main.version=$(VER) \
	-s -w" -v 

.PHONY: mipsle
mipsle:
	GOOS=linux GOARCH=mipsle $(GOBUILD) -o $(EXECUTABLE).mipsel \
	-ldflags \
	"-X main.build=$(DATE) \
	-X main.githash=$(GIT) \
	-X main.version=$(VER) \
	-s -w" -v

.PHONY: pkgs
pkgs: pkg-mips pkg-mipsel 

.PHONY: pkg-mips 
pkg-mips: clean deps mips coverage copyright docs readme
	cp $(EXECUTABLE).mips $(PAYLOAD)$(SCRIPTS)/$(EXECUTABLE) \
	&& ./make_deb $(EXECUTABLE) mips

.PHONY: pkg-mipsel
pkg-mipsel: clean deps mipsle coverage copyright docs readme
	cp $(EXECUTABLE).mipsel $(PAYLOAD)$(SCRIPTS)/$(EXECUTABLE) \
	&& ./make_deb $(EXECUTABLE) mipsel

.PHONY: readme 
readme: version
	cat README.header > README.md 
	godoc2md github.com/britannic/$(EXECUTABLE) >> README.md

.PHONY: tags
tags:
	git push origin --tags

.PHONY: version
version:
	$(GSED) '$(BADGE)' $(READMEHDR)
	$(GSED) '$(VERSIONS)' $(READMEHDR)
	cat ./VERSION > ./OLDVERSION

.PHONY: release
release: all commit push
	@echo Released $(TAG)

.PHONY: commit
commit:
	@echo Committing release $(TAG)
	git commit -am"Release $(TAG)"
	git tag $(TAG)

.PHONY: push
push:
	@echo Pushing release $(TAG) to master
	git push --tags
	git push

.PHONY: repo
repo:
	@echo Pushing repository $(TAG) to $(AWS)
	scp $(PKG)_$(VER)_*.deb $(AWS):/tmp
	./aws.sh $(AWS) $(PKG)_$(VER)_ $(TAG)

.PHONY: upload
upload: pkgs
	scp $(PKG)_$(VER)_mips.deb dev1:/tmp
	scp $(PKG)_$(VER)_mipsel.deb er-x:/tmp
	scp $(PKG)_$(VER)_mips.deb ubnt:/tmp