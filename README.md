# SuSy Wallet Allocator

Wallet allocator is a distinct oracle that creates token account for specified token mint and owner. If the token account exists it returns it. Otherwise, creates a new one.

## Executing

1. Run `go build -o allocator`
2. Execute `./allocator --port <port> --keypair <path to solana keypair file>`. If you don't have a keypair. 

You can create one via Solana CLI `solana-keygen new --no-bip39-passphrase -o ./allocator.json`.

## Reference

1. [Associated Token Account](https://spl.solana.com/associated-token-account)
2. [Solanoid](https://github.com/SuSy-One/solanoid)
3. [Solana](https://github.com/solana-labs/solana)
