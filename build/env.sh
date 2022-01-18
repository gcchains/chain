#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi


# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
gcchaindir="$workspace/src//gcchain"
if [ ! -L "$gcchaindir/chain" ]; then
    mkdir -p "$gcchaindir"
    cd "$gcchaindir"
    ln -s ../../../../../. chain
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$gcchaindir/chain"
PWD="$gcchaindir/chain"

# Launch the arguments with the configured environment.
exec "$@"
