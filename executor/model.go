package executor

import "encoding/json"




type WalletAllocRPCCommand struct {
	TokenMint string `json:"token_mint"`
	Owner string `json:"owner"`
}

type WalletAllocResponse struct {
	PublicKey string `json:"public_key"`
	TokenMint string `json:"token_mint"`
	TxSignature string `json:"tx_signature"`
}


type Response struct {
	Status string `json:"status"`
	Result string `json:"result"`
}


func (r *Response) Serialize() string {
	data, _ := json.Marshal(r)
	return string(data)
}

type TokenAccountsByOwnerResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Context struct {
			Slot int `json:"slot"`
		} `json:"context"`
		Value []struct {
			Account struct {
				Data struct {
					Parsed struct {
						Info struct {
							IsNative    bool   `json:"isNative"`
							Mint        string `json:"mint"`
							Owner       string `json:"owner"`
							State       string `json:"state"`
							TokenAmount struct {
								Amount         string `json:"amount"`
								Decimals       int    `json:"decimals"`
								UIAmount       float64    `json:"uiAmount"`
								UIAmountString string `json:"uiAmountString"`
							} `json:"tokenAmount"`
						} `json:"info"`
						Type string `json:"type"`
					} `json:"parsed"`
					Program string `json:"program"`
					Space   int    `json:"space"`
				} `json:"data"`
				Executable bool   `json:"executable"`
				Lamports   int    `json:"lamports"`
				Owner      string `json:"owner"`
				RentEpoch  int    `json:"rentEpoch"`
			} `json:"account"`
			Pubkey string `json:"pubkey"`
		} `json:"value"`
	} `json:"result"`
	ID int `json:"id"`
}