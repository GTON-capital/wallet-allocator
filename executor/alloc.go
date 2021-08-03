package executor

import (
	"encoding/json"
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
	account, err := solcmd.ReadOperatingAddress(nil, pkpath)
	if err != nil {
		return nil, err
	}
	return &Allocator{
		*account,
	}, nil
}

func GetTokenAccountsByOwner(tokenMint, owner common.PublicKey) (*[]string, error) {
	RPCEndpoint, _ := solcmd.InferSystemDefinedRPC()

	encodedBody := fmt.Sprintf(`
	{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "getTokenAccountsByOwner",
		"params": [
		  "%v",
		  {
			"mint": "%v"
		  },
		  {
			"encoding": "jsonParsed"
		  }
		]
	}
	`, owner.ToBase58(), tokenMint.ToBase58())

	req := fasthttp.AcquireRequest()
	req.SetBody([]byte(encodedBody))

	req.Header.SetMethodBytes([]byte("POST"))
	req.Header.Set("Content-Type", "application/json")

	req.SetRequestURIBytes([]byte(RPCEndpoint))
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		panic("handle error")
	}
	fasthttp.ReleaseRequest(req)

	body := res.Body()

	var resp TokenAccountsByOwnerResponse

	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	fasthttp.ReleaseResponse(res)

	var tokenAccounts []string

	for _, tokenAccount := range resp.Result.Value {
		tokenAccounts = append(tokenAccounts, tokenAccount.Pubkey)
	}

	return &tokenAccounts, nil
}

func (alloc *Allocator) Allocate(tokenMint, owner common.PublicKey) (*WalletAllocResponse, error) {
	RPCEndpoint, _ := solcmd.InferSystemDefinedRPC()

	tokenAccounts, err := GetTokenAccountsByOwner(tokenMint, owner)

	if err != nil {
		return nil, err
	}

	if len(*tokenAccounts) > 0 {
		return &WalletAllocResponse{
			PublicKey: (*tokenAccounts)[0],
			TokenMint: tokenMint.ToBase58(),
			TxSignature: "",
		}, nil
	}

	deployerExecutor, err := solexe.NewEmptyExecutor(alloc.account.PrivateKey, RPCEndpoint)
	if err != nil {
		return nil, err
	}

	ix, walletAddress, err := contract.CreateAssociatedTokenAccountIX(alloc.account.PublicKey, owner, tokenMint)
	if err != nil {
		return nil, err
	}

	response, err := deployerExecutor.InvokeIXList(
		[]types.Instruction { *ix },
	)
	if err != nil {
		return nil, err
	}

	return &WalletAllocResponse{
		PublicKey: walletAddress.ToBase58(),
		TokenMint: tokenMint.ToBase58(),
		TxSignature: response.TxSignature,
	}, nil
}

func handleError(ctx *fasthttp.RequestCtx, err error) {
	r := &Response{ Status: "error", Result: err.Error() }
	fmt.Fprint(ctx, r.Serialize())
	ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
}

// request handler in fasthttp style, i.e. just plain function.
func (alloc *Allocator) RequestHandle (ctx *fasthttp.RequestCtx) {

	switch string(ctx.RequestURI()) {
	case "/api/associated-token-account/alloc":
		if string(ctx.Method()) != "POST" {
			break
		}

		ctx.Response.Header.SetCanonical([]byte("Content-Type"), []byte("application/json"))
		// ctx.Response.Header.SetCanonical([]byte("Access-Control-Allow-Origin"), []byte("*"))
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		// ctx.Response.Header.SetCanonical([]byte("Access-Control-Allow-Methods"), []byte("GET, POST"))
		// ctx.Response.Header.SetCanonical([]byte("Access-Control-Allow-Headers"), []byte("Content-Type"))

		ctx.Response.SetStatusCode(200)

		requestBody := ctx.Request.Body()
		var allocRequest WalletAllocRPCCommand
		err := json.Unmarshal(requestBody, &allocRequest)

		if err != nil {
			handleError(ctx, err)
			return
		}

		response, err := alloc.Allocate(
			common.PublicKeyFromString(allocRequest.TokenMint),
			common.PublicKeyFromString(allocRequest.Owner),
		)
		if err != nil {
			handleError(ctx, err)
			return
		}

		encodedResult, err := json.Marshal(&response)
		if err != nil {
			handleError(ctx, err)
			return
		}
		ctx.Response.SetBody(encodedResult)

		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
}

