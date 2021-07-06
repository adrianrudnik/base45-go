# Base45 encoder/decoder for Go

A small and simple package to encode and decode Base45 data.

[![license](https://img.shields.io/github/license/adrianrudnik/base45-go.svg)](https://github.com/adrianrudnik/base45-go/blob/main/LICENSE)
[![lint and test](https://github.com/adrianrudnik/base45-go/actions/workflows/test.yaml/badge.svg)](https://github.com/adrianrudnik/base45-go/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/adrianrudnik/base45-go/branch/main/graph/badge.svg?token=O4B0TOQHM0)](https://codecov.io/gh/adrianrudnik/base45-go)
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

The current implementation works with the following benchmark results:

```
cpu: Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz
BenchmarkEncodeByte1-8          53444944                22.16 ns/op            2 B/op          1 allocs/op
BenchmarkEncodeByte8-8          14158418                84.98 ns/op           16 B/op          1 allocs/op
BenchmarkEncodeByte64-8          2111271               573.7 ns/op           128 B/op          1 allocs/op
BenchmarkEncodeByte512-8          270361              4493 ns/op            1024 B/op          1 allocs/op
BenchmarkEncodeByte1024-8         134451              8881 ns/op            2048 B/op          1 allocs/op
BenchmarkEncodeByte8192-8          17013             70445 ns/op           16384 B/op          1 allocs/op
BenchmarkDecodeChunk1-8         28931134                40.61 ns/op            2 B/op          1 allocs/op
BenchmarkDecodeChunk8-8          5612131               221.2 ns/op            21 B/op          5 allocs/op
BenchmarkDecodeChunk64-8          692574              1696 ns/op             160 B/op         33 allocs/op
BenchmarkDecodeChunk512-8          80316             15270 ns/op            1280 B/op        257 allocs/op
BenchmarkDecodeChunk1024-8         38046             30245 ns/op            2560 B/op        513 allocs/op
BenchmarkDecodeChunk8192-8          4898            243912 ns/op           20480 B/op       4097 allocs/op
```
