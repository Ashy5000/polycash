# Setup
You have three options to participate in the network:
1. Client: Send and receive currency
2. Node: Propagate transactions and support the network
3. Miner: Use your computing power to maintain the network's security and gain
   rewards

### Client Setup
This setup assumes a Linux/Unix-like command line.

First, clone the repository.

```bash
git clone https://github.com/Ashy5000/cryptocurrency
```

Find the correct executable.

|         | x86_64                         | arm64                          | arm                          | 386                          |
|---------|--------------------------------|--------------------------------|------------------------------|------------------------------|
| Linux   | builds/node/node_linux-amd64   | builds/node/node_linux-arm64   | builds/node/node_linux-arm   | builds/node/node_linux-386   |
| MacOS   | builds/node/node_darwin-amd64  | builds/node/node_darwin-arm64  | Combination not possible     | Combination not possible     |
| Windows | builds/node/node_windows-amd64 | builds/node/node_windows-arm64 | builds/node/node_windows-arm | builds/node/node_windows-386 |

Run the correct executable.

```bash
./builds/node/{executable}
```

```
You will now be in the BlockCMD console. It should look something like this:
Copyright (C) 2024 Asher Wrobel
This program comes with ABSOLUTELY NO WARRANTY. This is free software, and you are welcome to redistribute it under certain conditions.
To see the license, type `license`.
BlockCMD console (encrypted: true):
```

Type `bootstrap` and then enter. This will connect your device with other peers.

Type `keygen` and then enter. This will generate a new Dilithium2 keypair. If you wish, type `encrypt` and then enter. This will encrypt your keypair with a 16-character-long passcode. When you need to check your balance or send currency, type `decrypt` to make it usable. When you are done, use `encrypt` once more.

To use a graphical application to manage funds, first install the Rust programming language if you haven't already. Then, move into the `gui_wallet` directory and use the `cargo run` command to build and run the application.

If you prefer a text-based interface, use the `help` command in the BlockCMD console for more information.
