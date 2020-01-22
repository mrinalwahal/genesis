package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/genesis/blockchain"
)

func PrintChain(chain *blockchain.BlockChain) {
	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Hash %x\n", block.Hash)
		fmt.Printf("Data: %s\n", block.Data)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}

}

func main() {

	defer os.Exit(0)
	address := "wahal"
	chainI := blockchain.InitBlockChain(address)
	defer chainI.Database.Close()

	chain := blockchain.ContinueBlockChain(address)

	chain.AddBlock("Ekam mudra")
	chain.AddBlock("Dvi mudra")
	chain.AddBlock("Tritye mudra")
	chain.AddBlock("Chaturiya mudra")

	PrintChain(chain)
}
