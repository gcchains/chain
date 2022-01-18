It contains a simple example of gcchain that starts 4 nodes and issues a transaction from node 1 to
node 2.  node1 and node2 are the initial signers.

# Usage

- `gcchain-init.sh` initializes accounts and keystore.
- `gcchain-start.sh` launches 4 gcchain nodes, and node 1 and 2 are mining. logs of all nodes are printed in `data/logs`.
- `gcchain-stop.sh` stops all gcchain nodes.

# testing simple transaction
please install the deps, see install-deps.sh
check out transactions/*.py