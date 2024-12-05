package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/joho/godotenv"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

func main() {

	_ = godotenv.Load("../.env")
	words := strings.Split("piece panel addict sniff scare muffin never stem finish immune ship antenna pigeon toast evolve midnight useless cash blade mandate pitch cabin gasp lamp", " ")

	client := liteclient.NewConnectionPool()

	configUrl := "https://ton-blockchain.github.io/testnet-global.config.json"
	err := client.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		panic(err)
	}

	api := ton.NewAPIClient(client)

	block, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		panic(err)
	}

	w, err := wallet.FromSeed(api, words, wallet.V4R2)
	if err != nil {
		panic(err)
	}

	balance, err := w.GetBalance(context.Background(), block)
	if err != nil {
		panic(err)
	}
	account, err := api.GetAccount(context.Background(), block, w.WalletAddress())
	if err != nil {
		panic(err)
	}

	text := fmt.Sprintf("pay-%d-%d", 1200, 1)
	text = url.PathEscape(text)
	fmt.Println(
		fmt.Sprintf(
			`https://app.tonkeeper.com/transfer/%s?text=%s?&amount=%s`,
			w.WalletAddress().String(),
			text,
			tlb.MustFromTON("0.1").Nano().String()),
	)

	fmt.Println(balance)
	// fmt.Println(w.WalletAddress())
	lastLt := account.LastTxLT
	lastHash := account.LastTxHash

	dbLastlt := uint64(28748267000001)

mainFor:
	for {
		list, err := api.ListTransactions(context.Background(), w.WalletAddress(), 1, lastLt, lastHash)
		if err != nil {
			panic(err)
		}
		lastLt = list[0].PrevTxLT
		lastHash = list[0].PrevTxHash

		for _, transaction := range list {
			if transaction.LT <= dbLastlt {
				break mainFor
			}
			amount := transaction.IO.In.AsInternal().Amount
			comment := transaction.IO.In.AsInternal().Comment()

			fmt.Println(transaction)
			fmt.Println(comment, amount)
		}
	}

}
