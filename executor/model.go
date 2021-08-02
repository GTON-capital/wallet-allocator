package executor



type WalletAllocResponse struct {
	PublicKey string `json:"public_key"`
	TokenMint string `json:"token_mint"`
	TxSignature string `json:"tx_signature"`
}