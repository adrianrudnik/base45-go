# Base45 encoder/decoder for Go

A small and simple package to encode and decode Base45 data.

[![license](https://img.shields.io/github/license/adrianrudnik/base45-go.svg)](https://github.com/adrianrudnik/base45-go/blob/main/LICENSE)
[![lint and test](https://github.com/adrianrudnik/base45-go/actions/workflows/test.yaml/badge.svg)](https://github.com/adrianrudnik/base45-go/actions/workflows/test.yaml)
[![coverage](https://codecov.io/gh/adrianrudnik/base45-go/branch/main/graph/badge.svg?token=O4B0TOQHM0)](https://codecov.io/gh/adrianrudnik/base45-go)
[![go report card](https://goreportcard.com/badge/github.com/adrianrudnik/base45-go)](https://goreportcard.com/report/github.com/adrianrudnik/base45-go)

It implements the current [draft version 7](https://datatracker.ietf.org/doc/draft-faltstrom-base45/) with testing and security checks.

The target was not to optimize performance but keep the code readable and aligned to the draft.

## Usage

The following example app illustrates the usage:

```go
package main

import (
	"fmt"
	"github.com/adrianrudnik/base45-go"
)

func main() {
	// Encoding data
	encoded := base45.Encode([]byte("Hello!!"))
	fmt.Printf("Encoded: %s\n", encoded)

	urlEncoded := base45.EncodeURLSafe([]byte("Hello!!"))
	fmt.Printf("Encoded url safe: %s\n", urlEncoded)

	// Decoding data
	decoded, err := base45.Decode([]byte("%69 VD92EX0"))
	fmt.Printf("Decoded: %s, Error: %v\n", decoded, err)

	urlDecoded, err := base45.DecodeURLSafe("%2569%20VD92EX0")
	fmt.Printf("Decoded url safe: %s, Error: %v\n", urlDecoded, err)

	// Error handling
	_, err = base45.Decode([]byte("GGW"))

	if err == base45.ErrInvalidEncodedDataOverflow {
		fmt.Printf("Encountered invalid data")
	}
}
```

## Performance

Encoding is measured on input bytes. Decoding is measured on output bytes.

```
cpu: Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz
BenchmarkEncode1-8              52604754                22.48 ns/op            2 B/op          1 allocs/op
BenchmarkEncode128-8             1000000              1143 ns/op             256 B/op          1 allocs/op
BenchmarkEncode512-8              268015              4453 ns/op            1024 B/op          1 allocs/op
BenchmarkEncode1024-8             135018              8905 ns/op            2048 B/op          1 allocs/op
BenchmarkEncode8192-8              16906             70784 ns/op           16384 B/op          1 allocs/op
BenchmarkDecode1-8              28835797                40.95 ns/op            2 B/op          1 allocs/op
BenchmarkDecode128-8              337608              3614 ns/op             320 B/op         65 allocs/op
BenchmarkDecode512-8               80144             15233 ns/op            1280 B/op        257 allocs/op
BenchmarkDecode1024-8              39213             30357 ns/op            2560 B/op        513 allocs/op
BenchmarkDecode8192-8               4642            245702 ns/op           20480 B/op       4097 allocs/op
```
