package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Block represents the chain's block
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

// CreateBlock creates a new block
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}

	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	//block.DeriveHash()
	return block
}

// Genesis creates the first block in the chain
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func (b *Block) Serialize() []byte {

	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)
	Handle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {

	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	Handle(err)

	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
