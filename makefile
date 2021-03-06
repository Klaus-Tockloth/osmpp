# makefile to cross compile the certstate tool / utility
# list all cross compile possibilities: go tool dist list
# 
# makefile adapted from this example:
# http://stackoverflow.com/documentation/go/1020/cross-compilation#t=201703051136361578518
#
# releases:
# - v1.0.0 - 2019/11/23: initial release

appname := osmpp
sources := $(wildcard *.go)

build = GOOS=$(1) GOARCH=$(2) go build -o build/$(appname)$(3)
tar = cd build && tar -cvzf $(appname)_$(1)_$(2).tar.gz $(appname)$(3) && rm $(appname)$(3)
zip = cd build && zip $(appname)_$(1)_$(2).zip $(appname)$(3) && rm $(appname)$(3)

.PHONY: all windows darwin linux freebsd openbsd solaris clean

all: linux darwin windows freebsd openbsd solaris

clean:
	rm -rf build/

# ----- linux builds -----
linux: build/$(appname)_linux_amd64.tar.gz build/$(appname)_linux_arm64.tar.gz

build/$(appname)_linux_amd64.tar.gz: $(sources)
	$(call build,linux,amd64,)
	$(call tar,linux,amd64)

build/$(appname)_linux_arm64.tar.gz: $(sources)
	$(call build,linux,arm64,)
	$(call tar,linux,arm64)

# ----- darwin (macOS) build -----
darwin: build/$(appname)_darwin_amd64.tar.gz

build/$(appname)_darwin_amd64.tar.gz: $(sources)
	$(call build,darwin,amd64,)
	$(call tar,darwin,amd64)

# ----- windows build -----
windows: build/$(appname)_windows_amd64.zip

build/$(appname)_windows_amd64.zip: $(sources)
	$(call build,windows,amd64,.exe)
	$(call zip,windows,amd64,.exe)

# ----- freebsd build -----
linux: build/$(appname)_freebsd_amd64.tar.gz

build/$(appname)_freebsd_amd64.tar.gz: $(sources)
	$(call build,freebsd,amd64,)
	$(call tar,freebsd,amd64)

# ----- openbsd build -----
linux: build/$(appname)_openbsd_amd64.tar.gz

build/$(appname)_openbsd_amd64.tar.gz: $(sources)
	$(call build,openbsd,amd64,)
	$(call tar,openbsd,amd64)

# ----- solaris build -----
linux: build/$(appname)_solaris_amd64.tar.gz

build/$(appname)_solaris_amd64.tar.gz: $(sources)
	$(call build,solaris,amd64,)
	$(call tar,solaris,amd64)
