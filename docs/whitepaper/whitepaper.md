**Polycash: An Increased-Decentralization P2P Blockchain System**

**Asher Wrobel**

[https://github.com/Ashy5000](https://github.com/Ashy5000)

**Abstract**

Most widely-used cryptocurrency blockchains utilize the PoW (Proof of Work) or PoS (Proof of Stake) consensus algorithms. Unfortunately, these algorithms face a variety of issues. PoW decreases decentralization with mining pools and uses large amounts of energy, while PoS systems face lower resilience against so-called 51% attacks. In response to these issues, APoW(Adjusted Proof of Work) is introduced, using difficulty adjustment algorithms that increase decentralization and decrease power consumption in the blockchain. Combining this consensus algorithm with a new system of verification, time verifiers, creates Polycash, a blockchain designed to allow secure transactions and the operation of peer-to-peer applications to be conducted with increased decentralization and security as opposed to traditional solutions.

**1 Introduction**

In the current day, financial institutions have largely become adopted as a mediary for digital transactions. Using financial institutions relies on trust, an issue attempted to be resolved by blockchain-based cryptocurrencies. However, blockchains are becoming increasingly centralized when using the existing PoW of PoS consensus algorithms. In PoW, large mining pools have large influence over the state of the blockchain, and cooperation between multiple pools could lead to a 51% attack on the network. In PoS, individuals with large amounts of stake are able to have influence over the blockchain and the network in a similar manner. There are also other issues. PoW uses large amounts of energy, and PoS has low resistance to a 51% attack, as only a third of nodes must be malicious to take over the network and blockchain.

&nbsp;&nbsp;&nbsp;&nbsp;To resolve these issues, a new consensus algorithm is needed: the APoW (Adjusted Proof of Work) algorithm. Instead of giving blockchain power directly proportional to mining power, APoW calculates block speed on a per-miner basis. It decreases target time (base 1 minute) by up to 50% for a miner with infinite mining rate, and increases target time by up to 50% for a miner with zero mining rate. This leads to an increase of decentralization in the blockchain and the network, and allows for a lower energy requirement than with a standard PoW blockchain.

&nbsp;&nbsp;&nbsp;&nbsp;To adapt to the updated algorithm, multiple protocol changes must be put into place to ensure the security and reliability of the APoW blockchain.

**2 Adjusted Proof of Work**

To increase decentralization in the blockchain, a new consensus algorithm, Adjusted Proof of Work (APoW), is employed. The algorithm’s aim is to increase the probability of a miner with a lower hashrate having the ability to create a block, and decrease the probability of a miner with a higher hashrate from doing the same. While the system is designed to maintain the higher chance of a higher hashrate miner winning a block, the gap between higher hashrate miners and lower hashrate miners narrows when using APoW.

&nbsp;&nbsp;&nbsp;&nbsp;The mechanism utilized by APoW to increase decentralization in this manner is block difficulty. As opposed to the traditional PoW algorithm, APoW calculates difficulty on a per-miner basis, attempting to adjust the difficulty for that miner every block so the time it takes them to mine a block matches a target time. In order to encourage miners to contribute their computing power to the network while still giving lower hashrate miners a chance of successfully mining a block, a modified version of the sigmoid function is utilized to alter the target time for a given miner based on the difficulty of the previous block they mined:

$$\frac{1}{1+e^\frac{-(x-mdpm)}{mdpm}}$$

1 DPM (Difficulty Point per Minute) represents the speed of a miner that can mine a block with difficulty 1 in 1 minute.

_x_ represents the difficulty of the previous block the miner has created divided by the time it took to mine it, measured in DPM.

_mdpm_ is a constant representing 1 million DPM

&nbsp;&nbsp;&nbsp;&nbsp;Dividing the difficulty by the result of the above formula allows faster miners an advantage, but not to the degree of standard PoW systems.

&nbsp;&nbsp;&nbsp;&nbsp;To determine if a cryptographic proof of work is valid, the SHA-3-512 hash of the block, represented in this case as an unsigned 64-bit integer, must be less than the maximum value of an unsigned 64-bit integer (18446744073709551615) divided by the block’s difficulty. This means if block _A_ has difficulty _D<sub>a</sub>_ and block B has difficulty _D<sub>b</sub>_, and if _D<sub>b</sub> _/ _D<sub>a</sub>_ = _n_, then block B is expected to take _n_ times longer to mine than block _A_. In other words, difficulty is proportional to mining time.

**3 Time Verifiers**

When difficulty is adjusted based on the time it takes for a miner to create blocks, it becomes necessary to ensure miners are honest about the time they start and finish mining. If these two timestamps can be verified, it would also prevent all manipulation of past blocks in the blockchain, even if 51% of the mining power in the network is controlled. To implement this security feature, time verifiers are used in the network. Any node that has mined a block may become a time verifier. When a node starts or finishes mining a block, they must request and collect the signatures of time verifiers. Time verifiers will provide their signature if the current time is within a 10-second range of the time they need to verify. There may be no less than 75% of the time verifiers in the previous block in the current block, and each additional signature beyond the number of signatures in the previous block earns the miner a reward. Because of the reward gained from adding additional time verifiers, miners will attempt to gather signatures from as many nodes as possible. Therefore, the number of signatures will roughly match the number of verifiers in the network. As a given miner will never have control of over 50% of the verifiers in the network, they will fall short of the required number of signatures if they attempt to lie about the times they begin and end mining.

&nbsp;&nbsp;&nbsp;&nbsp;In addition to the number of signatures a miner needs to create a valid block, the majority of the verifiers the miner gathers signatures from must also be verifiers in the previous block. This almost completely prevents a miner from gaining malicious verification. Even if a miner controls 51% percent of time verifiers, these verifiers must not sign blocks- if they did, they would contribute to the number of verifiers in a block, and the number of signatures the miner could receive in their malicious block would still drop by about 50%, below the required 25%. This means a miner must control over 75% of miners in the network _and_ 50% of the computing power. Without the second rule, miners controlling malicious verifiers could have those verifiers not sign blocks, so they would not contribute to the verifier count. Then, they could be used to verify the malicious block. However, none or very few of those verifiers will have verified the previous blocks, so with the second rule the malicious block will be invalid. This security measure, along with the next, creates an extremely secure blockchain with significant advantages over one using standard PoW.

**4 Miner Count Limits**

In order to decrease the energy usage of the network and to slow the growth in the number of time verifiers, the number of miners in the network is limited to a maximum. If a miner with a public key not present as the miner in any previous block creates and broadcasts a new block and the number of unique public keys present as miners in the blockchain equals or exceeds the maximum, the block is rejected. After a miner has created one block, it can continue to create blocks even if the maximum number of miners have been reached- they are already registered. After they are registered, they do not have a strong motivation to use more energy to mine faster past a certain point, as the sigmoid difficulty function does not significantly reward faster miners after a point. This motivates miners to use a low to medium mining rate to optimize their profits- past a certain mining rate, energy costs outweigh mining reward benefits. This decreases the energy usage of the network.

&nbsp;&nbsp;&nbsp;&nbsp;To calculate the maximum number of miners at the current time, divide the number of blocks in the blockchain by 20, rounding up.

**5 Smart Contracts**

To enable cryptographically verified functionality on the blockchain, smart contracts are implemented. with a blockchain-specific assembly language. When a miner begins mining a block, they evaluate all smart contracts that would be triggered in that block, and, if they initiate any transactions, those transactions are placed in the block to be mined. The block is considered invalid if the smart contract transactions are missing or do not match the correct ones.

&nbsp;&nbsp;&nbsp;&nbsp;A smart contract may automatically make transactions on the behalf of another party, provided that party signs the smart contract when it is created, granting their approval of the instructions contained within it. This functionality is required to allow a contract to be deterministically executed without the need for trust of a party to fulfill their part of the agreement.

&nbsp;&nbsp;&nbsp;&nbsp;Deploying a smart contract requires a deployment fee, paid to the miner by the node deploying the smart contract.

**5.1 Virtual Machine**

Smart contracts running on the blockchain use a virtual machine with a custom instruction set to deterministically initiate transactions. It uses a Polycash-specific instruction set and assembly language to evaluate these contracts.

**5.1.1 Memory Management**

The memory of the Polycash virtual machine is contained within a set of virtual “buffers”. These are not strictly buffers in the traditional sense. In Rust, the language the VM is written in, they are represented as a value of type Vec&lt;u8>. This means they are expandable lists of bytes with no fixed length. Different data types may be represented by different lengths and organizations of buffers. There are an unlimited number of buffers, and each is an unlimited size, so the virtual machine may hold as much data as it needs in memory (up to its memory limit).

The buffer at hex address 0x00000000, used for global errors is pre-initialized during the VM boot process.

**5.1.2 Errors**

Errors are split into two types: local errors and global errors. Local errors write an error value (hex code 0x1) into a pre-initialized error buffer if an error occurs in an instruction. Global errors write an error value (also hex code 0x1) into the pre-initialized error buffer at 0x00000000. If the error buffer needed to throw a local error is not found, a global error is thrown.

See [Polycash VM Spec](https://github.com/ashy5000/cryptocurrency/blob/master/docs/whitepaper/instructions.md) for extended subsection 5.1.3 on the VM instruction set.

**6 Transaction Body**

To allow data to be sent through the blockchain, a body may be attached to each transaction containing data of any form. A data fee is paid to the miner of that block of a size proportional to the amount of data transferred, in bytes.

&nbsp;&nbsp;&nbsp;&nbsp;This transaction body feature allows for the system to be used to create new tokens, implement two-factor authentication, or to permanently and reliably store data, while utilizing the security and reliability of the blockchain.

**7 Conclusion**

We have seen that using an alternative consensus algorithm, APoW, with the help of the time verification protocol increases decentralization and decreases energy usage when compared to traditional PoW blockchains. It also avoids the security issues present in PoS or PoH blockchains, where controlling 1/3 of stake in the blockchain can cause the network to fail. We have outlined a new blockchain, Polycash, that implements these features alongside a smart contract system in order to enable secure digital transactions and decentralized applications.
