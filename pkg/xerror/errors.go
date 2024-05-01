package xerror

import "errors"

var (
	ErrCipherEncryptionFailure  = errors.New("cipher encryption failure")
	ErrCipherMethodNotSupported = errors.New("cipher method not supported")
	ErrNilFrame                 = errors.New("nil frame")
)
