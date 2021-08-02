package executor

import (
	"fmt"

	solcmd "github.com/Gravity-Tech/solanoid/commands"
	"github.com/Gravity-Tech/solanoid/commands/contract"
	solexe "github.com/Gravity-Tech/solanoid/commands/executor"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
	"github.com/valyala/fasthttp"
)


type Allocator struct {
	account solcmd.OperatingAddress
}

func NewAllocator(pkpath string) (*Allocator, error) {
	account, err := solcmd.NewOperatingAddress(nil, pkpath, &solcmd.OperatingAddressBuilderOptions{ Overwrite: false })
	if err != nil {
		return nil, err
	}
	return &Allocator{
		*account,
	}, nil
}

func (alloc *Allocator) Allocate(tokenMint common.PublicKey) (*WalletAllocResponse, error) {
	RPCEndpoint, _ := solcmd.InferSystemDefinedRPC()

	deployerExecutor, err := solexe.NewEmptyExecutor(alloc.account.PrivateKey, RPCEndpoint)
	if err != nil {
		return nil, err
	}

	targetWallet, ix := contract.CreateAssociatedTokenAccountIXNonFailing(alloc.account.PublicKey, tokenMint)

	deployerExecutor.SetAdditionalSigners([]solexe.GravityBftSigner { 
		*solexe.NewGravityBftSignerFromAccount(targetWallet),
	})

	response, err := deployerExecutor.InvokeIXList(
		[]types.Instruction { *ix },
	)
	if err != nil {
		return nil, err
	}

	return &WalletAllocResponse{
		PublicKey: targetWallet.PublicKey.ToBase58(),
		TokenMint: tokenMint.ToBase58(),
		TxSignature: response.TxSignature,
	}, nil
}


// request handler in fasthttp style, i.e. just plain function.
func (alloc *Allocator) RequestHandle (ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hi there! RequestURI is %q", ctx.RequestURI())
}