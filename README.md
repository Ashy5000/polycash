# ashy5000/cryptocurrency
The home of a power-efficient, scalable, and modern cryptocurrency blockchain designed to resolve the issues presented by the traditional PoW (Proof of Work) incentive mechanism while maintaining decentralization. Written in Golang and Rust for secure algorithm implementations, lightweight networking, memory safety, speed, and reliability.

## Modified PoW protocol
As opposed to the original Proof of Work (PoW) protocol, this blockchain does not calculate block difficulty on a global level. Instead, it is adjusted at a per-miner basis, removing the reward for using more computing power and thus more energy. If a miner has a higher hash rate and thus initially mines blocks faster, they will be provided with higher difficulty blocks and will end up mining blocks at a **constant rate** of *1 min/transaction*. In addition, the block difficulty is retargeted *every* time a miner mines a block, ensuring quick reactions to changes in miner hash rates. To prevent miners from registering multiple private keys in an attempt to avoid the difficulty constraints, there is planned to be a high initial difficulty, or the constant difficulty for the first block a miner mines.

## Directory structure
The root directory of this project is occupied by the Golang source code that nodes run in order to interact with each other and the decentralized blockchain. In the ```peer_server``` directory, there is Rust code that can be run by servers to maintain a list of peers in the network. Nodes can connect to these servers or maintain their own lists. You will probably only need to run the Golang code.