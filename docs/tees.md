# Teechan vs Starkware

TEE Approach: Encrypt / Encode data and enable only a few people to decrypt the code and verify its execution.

Problems:
- Power to decrypt is centralized between a few people.
- Some parts of execution not verifiable by everyone (some parts can only be guaranteed by the TEE)
- Side channel attacks on TEEs

ZKStarks Approach: Prover's work is done off chain and verification on chain.

Advantages:
- Inputs are private, private information can be done off chain
- All parts of the proof are verifiable by everyone
- More scalable than public blockchains since computation is done off chain.

What needs to be private in a STARK based environment can be defined by the person who defines the contract

> One component would be a noisy signal which all recover which when aggregated reveal important macro economic or state

I guess this need not be noisy - you could shard the message into multiple parts and spread it only to those entities you want to.

> The other component would be idiosyncratic information with parties might like to keep secret

Data which needs to be kept private can be kept private by the person who defines the contract.

> Or in another case individualized messages are scrambled and kept private by in the execution of the contract, allocations or transfers might need to be revealed.

Execution can be entirely private, so I don't think this should be a problem.
