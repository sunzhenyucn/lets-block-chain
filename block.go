package lbc

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

// built-in block version
const (
	blockVersion = 1
)

// Block represents a block in block-chain
type Block struct {
	// previous block 32bit hash
	Prev []byte
	// 32bit hash for this block
	Hash []byte
	// created time for this block
	Timestamp int64
	// Transactions for this block
	Transactions []*Transaction
	// block version number
	Version int32
	// block nonce
	Nonce int
}

// newGenesisBlock will generate a genesis block,
// block data will be constant data,
// genesis block's prev is nil
func newGenesisBlock(coinBaseTx *Transaction) *Block {
	return NewBlock([]*Transaction{coinBaseTx}, []byte{})
}

// NewBlock will create a new block with spec data and prev block hash
func NewBlock(transactions []*Transaction, prev []byte) *Block {
	_block := &Block{
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		Prev:         prev,
		Hash:         []byte{},
		Version:      blockVersion,
		Nonce:        0,
	}

	_pow := NewProofWork(_block)
	nonce, hash := _pow.Run()

	_block.Nonce = nonce
	_block.Hash = hash
	return _block
}

// serialize the block struct to bytes array
func (b *Block) serialize() []byte {
	var blobBuf bytes.Buffer
	enc := gob.NewEncoder(&blobBuf)

	if err := enc.Encode(b); err != nil {
		panic(fmt.Sprintf("Serialized block failed, error with:`%s`\n", err))
	}

	return blobBuf.Bytes()
}

// hash all transactions in this block
func (b *Block) hashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash := sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

// deserializeBlock to block struct with spec bytes array
func deserializeBlock(blob []byte) *Block {
	var b Block
	dec := gob.NewDecoder(bytes.NewReader(blob))

	if err := dec.Decode(&b); err != nil {
		panic(fmt.Sprintf("Deserialized block failed, error with:`%s`", err))
	}

	return &b
}
