package base45

import "fmt"

func ExampleEncode() {
	encoded := Encode([]byte("Hello!!"))
	fmt.Printf("Encoded: %s", encoded)
}

func ExampleEncodeURLSafe() {
	encoded := EncodeURLSafe([]byte("Hello!!"))
	fmt.Printf("Encoded url safe: %s", encoded)
}

func ExampleDecode() {
	decoded, _ := Decode([]byte("%69 VD92EX0"))
	fmt.Printf("Decoded: %s", decoded)
}

func ExampleDecodeURLSafe() {
	decoded, _ := DecodeURLSafe("%2569%20VD92EX0")
	fmt.Printf("Decoded url safe: %s", decoded)
}

func ExampleDecode_errorHandling() {
	_, err := Decode([]byte("GGW"))

	if err == ErrInvalidEncodedDataOverflow {
		fmt.Printf("Encountered invalid data")
	}
}
