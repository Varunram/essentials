# Shadow LN

So the idea is to have some sort of a private LN between hexa users so hexa-hexa transactions don't have any fee associated. The main things to tackle here are:

1. Having proof that x btc deposited in your account correlates with a specific amount in the hexa channel
2. Ensuring that others can't create "hexa btc" (hexa-hexa txs' currency)

This can be a very powerful primitive - coinbase does user balance accounting  with the help of a giant SQL database but we can do the same by having hexa channels between users. Say Alice has 1 btc and wants to send 0.5btc to Bob and Charlie. Alice first uses the shadow lightning network to create a channel with the Bithyve server and then she pushes 0.5 btc to Bob and Charlie. The transaction is onion encrypted so the origin/destination pair isn't known both together by the same entity.

Now Bob and Charlie have two options:
1. Spend the money that Alice just sent them
  a. To another hexa account - goto 1
  b. To external btc account - The channel between B/C and Bithyve server closes and this takes the last known state of the transaction (ie 0.5btc to B/C) on chain. Once the tx is confirmed on chain, B/C can spend their funds. Note that there are two transactions to be made instead of one here and that this would reuslt in an increased wait time.
2. Hold the coins - If Alice and Bob hold the coins, there is no requirement for an immediate tx to send the funds from Alice to B/C on chain. Instead this gives us a variety of interesting alternatives:
  a. Send funds on chain when the average tx fee is below x sat - We can settle the Alice - B/C tx on chain after a specific amount of time by closing the channel with the Bithyve server.
  b. Send funds on chain only when the user needs to spend them
  c. Send funds on chain every y blocks - can be done with the help of timelock magic.

The shadow LN would be similar to mainnet LN with a few changes:
1. No fee since we want hexa-hexa txs free of cost
2. No network routing problems since there is a hub-spoke model
3. Less reliance on splicing in/out since hexa can attest to a user's balance

Hexa could use LN as well but we'd have to worry about routing, rebalancing and other stuff which is not needed if we are to have a shadow LN. One main constraint of the shadow LN would be that all nodes/validators are connected to the bigger hub. But this hub can be both on Shadow and Mainnet LN and could switch between the two by splicing in/out when needed. The hub could also facilitate submarine swaps and stuff which may prove to be useful in the future. Mining pool payouts or company monthly payouts into shadow LN are also possible to save on tx fees.
