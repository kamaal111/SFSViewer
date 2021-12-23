#!/bin/sh

#  run.sh
#  SFSViewer
#
#  Created by Kamaal M Farah on 23/12/2021.
#  

set -e

GO_BIN=/usr/local/go/bin/go

if ! which "$GO_BIN" > /dev/null; then
    echo "error: go is not installed. Vistit https://go.dev to learn more."
    exit 1
else
    cd Scripts/reformat-names
    "$GO_BIN" run *.go
fi
