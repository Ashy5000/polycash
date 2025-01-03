# Polycash

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Rust](https://img.shields.io/badge/rust-%23000000.svg?style=for-the-badge&logo=rust&logoColor=white)
![C++](https://img.shields.io/badge/c++-%2300599C.svg?style=for-the-badge&logo=c%2B%2B&logoColor=white)


<img src="https://github.com/Ashy5000/polycash/blob/master/assets/logos/logo_new.png?raw=true" alt="Polycash logo" width="200" />

[Radicle](https://radicle.xyz) ID: `rad:z2E1qt896CsMnTcRRhqt8sQDCgKp`

The home of a power-efficient, secure, quantum-resistant, and modern blockchain and smart contract platform designed to resolve the issues presented by existing consensus mechanisms while improving upon decentralization. Written in Golang, Rust, and C++.

**Polycash is now ZK!** All state changes are now verified and sealed using SNARKs. This greatly improves the scalability potential for Polycash.

For the whitepaper, look [here](https://github.com/Ashy5000/polycash/blob/master/docs/whitepaper/whitepaper.pdf).

Technical docs are located in the [wiki](https://github.com/Ashy5000/polycash/wiki).

>[!NOTE]
>Polycash has no association with the Polygon blockchain, despite the similar naming.

**Comparison**

|                    | Bitcoin | Ethereum        | Polycash          |
| ------------------ | ------- | --------------- | ----------------- |
| Consensus          | PoW     | PoS             | APoW              |
| Scriptability      | Limited | Turing-complete | Turing-complete   |
| VM                 | N/A     | EVM             | PVM               |
| Signatures         | ECC     | ECC             | Lattice-based     |
| Token              | BTC     | ETH             | PCSH              |
| Pre-mined tokens   | No      | Yes             | No                |
| Builtin Oracles    | No      | No              | Yes               |
| Hash function      | SHA256  | Keccak-256      | SHA3-512 & SHA256 |
| ZK                 | N/A     | N/A             | SNARK             |


>[!TIP]
>[Discussions](https://github.com/Ashy5000/cryptocurrency/discussions) are now available! You can use discussions for asking questions, giving feedback, or really anything else related to the project.

## Table of Contents
1. [Modified PoW protocol](./#Modified-PoW-protocol)
2. [Time Verification](./#Time-Verification)
3. [Quantum-Resistant Signatures](./#Quantum-Resistant-Signatures)
4. [Roadmap](./#Roadmap)
5. [Directory Structure](./#Directory-Structure)
6. [Setup & Usage](./#Setup--Usage)
7. [License](./#License)

## Modified PoW protocol

As opposed to the original Proof of Work (PoW) protocol, Polycash does not calculate block difficulty on a global level. Instead, it is adjusted at a per-miner basis, mostly removing the reward for using more computing power and thus more energy. If a miner has a higher hash rate and thus initially mines blocks faster, they will be provided with higher difficulty blocks and will end up mining blocks at a **nearly constant rate** of _1 min/transaction_. The reason this is not exactly constant is that miners with a higher mining rate do receive slightly lower difficulties than necessary for the constant rate. At infinite computing power, a miner would theoretically have a 50% lower difficulty than expected: this is the lower boundary. A miner with zero computing power would theoretically have a 50% higher difficulty than expected: this is the upper boundary. Keep in mind that the difficulty is still lower for miners with a lower mining rate, but the ratio of difficulty to mining rate is higher with a lower mining rate. This motivates miners to contribute computing power to the network to keep it secure while still maintaining an increase in decentralization. In addition, the block difficulty is retargeted _every_ time a miner mines a block, ensuring quick reactions to changes in miner hash rates. To prevent miners from registering multiple private keys in an attempt to avoid the difficulty constraints, there are several safeguards put in place. First, miners are not rewarded for the first 3 blocks they mine, so creating multiple keys would cost a significant barrier to entry. Also, the block reward decreases by 1% for each new mining keypair that joins the network, so creating multiple keys becomes unprofitable, as it would damage the revenue of existing keys.

## Time Verification

Using a new verification method, time verification, both the security and the performance of the blockchain are improved as compared to Proof of Work. This protocol can prevent malicious forks from occuring using nodes' verification of new blocks using the curent timestamp, requiring signatures from miner nodes to become valid. Setting rewards and limits for the number of valid signatures in a block strictly enforces the security of this protocol.

## Quantum-Resistant Signatures

This blockchain utilizes the Dilithium3 signature algorithm, a quantum-resistant algorithm chosen as a winner for the NIST Post Quantum Cryptography standardization process. This is designed to future-proof the blockchain from the developing field of quantum computing.

## ZK

Every smart contract invocation in a given block can be verified with a SNARK proof, with only one necessary per block. This greatly improves verification times for large batches of transactions.

## Roadmap

The roadmap for this repository looks something like this:

**Initial Development** (COMPLETE): Implementing a blockchain and networking system, along with the specific changes to the PoW and networking protocols in order to increase the effectiveness, speed, reliability, and efficiency of the blockchain. Creating a README that documents the design decisions of the blockchain and network and how to use the software. Licensing the software with the GNU General Purpose License v3.0. Performing small scale, local tests on various operating systems and architectures to ensure the software works correctly.

**Testnet** (DEPRECATED): Running a testing network to verfify the correct operation of the software. Expirementing with its ability to verify large amounts of transactions at once. Making adjustments if neccessary to increase the speed, reliability, and efficiency of the blockchain. Running small-scale tests on the cryptocurrency's financial model.

**ZK Refactor** (ACTIVE): Refactoring the entire Polycash blockchain to use a zero-knowledge proof system with the goal of enhancing scalability, decentralization, and security.

**ZK Testnet** (COMING UP): Launching a post-refactor testnet that will replace the first testnet. **The bootstrap node will switch to this network momentarily.**

**Mainnet** (COMING UP): Launching the final network, fixing any issues if and when they arrive. Uses the same protocol as the ZK testnet.

## Directory structure

The root directory of this project is occupied by the Golang source code that nodes run in order to interact with each other and the decentralized blockchain. The `gui_wallet` directory contains Rust code for a GUI wallet to make transactions easier to make. Note that you will still have to generate a key using the `keygen` command in the interactive console (see _To run as a client_ below). The `vm` directory contains a Rust project capable of running smart contracts on the blockchain's virtual machine. The CMake project in the `language` directory forms the compiler for Polylang, a high-level language that compiles to `.blockasm` files, which can be run on the blockchain's virtual machine.

## Setup & Usage
>[!IMPORTANT]
>A node must have a publicly accessible IP address in order to join the network. You may have to set up port forwarding on your router.

Prerequisites:
- [go](https://go.dev/doc/install)
- [liboqs-go](https://github.com/open-quantum-safe/liboqs-go#installation)
- [rust](https://www.rust-lang.org/tools/install)

Run the install script:
```bash
curl https://raw.githubusercontent.com/Ashy5000/cryptocurrency/master/install.sh | bash
```

### To run as a client:

To use an interactive console for viewing and adding to the blockchain, run:

```bash
./builds/node/node
```

Commands:

- `help`: see a list of all commands
- `license`: display the software's license (GNU GPL v3)
- `sync`: update the blockchain and all balances and transactions
- `keygen`: generate a key pair so you can send and receive tokens
- `showPublicKey`: print your public key to give to people or services that need to pay you
- `encrypt`: encrypt the private key so you can store it safely
- `decrypt`: decrypt the private key so you can use it
- `send {recipient} {amount}`: send {amount} tokens to {recipient}
- `sendL2 {recipient} {amount}`: send {amount} tokens to {recipient} via the layer 2 rollup system (alpha)
- `balance {key}`: get the balance associated with the public key {key} (tip: running `balance` without passing {key} will get your own balance)
- `savestate`: save a backup of the current state of the blockchain to a file
- `loadstate`: load a backup of the current state of the blockchain from a file
- `addpeer {ip}`: connect to a peer
- `startAnalysisConsole`: start a console for analyzing the status and history of the blockchain and network
- `bootstrap`: connect to your peers' peers for increased speed, reliability, and decentralization
- `exit`: exit the console

To get started, run `keygen` to generate a new key. To get your balance, find your public key in the `key.json` file (the long number following `"Y":`), and run `balance {YOUR KEY HERE}`. To send currency, type `send {RECIPIENT PUBLIC KEY} {AMOUNT}`. You'll have to ask the recipient for their public key. When you're done, type `encrypt` to encrypt your private key and store it safely. You can decrypt it later to use it again with `decrypt`. You must use a passcode that is a multiple of 16 characters long for encryption and decryption.

>[!CAUTION]
>Write your passcode down somewhere safe. If you lose it, there is no "forgot my password" option like in many Web 2.0 applications-- your funds will be irrecoverable.

### To run a node:

To run the node software, which keeps the blockchain distributed across the p2p network, run:

```bash
./builds/node/node -serve -port [PORT]
```

### To run a miner:

To run the mining software, which adds new blocks to the blockchain in exchange for a reward, run:

```bash
./builds/node/node -serve -mine -port [PORT]
```

### To connect to a peer:

To connect to a peer, enter the BlockCMD console and run:

```
addPeer http://PEER_IP:PORT
```

## License

This software is released under the GNU General Public License v3.0.

## Contributing

Contributions are greatly appreciated! See the CONTRIBUTING file for more details.
