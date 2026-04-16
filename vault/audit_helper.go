package vault

import "bytes"

// newBytesReader wraps a byte slice in a *bytes.Reader for use with json.NewDecoder.
func newBytesReader(data []byte) *bytes.Reader {
	return bytes.NewReader(data)
}
