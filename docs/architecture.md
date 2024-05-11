# Architecture

The architecture of the blockchain revolves around the idea of time verififcation, Proof of Work as a spacing mechanism, and difficulty adjustment.

### Time Verification

To ensure miners are honest about the time it took them to mine a block, each block must be signed by at least 3/4 of the number of nodes that signed the previous block. For each additional time verifier past the number of verifiers in the previous block, the miner for that block earns an additional 0.1 tokens in addition to the 1.0 tokens from their mining reward. This ensures miners always attempt to get the signatures for as many nodes as possible. These nodes are called the 'time verifiers'. This proccess occurs twice for every block. First, the time verifiers ensure that the block began to be mined within a 10 second range of the current time (when their verification was requested). This ensures nodes joining the network mid-block will have their difficulties correctly calcualted. Then they ensure that the current time is less than 10 seconds after the time the miner claimed to have finished mining. If the miner says it took longer than it actually did to mine a block, the time verifiers will see that the finishing time is in the future and refuse to sign the block. If the miner attempts to wait before sending it to a time verifier, they will most likely lose the block before it is sent. This system creates a history of blocks with verified timestamps, ensuring the honesty of miners and their mining times, preventing miners from manipulating block difficulty, and preventing forks.

### Maximum block size

To enable high speeds and throughput of the network, there is no maximum block size.

### Proof of Work as a spacer

In this blockchain, PoW is not used as simply a way to prove the validity of a block or fork. It is also used as a spacer. Time verifiers will accept blocks with a 10-second room for error. If these errors overlap, vulnerabilities will occur. Attackers could theoretically have time verifiers accept invalid blocks. To prevent this, the network uses PoW to maintain a 30 sec-1 min, 30 sec gap between each block. Combining Proof of Work's security with the security of time verification leads to a two-layer security system with high resistance to attacks.

### Difficulty Adjustment

One major change is made to the existing method to increase the efficiency and lower the power usage of the network. Instead of a block's difficulty being calculated based on the speed of the network, it is calculated based on the speed **of the miner that is mining it**. In this way, a miner with more hashing power will gain only a small benefit over a miner with less hashing power. If a block is found more quickly, the next one will be more difficult to mine, and the speed will balance out (for the most part). However, the target time, on average 1 minute, does vary based on your mining rate- a higher mining rate means a lower target time. The lowest target time is 30 seconds and the highest 1 minute 30 seconds, so there's no substantial benefit from mining faster, but there's still a little. It's important to encourage miners to contribute more computing power to the network, even if you don't provide a large benefit. This method will strengthen the network while also increasing decentralization and efficiency. The difficulty calculations are made each block for every miner on the network, so attempting to "trick" the difficulty calculations is not possible. However, without further modification, this approach is impossible to enforce. As miners may create an arbitrary amount of public/private key pairs, they may bypass the difficulty adjustment feature and create a new key pair for every block they mine. To prevent unlimited key pairs from being generated, this blockchain sets a maximum miner limit that increases with the number of blocks in the blockchain. This ensures that unlimited miners cannot be created. Although this feature was initially created to prevent hashrate forking, there are now additional measures to discourage that behavior. Now, it stands as a way to prevent unlimited time verifiers from being created, since all miners may be time verifiers.

### Smart Contracts

To execute smart contracts, the Polycash blockchain uses its own virtual machine. See VM docs [here](vm.md).
