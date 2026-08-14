package validators

import (
	"regexp"
	"strings"

	"blob/src/functions"

	"github.com/asaskevich/govalidator"
)

// UploadFields are the raw form values for a single-file upload.
type UploadFields struct {
	Bucket    string
	Filename  string
	Public    string
	ExpiresAt string
	Metadata  string
}

var bucketPattern = regexp.MustCompile(`^[a-zA-Z0-9/_-]+$`)

// ValidateUploadFields validates multipart/form-data upload fields and returns
// a map of field name to error message. An empty map means valid.
func ValidateUploadFields(f UploadFields) map[string]string {
	errors := make(map[string]string)

	if govalidator.IsNull(f.Bucket) {
		errors["bucket"] = "bucket is required"
	} else if !govalidator.StringLength(f.Bucket, "1", "64") {
		errors["bucket"] = "bucket must be 1-64 chars"
	} else if strings.Contains(f.Bucket, "..") {
		errors["bucket"] = "bucket cannot contain '..'"
	} else if !bucketPattern.MatchString(f.Bucket) {
		errors["bucket"] = "bucket can only contain letters, numbers, '/', '-', and '_'"
	}

	if f.Filename != "" && !govalidator.StringLength(f.Filename, "1", "255") {
		errors["filename"] = "filename must be 1-255 chars"
	}

	if f.Public != "" && !functions.StringInSlice(f.Public, []string{"true", "false", "0", "1"}) {
		errors["public"] = "public must be true, false, 0 or 1"
	}

	if f.ExpiresAt != "" && !govalidator.IsRFC3339(f.ExpiresAt) {
		errors["expires_at"] = "expires_at must be RFC3339 date"
	}

	if f.Metadata != "" && !govalidator.IsJSON(f.Metadata) {
		errors["metadata"] = "metadata must be valid JSON"
	}

	return errors
}

// BucketName reports whether the bucket name is acceptable for multipart uploads.
var multipartBucketPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,63}$`)

// ValidateMultipartBucket validates a bucket name for multipart sessions.
func ValidateMultipartBucket(bucket string) bool {
	if !multipartBucketPattern.MatchString(bucket) {
		return false
	}
	return !strings.Contains(bucket, "..")
}
