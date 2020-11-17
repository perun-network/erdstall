#!/bin/bash

# SPDX-License-Identifier: Apache-2.0

set -e

# Download solc.
wget -nc "https://github.com/ethereum/solidity/releases/download/v0.7.4/solc-static-linux"
chmod +x solc-static-linux
echo -e "Ensure that the newest version of abigen is installed"

solpath="contracts"

# Generates optimized golang bindings and runtime binaries for sol contracts.
# $1  solidity file path, relative to $solpath/.
# $1  golang package name.
# $2… list of contract names.
function generate() {
    file=$1; pkg=$2
    shift; shift   # skip the first two args.
    for contract in "$@"; do
        abigen --pkg $pkg --sol $solpath/$file.sol --out $pkg/$file.go --solc ./solc-static-linux
        ./solc-static-linux --bin-runtime --optimize --allow-paths *, $solpath/$file.sol --overwrite -o $pkg/
        echo -e "package $pkg\n\n // ${contract}BinRuntime is the runtime part of the compiled bytecode used for deploying new contracts.\nvar ${contract}BinRuntime = \`$(<${pkg}/${contract}.bin-runtime)\`" > "$pkg/${contract}BinRuntime.go"
    done
}

# Generate bindings
generate "Erdstall" "bindings" "Erdstall"

abigen --version --solc ./solc-static-linux
echo -e "Generated bindings"
