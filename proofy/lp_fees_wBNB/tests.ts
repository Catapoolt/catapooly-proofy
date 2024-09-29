import {Field, ProofRequest, Prover, ReceiptData, TransactionData} from "brevis-sdk-typescript";

async function test_amount_1_wBNB() {
  const prover = new Prover('localhost:32001');

  const proofReq = new ProofRequest();

  proofReq.addTransaction(
    new TransactionData({
      hash: "0xdd94da45b82440d87911c0f1a7972fdae1f7b259dba5ac742d7a564c746bf511",
      chain_id: 97,
      block_num: 44116540,
      nonce: 9,
      gas_tip_cap_or_gas_price: "9700000000",
      gas_fee_cap: "0",
      gas_limit: 206531,
      from: "0xc8fb199f3d4f3ebcd112023ceb9536ca12a4d198",
      to: "0x427bf5b37357632377ecbec9de3626c71a5396c1",
      value: "0",
    }),
  );

  proofReq.addReceipt(
    new ReceiptData({
      block_num: 44116540,
      tx_hash:
        "0xdd94da45b82440d87911c0f1a7972fdae1f7b259dba5ac742d7a564c746bf511",
      fields: [
        new Field({
          contract: "0x427bF5b37357632377eCbEC9de3626C71A5396c1",
          log_index: 3,
          event_id:
            "0x40d0efd1a53d60ecbf40971b9daf7dc90178c3aadc7aab1765632738fa8b8f01",
          is_topic: false,
          field_index: 2,
          value: "1807582057",
        }),
      ],
    })
  );

  const proofRes = await prover.prove(proofReq);

  console.log(JSON.stringify(proofRes))
}

test_amount_1_wBNB().then(() => console.log("Finished test test_amount_1_wBNB!"));
