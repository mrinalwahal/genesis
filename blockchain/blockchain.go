package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

// BlockChain designs the chhain of type Blocks
type BlockChain struct {
	//Blocks []*Block

	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

/*
// DeriveHash calculates and returns the hash
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}
*/

// AddBlock adds a new block to the chain
func (chain *BlockChain) AddBlock(data string) {

	//prevBlock := chain.Blocks[len(chain.Blocks)-1]
	//new := CreateBlock(data, prevBlock.Hash)
	//chain.Blocks = append(chain.Blocks, new)

	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			// This func with val would only be called if item.Value encounters no error.

			// Accessing val here is valid.
			//fmt.Printf("The answer is: %s\n", val)

			// Copying or parsing val is valid.
			lastHash = append([]byte{}, val...)
			return nil
		})
		return err
	})

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err)

}

// InitBlockChain creates thhe chain
func InitBlockChain() *BlockChain {
	var lastHash []byte

	db, err := badger.Open(badger.DefaultOptions(dbPath))
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existiinig blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis awakened and proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)

			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash
			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			err = item.Value(func(val []byte) error {
				// This func with val would only be called if item.Value encounters no error.

				// Accessing val here is valid.
				//fmt.Printf("The answer is: %s\n", val)

				// Copying or parsing val is valid.
				lastHash = append([]byte{}, val...)
				return nil
			})
			return err
		}
	})

	//return &BlockChain{[]*Block{Genesis()}}

	Handle(err)
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}
