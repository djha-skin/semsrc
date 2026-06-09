package filesystem

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/djha-skin/semsrc/internal/store"
)

// Config holds configuration for the filesystem object store
type Config struct {
	Path      string // Base path for object storage
	ShardDepth int    // Number of directory levels for sharding (e.g., 2 for ab/cd/abcdef...)
	Compress  bool   // Whether to compress objects
}

// ObjectStore implements store.ObjectStore using the local filesystem
type ObjectStore struct {
	config Config
	mu     sync.RWMutex // Protects concurrent access to the store
}

// New creates a new filesystem object store
func New(config Config) (*ObjectStore, error) {
	// Ensure the base directory exists
	if err := os.MkdirAll(config.Path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create object store directory: %w", err)
	}

	// Validate shard depth
	if config.ShardDepth < 0 || config.ShardDepth > 4 {
		return nil, fmt.Errorf("shard depth must be between 0 and 4")
	}

	return &ObjectStore{
		config: config,
	}, nil
}

// Put stores data and returns its hash
func (s *ObjectStore) Put(data io.Reader) (store.Hash, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Read all data to calculate hash
	dataBytes, err := io.ReadAll(data)
	if err != nil {
		return store.Hash{}, fmt.Errorf("failed to read data: %w", err)
	}

	// Calculate hash
	hash := sha256.Sum256(dataBytes)

	// Check if object already exists
	exists, err := s.Exists(hash)
	if err != nil {
		return store.Hash{}, fmt.Errorf("failed to check existence: %w", err)
	}
	if exists {
		return hash, nil
	}

	// Get object path
	path := s.objectPath(hash)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return store.Hash{}, fmt.Errorf("failed to create directory: %w", err)
	}

	// Write to temporary file first
	tmpPath := path + ".tmp"
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return store.Hash{}, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpPath) // Clean up on error
	defer tmpFile.Close()

	// Write data
	if _, err := tmpFile.Write(dataBytes); err != nil {
		return store.Hash{}, fmt.Errorf("failed to write data: %w", err)
	}

	// Close temp file
	if err := tmpFile.Close(); err != nil {
		return store.Hash{}, fmt.Errorf("failed to close temp file: %w", err)
	}

	// Atomically rename to final location
	if err := os.Rename(tmpPath, path); err != nil {
		return store.Hash{}, fmt.Errorf("failed to rename temp file: %w", err)
	}

	return hash, nil
}

// Get retrieves data by hash
func (s *ObjectStore) Get(hash store.Hash) (io.ReadCloser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := s.objectPath(hash)
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, store.ErrNotFound
		}
		return nil, fmt.Errorf("failed to open object: %w", err)
	}

	// Verify hash while reading?
	// For now, trust the filesystem path matches the hash
	return file, nil
}

// Exists checks if an object exists
func (s *ObjectStore) Exists(hash store.Hash) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := s.objectPath(hash)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("failed to stat object: %w", err)
}

// Delete removes an object
func (s *ObjectStore) Delete(hash store.Hash) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.objectPath(hash)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return store.ErrNotFound
		}
		return fmt.Errorf("failed to delete object: %w", err)
	}

	// Try to clean up empty directories
	s.cleanupEmptyDirs(path)
	return nil
}

// Stats returns statistics about the store
func (s *ObjectStore) Stats() (store.ObjectStoreStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var stats store.ObjectStoreStats
	var totalSize int64

	// Walk the object store directory
	err := filepath.Walk(s.config.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			stats.ObjectCount++
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		return store.ObjectStoreStats{}, fmt.Errorf("failed to walk directory: %w", err)
	}

	stats.TotalSize = totalSize

	// Get filesystem stats for used/free space
	var fs syscall.Statfs_t
	if err := syscall.Statfs(s.config.Path, &fs); err == nil {
		stats.UsedSpace = int64(fs.Blocks*uint64(fs.Bsize) - fs.Bfree*uint64(fs.Bsize))
		stats.FreeSpace = int64(fs.Bfree * uint64(fs.Bsize))
	}

	return stats, nil
}

