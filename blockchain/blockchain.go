package blockchain

import (
	"encoding/hex"
	"log"
	"os"
	"runtime"

	"github.com/prologic/bitcask"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/MANIFEST"
	genesisData = "First txn from Genesis"
)

// BlockChain designs the chhain of type Blocks
type BlockChain struct {
	//Blocks []*Block

	LastHash []byte
	Database *bitcask.Bitcask
}

// BlockChainIterator runs the iter for blockchain
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *bitcask.Bitcask
}

/*
// DeriveHash calculates and returns the hash
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}
*/

func DBExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// AddBlock adds a new block to the chain
func (chain *BlockChain) AddBlock(transactions []*Transaction) {

	//prevBlock := chain.Blocks[len(chain.Blocks)-1]
	//new := CreateBlock(data, prevBlock.Hash)
	//chain.Blocks = append(chain.Blocks, new)

	var lastHash []byte

	item, err := chain.Database.Get([]byte("last_hash"))
	Handle(err)
	lastHash = item

	newBlock := CreateBlock(transactions, lastHash)

	err = chain.Database.Put(newBlock.Hash, newBlock.Serialize())
	Handle(err)

	err = chain.Database.Put([]byte("last_hash"), newBlock.Hash)
	Handle(err)

	chain.LastHash = newBlock.Hash

}

// InitBlockChain creates thhe chain
func InitBlockChain(address string) *BlockChain {

	if DBExists() {
		log.Println("existing blockchain detected")
		runtime.Goexit()
	}

	var lastHash []byte

	db, err := bitcask.Open(dbPath)
	Handle(err)

	item, err := db.Get([]byte("last_hash"))

	if err != nil {
		cbtx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		log.Println("Genesis awakened and proved")

		//
		err := db.Put(genesis.Hash, genesis.Serialize())
		Handle(err)

		err = db.Put([]byte("last_hash"), genesis.Hash)
		lastHash = genesis.Hash
		Handle(err)
	}

	lastHash = item

	//return &BlockChain{[]*Block{Genesis()}}

	Handle(err)
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func ContinueBlockChain(address string) *BlockChain {

	if !DBExists() {
		log.Println("no blockchain detected")
		runtime.Goexit()
	}

	var lastHash []byte

	db, err := bitcask.Open(dbPath)
	Handle(err)

	lastHash, err = db.Get([]byte("last_hash"))
	Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}
	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block
	//var encodedBlock []byte

	encodedBlock, err := iter.Database.Get(iter.CurrentHash)
	Handle(err)
	block = Deserialize(encodedBlock)

	iter.CurrentHash = block.PrevHash
	return block
}

func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction

	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, txn := range block.Transactions {
			txnID := hex.EncodeToString(txn.ID)

		Outputs:
			for outIdx, out := range txn.Outputs {
				if spentTXOs[txnID] != nil {
					for _, spentOut := range spentTXOs[txnID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *txn)
				}
			}
			if !txn.IsCoinbase() {
				for _, in := range txn.Inputs {
					if in.CanUnlock(address) {
						inTxnID := hex.EncodeToString(in.ID)
						spentTXOs[inTxnID] = append(spentTXOs[inTxnID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxs
}

func (chain *BlockChain) FindUTXNO(address string) []TxOutput {
	var UTXOs []TxOutput

	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, txn := range unspentTransactions {
		for _, out := range txn.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxns := chain.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, txn := range unspentTxns {
		txnID := hex.EncodeToString(txn.ID)

		for outIdx, out := range txn.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txnID] = append(unspentOuts[txnID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}
