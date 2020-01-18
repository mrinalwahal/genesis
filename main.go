package main

import (
	"fmt"
	"strconv"

	"github.com/genesis/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()

	chain.AddBlock("Ekam mudra")
	chain.AddBlock("Dvi mudra")
	chain.AddBlock("Tritye mudra")

	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Hash %x\n", block.Hash)
		fmt.Printf("Data: %s\n", block.Data)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
