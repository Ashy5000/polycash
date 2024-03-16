# ashy5000/cryptocurrency

![](https://github.com/Ashy5000/cryptocurrency/actions/workflows/go.yml/badge.svg)

The home of a power-efficient, scalable, and modern cryptocurrency blockchain designed to resolve the issues presented by the traditional PoW (Proof of Work) incentive mechanism while maintaining decentralization. Written in Golang and Rust for secure algorithm implementations, lightweight networking, memory safety, speed, and reliability.

## Modified PoW protocol
As opposed to the original Proof of Work (PoW) protocol, this blockchain does not calculate block difficulty on a global level. Instead, it is adjusted at a per-miner basis, removing the reward for using more computing power and thus more energy. If a miner has a higher hash rate and thus initially mines blocks faster, they will be provided with higher difficulty blocks and will end up mining blocks at a **constant rate** of *1 min/transaction*. In addition, the block difficulty is retargeted *every* time a miner mines a block, ensuring quick reactions to changes in miner hash rates. To prevent miners from registering multiple private keys in an attempt to avoid the difficulty constraints, there is planned to be a high initial difficulty, or the constant difficulty for the first block a miner mines. In addition, the total number of miners at a time is limited to a maximum that increases as the blockchain's length grows. This prevents miners from using large amounts of computing resources and energy to create a large number of key pairs that have overcome the initial block.

## Roadmap
The roadmap for this repository looks something like this:

**Development**: Implementing a blockchain and networking system, along with the specific changes to the PoW and networking protocols in order to increase the effectiveness, speed, reliability, and efficiency of the blockchain. Creating a README that documents the design decisions of the blockchain and network and how to use the software. Licensing the software with the GNU General Purpose License v3.0. Performing small scale, local tests on various operating systems and architectures to ensure the software works correctly.
**Testnet**: Running a testing network to verfify the correct operation of the software. Expirementing with its ability to verify large amounts of transactions at once. Making adjustments if neccessary to increase the speed, reliability, and efficiency of the blockchain. Running small-scale tests on the cryptocurrency's financial model.
**Production**: Launching the final product, fixing any issues if and when they arrive.

## Directory structure
The root directory of this project is occupied by the Golang source code that nodes run in order to interact with each other and the decentralized blockchain. In the ```peer_server``` directory, there is Rust code that can be run by servers to maintain a list of peers in the network. Nodes can connect to these servers or maintain their own lists. You will probably only need to run the Golang code. There is also a build directory that contains builds of the node and peer server for various platforms. By launch, almost all of the major platforms will be supported. Alternatively, you may build from source.

## About decentralized peer lists
There are two ways to run a node: using a peer server or a local peer list. With a peer server, there is far less configuration. To make yourself known to the network, use the `addpeer [YOUR IP]` command in the BlockCMD console (see `# To run as a client:` below to open the console). With a local peer list, you fully decentralize your connections with other computers to completely remove trust from the system. However, dealing with the configuration for this option can take a bit longer. I personally recommend using a local peer list. The security and reliability that decentralized peer tracking provides simply can't be reproduced with a centralized peer server, but if it seems daunting to manually configure IP addresses, the peer server's an okay way to start. The two options both use the same network and the same blockchain, so your balance will be preserved if you switch over.

## Setup & Usage

*Note: A node must have a publicly accessible IP address in order to join the network.*

Create and clone the repository:

```bash
git clone https://github.com/ashy5000/cryptocurrency
cd cryptocurrency
```

### To run as a client:
To use an interactive console for viewing and adding to the blockchain, run:
```bash
./builds/node/node_linux_x86_64 # replace for your os and architecture
````
Commands:
- `help`: see a list of all commands
- `sync`: update the blockchain and all balances and transactions
- `keygen`: generate a key pair so you can send and receive tokens
- `send {recipient} {amount}`: send {amount} tokens to {recipient}
- `balance {key}`: get the balance associated with the public key {key}
- `savestate`: save a backup of the current state of the blockchain to a file
- `loadstate`: load a backup of the current state of the blockchain from a file
- `exit`: exit the console
- `addpeer {ip}`: by default, connects to a peer. If using a centralized peer server, makes yourself known to the network.

### To run a node:
To run the node software, which keeps the blockchain distributed across the p2p network, run:
```bash
./builds/node/node_linux_x86_64 -serve -port 8080  # replace for your os and architecture
# In a new terminal window: (optional, starts a peer server so you it is faster to find new nodes)
# This is not at all required.
cd peer_server
cargo run
```


### To run a miner:
To run the mining software, which adds new blocks to the blockchain in exchange for a reward, run:
```bash
./builds/node/node_linux_x86_64 -serve -mine -port 8080
# In a new terminal window: (optional, starts a peer server so you it is faster to find new nodes)
# This is not at all required.
cd peer_server
cargo run
```

### To connect to a peer via their peer server:
To connect to a peer that is also running a peer server (ran the commands after `# In a new terminal window:`), run:
```bash
curl http://[PEER IP]:6060 -d 'http://[YOUR IP]:8080'
```

## License
This software is released under the GNU General Public License v3.0.

## Contributing
I'm currently not accepting pull requests for this repository, but you're more than welcome to open an issue if there's something you want improved, added, or fixed.
