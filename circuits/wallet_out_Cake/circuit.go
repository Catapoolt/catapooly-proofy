package wallet_out_Cake

import (
	"fmt"
	"github.com/brevis-network/brevis-sdk/sdk"
	"github.com/consensys/gnark/frontend"
	"github.com/ethereum/go-ethereum/common"
)

type CakeFromWallet struct {
	walletId sdk.Uint248
	amount   sdk.Uint248
}

func (c CakeFromWallet) Values() []frontend.Variable {
	var ret []frontend.Variable
	ret = append(ret, c.walletId.Values()...)
	ret = append(ret, c.amount.Values()...)
	return ret
}

func (c CakeFromWallet) FromValues(vs ...frontend.Variable) sdk.CircuitVariable {
	nf := CakeFromWallet{}

	start, end := uint32(0), c.walletId.NumVars()
	nf.walletId = c.walletId.FromValues(vs[start:end]...).(sdk.Uint248)

	start, end = end, end+c.amount.NumVars()
	nf.amount = c.amount.FromValues(vs[start:end]...).(sdk.Uint248)

	return nf
}

func (c CakeFromWallet) NumVars() uint32 {
	return c.walletId.NumVars() + c.amount.NumVars()
}

func (c CakeFromWallet) String() string {
	return ""
}

var _ sdk.CircuitVariable = CakeFromWallet{}

type WalletOutCakeCircuit struct{}

var _ sdk.AppCircuit = &WalletOutCakeCircuit{}

var cake = sdk.ConstUint248(common.HexToAddress("0xD3677F083B127a93c825d015FcA7DD0e45684AcA"))
var zero = sdk.ConstUint248(0)

func (c *WalletOutCakeCircuit) Allocate() (maxReceipts, maxStorage, maxTransactions int) {
	// Receipts and Transactions need to be equal
	return 10, 0, 0
}

func (c *WalletOutCakeCircuit) Define(api *sdk.CircuitAPI, in sdk.DataInput) error {
	receipts := sdk.NewDataStream(api, in.Receipts)

	fmt.Print("Receipts: ")
	fmt.Println(sdk.Count(receipts))

	sdk.AssertEach(receipts, func(cur sdk.Receipt) sdk.Uint248 {
		return api.Uint248.And(
			api.Uint248.IsEqual(cur.Fields[0].Contract, cake),
			api.Uint248.IsEqual(cur.Fields[1].Contract, cake),
		)
	})

	// todo: we assume each receipt fields[0] is wallet id and fields[1] is the amount
	walletAddress := api.ToUint248(sdk.GetUnderlying(receipts, 0).Fields[0].Value)
	amount := sdk.Reduce(receipts, zero, func(accumulator sdk.Uint248, current sdk.Receipt) (newAccumulator sdk.Uint248) {
		fmt.Print("Wallet address: ")
		fmt.Println(current.Fields[0].Value.String())
		fmt.Print("Wallet amount: ")
		fmt.Println(current.Fields[1].Value.String())
		return api.Uint248.Add(accumulator, api.ToUint248(current.Fields[1].Value))
	})

	fmt.Print("Final wallet address: ")
	fmt.Println(walletAddress.String())
	fmt.Print("Final wallet amount: ")
	fmt.Println(amount.String())

	api.OutputAddress(walletAddress)
	api.OutputUint(248, amount)

	return nil
}
