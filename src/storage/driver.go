package storage

import (
	"context"
	"io"
	"time"
)

// StorageDriver abstracts the physical persistence of blob data. It allows the
// service layer to operate on buckets/keys without caring about where the data
// actually lives (filesystem today, S3/MinIO/R2 in the future).
type StorageDriver interface {
	// Name returns the driver identifier (e.g. "filesystem").
	Name() string

	// Put stores a stream under bucket/key. size may be 0 when unknown.
	Put(ctx context.Context, bucket, key string, r io.Reader, size int64) error

	// Create opens a writer for bucket/key, overwriting any existing object.
	// Used by the multipart assembler to concatenate chunks.
	Create(ctx context.Context, bucket, key string) (io.WriteCloser, error)

	// Open returns a seekable reader for the object, used for streaming
	// (Range requests) and read-only access.
	Open(ctx context.Context, bucket, key string) (io.ReadSeekCloser, error)

	// Stat returns metadata about the object. Returns os.ErrNotExist when missing.
	Stat(ctx context.Context, bucket, key string) (*FileInfo, error)

	// Exists reports whether the object is present.
	Exists(ctx context.Context, bucket, key string) (bool, error)

	// Delete removes a single object. Missing objects are not an error.
	Delete(ctx context.Context, bucket, key string) error

	// Move relocates an object between bucket/key pairs.
	Move(ctx context.Context, srcBucket, srcKey, dstBucket, dstKey string) error

	// EnsureBucket creates the bucket namespace if it does not exist.
	EnsureBucket(ctx context.Context, bucket string) error

	// RemoveTree deletes an object prefix recursively (used for temp cleanup).
	RemoveTree(ctx context.Context, bucket, keyPrefix string) error
}

// FileInfo describes an object stored by a StorageDriver.
type FileInfo struct {
	Size    int64
	ModTime time.Time
}
