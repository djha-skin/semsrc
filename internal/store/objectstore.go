package store

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

// Hash represents a SHA-256 hash of object content
type Hash [sha256.Size]byte

// String returns the hex string representation of the hash
func (h Hash) String() string {
	return fmt.Sprintf("%x", h[:])
}

// ParseHash parses a hex string into a Hash
func ParseHash(s string) (Hash, error) {
	var h Hash
	if len(s) != sha256.Size*2 {
		return h, errors.New("invalid hash length")
	}
	_, err := fmt.Sscanf(s, "%x", &h)
	return h, err
}

// ObjectStore defines the interface for content-addressable storage
type ObjectStore interface {
	// Put stores data and returns its hash
	Put(data io.Reader) (Hash, error)
	
	// Get retrieves data by hash
	Get(hash Hash) (io.ReadCloser, error)
	
	// Exists checks if an object exists
	Exists(hash Hash) (bool, error)
	
	// Delete removes an object (may be no-op for immutable stores)
	Delete(hash Hash) error
	
	// Stats returns statistics about the store
	Stats() (ObjectStoreStats, error)
	
	// Close releases resources
	Close() error
}

// ObjectStoreStats contains statistics about the object store
type ObjectStoreStats struct {
	ObjectCount int64
	TotalSize   int64
	UsedSpace   int64
	FreeSpace   int64
}

// ObjectEntry represents a single object for batch operations
type ObjectEntry struct {
	Hash Hash
	Data io.Reader
	Size int64 // Optional size hint for optimization
}

// BatchObjectStore extends ObjectStore with batch operations for efficiency
type BatchObjectStore interface {
	ObjectStore
	
	// PutBatch stores multiple objects in a single operation
	PutBatch(entries []ObjectEntry) error
	
	// GetBatch retrieves multiple objects in a single operation
	GetBatch(hashes []Hash) (map[Hash]io.ReadCloser, error)
	
	// ExistsBatch checks existence of multiple objects
	ExistsBatch(hashes []Hash) (map[Hash]bool, error)
}

// StreamingObjectStore extends ObjectStore with streaming support
type StreamingObjectStore interface {
	ObjectStore
	
	// StreamPut returns a WriteCloser that calculates hash as it writes
	StreamPut() (StreamWriter, error)
	
	// StreamGet returns a ReadCloser that verifies hash as it reads
	StreamGet(hash Hash) (io.ReadCloser, error)
}

// StreamWriter writes data while calculating its hash
type StreamWriter interface {
	io.WriteCloser
	
	// Hash returns the hash of the written data
	Hash() Hash
	
	// Written returns the number of bytes written
	Written() int64
}

// Common errors
var (
	ErrNotFound      = errors.New("object not found")
	ErrHashMismatch  = errors.New("hash mismatch")
	ErrAlreadyExists = errors.New("object already exists")
	ErrInvalidHash   = errors.New("invalid hash")
)

// CalculateHash computes the SHA-256 hash of data
func CalculateHash(data []byte) Hash {
	return sha256.Sum256(data)
}

// CalculateHashFromReader computes hash from a reader
func CalculateHashFromReader(r io.Reader) (Hash, error) {
	h := sha256.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return Hash{}, err
	}
	var hash Hash
	copy(hash[:], h.Sum(nil))
	return hash, nil
}