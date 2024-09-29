import {GoldRushClient} from "@covalenthq/client-sdk";
import {
  Transaction,
  TransactionResponse
} from "@covalenthq/client-sdk/dist/esm/src/utils/types/TransactionService.types";
import {Field, ProofRequest, ReceiptData, TransactionData} from "brevis-sdk-typescript";

const collectTopicHash = "0x40d0efd1a53d60ecbf40971b9daf7dc90178c3aadc7aab1765632738fa8b8f01"

async function get_collect_transactions_for_wallet(walletId: string, startingBlock: number): Promise<Transaction[]> {
  const client = new GoldRushClient("cqt_rQ4B9bxHBDmvKCgcDG3fXBhDFpWm");

  const logEvents = client.BaseService.getLogEventsByTopicHash("bsc-testnet", collectTopicHash, {
    startingBlock: startingBlock,
    endingBlock: "latest"
  })

  //TODO: uncomment this after testing is done
  // const tx_hashes = []
  // for await (const event of logEvents) {
  //   for (let item of event.data.items) {
  //     tx_hashes.push(item.tx_hash)
  //   }
  // }

  const tx_hashes = [
    "0x1239fcb2c8c366707004e74ab103d1c6bcce19fccd6cb0633850b9b3a2d1e237",
    "0xdd94da45b82440d87911c0f1a7972fdae1f7b259dba5ac742d7a564c746bf511"
  ]
  console.log(JSON.stringify(tx_hashes))

  const collect_txs_for_wallet = []
  for (let tx_hash of tx_hashes) {
    const response = await client.TransactionService.getTransaction("bsc-testnet", tx_hash, {
      noLogs: false
    })

    const data: TransactionResponse | null = response.data;
    if (data == null) {
      continue
    }

    for (let tx of data.items) {
      if (tx.from_address.toLowerCase() === walletId.toLowerCase()) {
        collect_txs_for_wallet.push(tx)
      }
    }
  }
  // TODO: and now filter all transactions for which amount1 collect log is wBNB contract

  return collect_txs_for_wallet
}

export async function proof_request(): Promise<ProofRequest> {
  const proofReq = new ProofRequest();

  const txs = await get_collect_transactions_for_wallet("0xc8fB199F3d4F3ebCD112023CEB9536cA12A4D198", 44116540);

  for (let tx_index = 0; tx_index < txs.length; tx_index++) {
    const tx = txs[tx_index];
    let transactionData = new TransactionData({
      hash: tx.tx_hash,
      chain_id: 97,
      block_num: tx.block_height,
      nonce: 0,  // TODO: no idea with goldrush.dev
      gas_tip_cap_or_gas_price: tx.gas_price.toString(),
      gas_fee_cap: "0",
      gas_limit: tx.gas_offered,
      from: tx.from_address,
      to: tx.to_address,
      value: tx.value.toString(),
    });
    console.log("TransactionData: ", JSON.stringify(transactionData.toObject()));
    proofReq.addTransaction(
      transactionData,
      tx_index
    );

    const smallest_log_offset = Math.min(...tx.log_events.map((log_event) => log_event.log_offset))
    const log_events = tx.log_events.filter((log) => log.raw_log_topics.includes(collectTopicHash))

    if (log_events.length != 1) {
      throw new Error("Each transaction should have only 1 collect log event!")
    }

    const log = log_events[0];
    for (let i = 0; i < log.decoded.params.length; i++) {
      const param = log.decoded.params[i];
      if (param.name == "amount1") {
        let receiptData = new ReceiptData({
          block_num: log.block_height,
          tx_hash: log.tx_hash,
          fields: [
            new Field({
              contract: log.sender_address,
              log_index: log.log_offset % smallest_log_offset,
              event_id: "0x40d0efd1a53d60ecbf40971b9daf7dc90178c3aadc7aab1765632738fa8b8f01",
              is_topic: false,
              field_index: i,
              value: param.value,
            }),
          ],
        });
        console.log("ReceiptData: ", JSON.stringify(receiptData.toObject()));
        proofReq.addReceipt(
          receiptData,
          tx_index
        );
      }
    }
  }

  let transactions_number = proofReq.getTransactions().length;
  let receipts_number = proofReq.getReceipts().length;
  console.log(`Found ${transactions_number} transactions and ${receipts_number} receipts!`);
  if (transactions_number != receipts_number) {
    throw new Error("The number of transaction and receipts should be the same!")
  }

  return proofReq
}
