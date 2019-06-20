#!/bin/bash

go-bindata -nomemcopy -pkg=assets \
    -debug=false \
    -o=server/assets/assets_release.go \
    assets/server/...

go build -o output/ngrokd/ngrokd github.com/x-insane/ngrokex/main/ngrokd

cp main/ngrokd/config.yml output/ngrokd/config.yml
