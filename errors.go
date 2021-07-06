package base45

import "errors"

var InvalidEncodingCharactersError = errors.New("invalid characters in encoded string")
var InvalidLengthError = errors.New("invalid input length")
var InvalidUrlSafeEscapingError = errors.New("invalid escaped input given")
var InvalidEncodedDataOverflowError = errors.New("invalid encoded data leads to unexpected overflow")
