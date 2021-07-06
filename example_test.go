package base45

import "fmt"

func ExampleEncode() {
	encoded := Encode([]byte("Hello!!"))
	fmt.Printf("Encoded: %s\n", encoded)
}

func ExampleEncodeURLSafe() {
	encoded := EncodeURLSafe([]byte("Hello!!"))
	fmt.Printf("Encoded url safe: %s\n", encoded)
}

func ExampleDecode() {
	decoded, err := Decode([]byte("%69 VD92EX0"))
	fmt.Printf("Decoded: %s, Error: %v\n", decoded, err)
}

func ExampleDecodeURLSafe() {
	decoded, err := DecodeURLSafe("%2569%20VD92EX0")
	fmt.Printf("Decoded url safe: %s, Error: %v\n", decoded, err)
}

func ExampleDecode_errorHandling() {
	_, err := Decode([]byte("GGW"))

	if err == ErrInvalidEncodedDataOverflow {
		fmt.Printf("Encountered invalid data")
	}
}
