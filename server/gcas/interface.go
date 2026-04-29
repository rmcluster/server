// The GCAS package provides content-addressable storage systems.

package gcas

import (
	"context"
	"crypto/sha256"
)

type Hash [sha256.Size]byte

type CAS interface {
	Put(ctx context.Context, hash Hash, data []byte) error
	Get(ctx context.Context, hash Hash) ([]byte, error)
	Delete(ctx context.Context, hash Hash) error
	List(ctx context.Context) (<-chan Hash, error)
	FreeSpace(ctx context.Context) (int64, error)
}
