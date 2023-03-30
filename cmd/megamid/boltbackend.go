package main

import (
	"strconv"

	"go.etcd.io/bbolt"
)

type BoltBackend struct {
	BucketName string
	SizeBytes  int64

	db *bbolt.DB
}

func toBlockID(block int64) string {
	return strconv.FormatInt(block, 16)
}

func (bb *BoltBackend) getOrCreateBlock(bkt *bbolt.Bucket, blockID string) ([]byte, error) {
	result := bkt.Get([]byte(blockID))
	if result != nil {
		return result, nil
	}

	result = make([]byte, blockSize)

	return result, bkt.Put([]byte(blockID), result)
}

func (bb *BoltBackend) ReadAt(p []byte, off int64) (int, error) {
	size := len(p)
	nBlocks := size / blockSize

	var n int
	if err := bb.db.Update(func(tx *bbolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(bb.BucketName))
		if err != nil {
			return err
		}

		for i := 0; i <= nBlocks; i++ {
			id := off + int64(i*blockSize)
			block, err := bb.getOrCreateBlock(bkt, toBlockID(id))
			if err != nil {
				return err
			}

			n += copy(p[i*blockSize:((i+1)*blockSize)-1], block)
		}

		return nil
	}); err != nil {
		return -1, err
	}
	return n, nil
}

func (bb *BoltBackend) WriteAt(p []byte, off int64) (int, error) {
	size := len(p)
	nBlocks := size / blockSize

	if err := bb.db.Update(func(tx *bbolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(bb.BucketName))
		if err != nil {
			return err
		}

		for i := 0; i <= nBlocks; i++ {
			id := off + int64(i*blockSize)
			block := p[i*blockSize : ((i+1)*blockSize)-1]
			if err := bkt.Put([]byte(toBlockID(id)), block); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return -1, err
	}

	return size, nil
}

func (bb *BoltBackend) Size() (int64, error) {
	return bb.SizeBytes, nil
}

func (bb *BoltBackend) Sync() error {
	return bb.db.Sync()
}
