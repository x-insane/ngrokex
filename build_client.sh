#!/bin/bash

go-bindata -nomemcopy -pkg=assets \
    -debug=false \
    -o=client/assets/assets_release.go \
    assets/client/...

go build -o output/ngrok/ngrok github.com/x-insane/ngrokex/main/ngrok

cp main/ngrok/config.yml output/ngrok/config.yml
