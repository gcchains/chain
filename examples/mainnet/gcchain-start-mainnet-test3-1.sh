#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"

#set log level by add parameter:--verbosity 4
# or spec env like this:env CPC_VERBOSITY=4  ./gcchain-start-mainnet-test.sh
#PanicLevel	0
#FatalLevel	1
#ErrorLevel	2
#WarnLevel	3
#InfoLevel	4
#DebugLevel	5


set -u
set -e

validator_ip=""
if [ $# == 0 ]; then
    validator_ip='127.0.0.1'
else
    validator_ip="$1"
fi

source ./gcchain-start-mainnet-config.sh ${validator_ip}

./gcchain-start-mainnet-init.sh

./gcchain-start-mainnet-bootnode.sh

# TODO change address
./gcchain-start-mainnet-validator.sh

./gcchain-start-mainnet-proposer-1.sh

./gcchain-start-mainnet-proposer-2.sh

./gcchain-start-mainnet-proposer-3.sh

./gcchain-start-mainnet-civilian.sh

./gcchain-start-mainnet-contract-admin.sh

./gcchain-start-mainnet-deploy-contract.sh ${validator_ip}


