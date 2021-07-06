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

	urlEncoded := base45.EncodeUrlSafe([]byte("Hello!!"))
	fmt.Printf("Encoded url safe: %s\n", urlEncoded)

	// Decoding data
	decoded, err := base45.Decode([]byte("%69 VD92EX0"))
	fmt.Printf("Decoded: %s, Error: %v\n", decoded, err)

	urlDecoded, err := base45.DecodeUrlSafe("%2569%20VD92EX0")
	fmt.Printf("Decoded url safe: %s, Error: %v\n", urlDecoded, err)

	// Error handling
	_, err = base45.Decode([]byte("GGW"))

	if err == base45.InvalidEncodedDataOverflowError {
		fmt.Printf("Encountered invalid data")
	}
}
```

## Performance

The current implementation works with the following benchmark results:

```
cpu: Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz
BenchmarkEncode
BenchmarkEncode-8          	138791397	         8.723 ns/op
BenchmarkEncodeURLSafe
BenchmarkEncodeURLSafe-8   	36777698	        32.39 ns/op
BenchmarkDecode
BenchmarkDecode-8          	25810398	        46.26 ns/op
BenchmarkDecodeURLSafe
BenchmarkDecodeURLSafe-8   	20122197	        58.18 ns/op
```
