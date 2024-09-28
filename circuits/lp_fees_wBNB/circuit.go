package lp_fees_wBNB

import (
	"fmt"
	"github.com/brevis-network/brevis-sdk/sdk"
	"github.com/consensys/gnark/frontend"
	"github.com/ethereum/go-ethereum/common"
)

type LPFeesWBNBCircuit struct{}

var _ sdk.AppCircuit = &LPFeesWBNBCircuit{}

func (c *LPFeesWBNBCircuit) Allocate() (maxReceipts, maxStorage, maxTransactions int) {
	// Receipts and Transactions need to be equal
	return 10, 0, 10
}

type CollectedFee struct {
	walletId sdk.Uint248
	amount   sdk.Uint248
}

func (c CollectedFee) Values() []frontend.Variable {
	var ret []frontend.Variable
	ret = append(ret, c.walletId.Values()...)
	ret = append(ret, c.amount.Values()...)
	return ret
}

func (c CollectedFee) FromValues(vs ...frontend.Variable) sdk.CircuitVariable {
	nf := CollectedFee{}

	start, end := uint32(0), c.walletId.NumVars()
	nf.walletId = c.walletId.FromValues(vs[start:end]...).(sdk.Uint248)

	start, end = end, end+c.amount.NumVars()
	nf.amount = c.amount.FromValues(vs[start:end]...).(sdk.Uint248)

	return nf
}

func (c CollectedFee) NumVars() uint32 {
	return c.walletId.NumVars() + c.amount.NumVars()
}

func (c CollectedFee) String() string {
	return ""
}

var _ sdk.CircuitVariable = CollectedFee{}

var zero = sdk.ConstUint248(0)
var wBNB = sdk.ConstUint248(common.HexToAddress("0xae13d989daC2f0dEbFf460aC112a837C89BAa7cd"))

func (c *LPFeesWBNBCircuit) Define(api *sdk.CircuitAPI, in sdk.DataInput) error {
	txs := sdk.NewDataStream(api, in.Transactions)
	rxs := sdk.NewDataStream(api, in.Receipts)

	fmt.Print("Transactions: ")
	fmt.Println(sdk.Count(txs))
	fmt.Print("Receipts: ")
	fmt.Println(sdk.Count(rxs))

	walletFees := sdk.ZipMap2(txs, in.Receipts.Raw, func(a sdk.Transaction, b sdk.Receipt) CollectedFee {
		// assert event id: "0x40d0efd1a53d60ecbf40971b9daf7dc90178c3aadc7aab1765632738fa8b8f01"
		return CollectedFee{
			walletId: a.From,
			amount:   api.ToUint248(b.Fields[0].Value),
		}
	})

	walletAddress := sdk.GetUnderlying(walletFees, 0).walletId
	finalAmount := zero
	sdk.Map(walletFees, func(fee CollectedFee) CollectedFee {
		api.Uint248.IsEqual(walletAddress, fee.walletId)

		finalAmount = api.Uint248.Add(finalAmount, fee.amount)
		fmt.Print("Wallet address: ")
		fmt.Println(fee.walletId.String())
		fmt.Print("Wallet amount: ")
		fmt.Println(fee.amount.String())
		return fee
	})

	//by, err := sdk.GroupBy(
	//	walletFees,
	//	func(accumulator CollectedFee, current CollectedFee) CollectedFee {
	//		accumulator.walletId = current.walletId
	//		accumulator.amount = api.Uint248.Add(accumulator.amount, current.amount)
	//		return accumulator
	//	},
	//	CollectedFee{walletId: sdk.ConstUint248(0), amount: sdk.ConstUint248(0)},
	//	func(current CollectedFee) sdk.Uint248 {
	//		return current.walletId
	//	},
	//)
	//if err != nil {
	//	return err
	//}

	//tx := sdk.GetUnderlying(txs, 0)
	//// This is our main check logic
	//api.Uint248.AssertIsEqual(tx.Nonce, sdk.ConstUint248(0))
	//
	//// Output variables can be later accessed in our app contract
	fmt.Print("Final wallet address: ")
	fmt.Println(walletAddress.String())
	fmt.Print("Final wallet amount: ")
	fmt.Println(finalAmount.String())
	api.OutputAddress(walletAddress)
	api.OutputAddress(wBNB)
	api.OutputUint(248, finalAmount)

	return nil
}
