.PHONY: default server client deps fmt clean all release-all assets client-assets server-assets

BUILDTAGS=debug
default: all

deps: assets
	go get -tags '$(BUILDTAGS)' -d -v ./...

server: deps
	go install -tags '$(BUILDTAGS)' main/ngrokd

fmt:
	go fmt ./...

client: deps
	go install -tags '$(BUILDTAGS)' main/ngrok

assets: client-assets server-assets

bin/go-bindata:
	GOOS="" GOARCH="" go get github.com/jteeuwen/go-bindata/go-bindata

client-assets: bin/go-bindata
	go-bindata -nomemcopy -pkg=assets -tags=$(BUILDTAGS) \
		-debug=$(if $(findstring debug,$(BUILDTAGS)),true,false) \
		-o=client/assets/assets_$(BUILDTAGS).go \
		assets/client/...

server-assets: bin/go-bindata
	go-bindata -nomemcopy -pkg=assets -tags=$(BUILDTAGS) \
		-debug=$(if $(findstring debug,$(BUILDTAGS)),true,false) \
		-o=server/assets/assets_$(BUILDTAGS).go \
		assets/server/...

release-client: BUILDTAGS=release
release-client: client

release-server: BUILDTAGS=release
release-server: server

release-all: fmt release-client release-server

all: fmt client server

clean:
	go clean -i -r github.com/x-insane/ngrokex/...
	rm -rf client/assets/ server/assets/
