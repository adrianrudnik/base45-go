package base45

import "errors"

// ErrInvalidEncodingCharacters means that the encoded string did contain
// invalid characters not supported by the base 45 alphabet.
var ErrInvalidEncodingCharacters = errors.New("invalid characters in encoded string")

// ErrInvalidLength means that the encoded string did not match the expected
// length. Normal tuples have a length of 3 bytes, shortend (trailing) ones of 2 bytes.
// If you encounter an base 45 encoded string with a length of of 1 or 4 bytes
// it can not be valid.
var ErrInvalidLength = errors.New("invalid input length")

// ErrInvalidURLSafeEscaping means the given URL escaped content could not safely be unescaped.
var ErrInvalidURLSafeEscaping = errors.New("invalid escaped input given")

// ErrInvalidEncodedDataOverflow means the decoder encountered an invalid byte combination
// like "GGW" which would lead to an overflow of a uint16 (with the value 0xffff + 1).
var ErrInvalidEncodedDataOverflow = errors.New("invalid encoded data leads to unexpected overflow")
