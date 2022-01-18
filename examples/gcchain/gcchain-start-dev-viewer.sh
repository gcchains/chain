#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"

set -u
set -e

proj_dir=../..

echo "[*] Starting gcchain viewer nodes"
#set log level by add parameter:--verbosity 4
# or spec env like this:env CPC_VERBOSITY=4  ./gcchain-start.sh
#PanicLevel	0
#FatalLevel	1
#ErrorLevel	2
#WarnLevel	3
#InfoLevel	4
#DebugLevel	5


args="run --networkid 1 --rpcapi personal,eth,gcc,admission,net,web3,db,txpool,miner --linenumber --runmode dev "

gcchain=$proj_dir/build/bin/gcchain
ipc_path_base=data/gcc-

echo "start civilians unlock"
nohup $gcchain $args --ipcaddr ${ipc_path_base}11 --datadir data/data11  --rpcaddr 127.0.0.1:8511 --port 30321 \
         --unlock "0xbc131722d837b7d867212568baceb3a981181443"  --password conf-dev/passwords/password --logfile data/logs/11.log 2>data/logs/11.err.log &

echo "start civilians no unlock"
nohup $gcchain $args --ipcaddr ${ipc_path_base}12 --datadir data/data12  --rpcaddr 127.0.0.1:8512 --port 30322 \
    --logfile data/logs/12.log 2>data/logs/12.err.log &

echo "start bank node"
nohup $gcchain $args --ipcaddr ${ipc_path_base}22 --datadir data/data22  --rpcaddr 127.0.0.1:8522 --port 30332 \
         --unlock "0xabb528bffc707c2c507307e426ce810a7ad93ed6"  --password conf-dev/passwords/password --logfile data/logs/22.log 2>data/logs/22.err.log &