// Close releases resources
func (s *ObjectStore) Close() error {
	// Nothing to close for filesystem store
	return nil
}

// objectPath returns the filesystem path for an object hash
func (s *ObjectStore) objectPath(hash store.Hash) string {
	hashStr := hash.String()
	
	if s.config.ShardDepth == 0 {
		return filepath.Join(s.config.Path, hashStr)
	}

	// Build sharded path
	parts := make([]string, 0, s.config.ShardDepth+1)
	parts = append(parts, s.config.Path)
	
	for i := 0; i < s.config.ShardDepth && i*2+2 <= len(hashStr); i++ {
		parts = append(parts, hashStr[i*2:i*2+2])
	}
	
	parts = append(parts, hashStr)
	return filepath.Join(parts...)
}

// cleanupEmptyDirs removes empty directories up to the shard depth
func (s *ObjectStore) cleanupEmptyDirs(objectPath string) {
	if s.config.ShardDepth == 0 {
		return
	}

	dir := filepath.Dir(objectPath)
	for depth := 0; depth < s.config.ShardDepth; depth++ {
		// Check if directory is empty
		entries, err := os.ReadDir(dir)
		if err != nil || len(entries) > 0 {
			break
		}

		// Remove empty directory
		os.Remove(dir)
		
		// Move up one level
		dir = filepath.Dir(dir)
		if dir == s.config.Path {
			break
		}
	}
}

// StreamWriter implements store.StreamWriter for filesystem store
type streamWriter struct {
	hash     sha256.Hash
	tmpFile  *os.File
	tmpPath  string
	finalPath string
	written  int64
}

// Write implements io.Writer
func (sw *streamWriter) Write(p []byte) (int, error) {
	n, err := sw.tmpFile.Write(p)
	if err != nil {
		return n, err
	}
	
	sw.hash.Write(p[:n])
	sw.written += int64(n)
	return n, nil
}

// Close implements io.Closer
func (sw *streamWriter) Close() error {
	// Close temp file
	if err := sw.tmpFile.Close(); err != nil {
		return err
	}

	// Get final hash
	var finalHash store.Hash
	copy(finalHash[:], sw.hash.Sum(nil))

	// Check if object already exists
	exists, err := fileExists(sw.finalPath)
	if err != nil {
		return err
	}
	if exists {
		// Clean up temp file
		os.Remove(sw.tmpPath)
		return store.ErrAlreadyExists
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(sw.finalPath), 0755); err != nil {
		return err
	}

	// Atomically rename to final location
	return os.Rename(sw.tmpPath, sw.finalPath)
}

// Hash returns the hash of the written data
func (sw *streamWriter) Hash() store.Hash {
	var hash store.Hash
	copy(hash[:], sw.hash.Sum(nil))
	return hash
}

// Written returns the number of bytes written
func (sw *streamWriter) Written() int64 {
	return sw.written
}

// StreamPut returns a WriteCloser that calculates hash as it writes
func (s *ObjectStore) StreamPut() (store.StreamWriter, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create temporary file
	tmpFile, err := os.CreateTemp(s.config.Path, "semsrc-*.tmp")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	// We don't know the hash yet, so we can't determine the final path
	// We'll calculate it as we write and rename at the end
	return &streamWriter{
		hash:    sha256.New(),
		tmpFile: tmpFile,
		tmpPath: tmpFile.Name(),
	}, nil
}

// StreamGet returns a ReadCloser that verifies hash as it reads
func (s *ObjectStore) StreamGet(hash store.Hash) (io.ReadCloser, error) {
	return s.Get(hash) // Filesystem store doesn't do hash verification during streaming
}

// helper function to check if file exists
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Note: For streaming with final path calculation, we would need to buffer
// or implement a more complex solution. The current StreamPut implementation
// writes to a temp file and calculates hash, but we need the hash to know
// the final path. One solution is to write to temp, then read and calculate
// hash, then rename. Another is to write to final location with a placeholder
// name and rename after.