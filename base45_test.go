package base45

import (
	"bytes"
	"math/rand"
	"strings"
	"testing"
)

var validRfcExamples = []struct {
	decoded []byte
	encoded []byte
}{
	/*
		[1] Chapter 4.3:

		Encoding example 1:

		The string "AB" is the byte sequence [[65 66]].  If we look at all
		16 bits, we get 65 * 256 + 66 = 16706.  16706 equals 11 + (11 *
		45) + (8 * 45 * 45), so the sequence in base 45 is [11 11 8].
		Referring to Table 1, we get the encoded string "BB8".

		                  +-----------+------------------+
		                  | AB        | Initial string   |
		                  +-----------+------------------+
		                  | [[65 66]] | Decimal value    |
		                  +-----------+------------------+
		                  | [16706]   | Value in base 16 |
		                  +-----------+------------------+
		                  | [11 11 8] | Value in base 45 |
		                  +-----------+------------------+
		                  | BB8       | Encoded string   |
		                  +-----------+------------------+

		                    Table 2: Example 1 in Detail
	*/
	{[]byte("AB"), []byte("BB8")},

	/*
		[1] Chapter 4.3:

		Encoding example 2:

		The string "Hello!!" as ASCII is the byte sequence [[72 101] [108
		108] [111 33] [33]].  If we look at this 16 bits at a time, we get
		[18533 27756 28449 33].  Note the 33 for the last byte.  When
		looking at the values in base 45, we get [[38 6 9] [36 31 13] [9 2
		14] [33 0]], where the last byte is represented by two values.
		The resulting string "%69 VD92EX0" is created by looking up these
		values in Table 1.  It should be noted it includes a space.

		 +---------------------------------------+------------------+
		 | Hello!!                               | Initial string   |
		 +---------------------------------------+------------------+
		 | [[72 101] [108 108] [111 33] [33]]    | Decimal value    |
		 +---------------------------------------+------------------+
		 | [18533 27756 28449 33]                | Value in base 16 |
		 +---------------------------------------+------------------+
		 | [[38 6 9] [36 31 13] [9 2 14] [33 0]] | Value in base 45 |
		 +---------------------------------------+------------------+
		 | %69 VD92EX0                           | Encoded string   |
		 +---------------------------------------+------------------+

		                   Table 3: Example 2 in Detail
	*/
	{[]byte("Hello!!"), []byte("%69 VD92EX0")},

	/*
		[1] Chapter 4.3:

		Encoding example 3:

		The string "base-45" as ASCII is the byte sequence [[98 97] [115
		101] [45 52] [53]].  If we look at this two bytes at a time, we
		get [25185 29541 11572 53].  Note the 53 for the last byte.  When
		looking at the values in base 45, we get [[30 19 12] [21 26 14] [7
		32 5] [8 1]] where the last byte is represented by two values.
		Referring to Table 1, we get the encoded string "UJCLQE7W581".

		 +----------------------------------------+------------------+
		 | base-45                                | Initial string   |
		 +----------------------------------------+------------------+
		 | [[98 97] [115 101] [45 52] [53]]       | Decimal value    |
		 +----------------------------------------+------------------+
		 | [25185 29541 11572 53]                 | Value in base 16 |
		 +----------------------------------------+------------------+
		 | [[30 19 12] [21 26 14] [7 32 5] [8 1]] | Value in base 45 |
		 +----------------------------------------+------------------+
		 | UJCLQE7W581                            | Encoded string   |
		 +----------------------------------------+------------------+

						Table 4: Example 3 in Detail
	*/
	{[]byte("base-45"), []byte("UJCLQE7W581")},

	/*
		[1] Chapter 4.4:

		Decoding example 1:

		The string "QED8WEX0" represents, when looked up in Table 1, the
		values [26 14 13 8 32 14 33 0].  We arrange the numbers in chunks
		of three, except for the last one which can be two numbers, and
		get [[26 14 13] [8 32 14] [33 0]].  In base 45, we get [26981
		29798 33] where the bytes are [[105 101] [116 102] [33]].  If we
		look at the ASCII values, we get the string "ietf!".

		 +-------------------------------+------------------------+
		 | QED8WEX0                      | Initial string         |
		 +-------------------------------+------------------------+
		 | [26 14 13 8 32 14 33 0]       | Looked up values       |
		 +-------------------------------+------------------------+
		 | [[26 14 13] [8 32 14] [33 0]] | Groups of three        |
		 +-------------------------------+------------------------+
		 | [26981 29798 33]              | Interpreted as base 45 |
		 +-------------------------------+------------------------+
		 | [[105 101] [116 102] [33]]    | Values in base 8       |
		 +-------------------------------+------------------------+
		 | ietf!                         | Decoded string         |
		 +-------------------------------+------------------------+

						Table 5: Example 4 in Detail
	*/
	{[]byte("ietf!"), []byte("QED8WEX0")},
}

func TestValidEncodeWithRfcExamples(t *testing.T) {
	for _, entry := range validRfcExamples {
		got := Encode(entry.decoded)

		if !bytes.Equal(got, entry.encoded) {
			t.Errorf("Unexpected encoding result for \"%s\", expected %v, got %v", entry.decoded, entry.encoded, got)
		}
	}
}

func TestValidDecodeWithRfcExamples(t *testing.T) {
	for _, entry := range validRfcExamples {
		got, err := Decode(entry.encoded)

		if err != nil {
			t.Errorf("Expected decoded string, got error \"%s\"", err)
		}

		if !bytes.Equal(got, entry.decoded) {
			t.Errorf("Unexpected decoding result for \"%s\", expected %v, got %v", entry.encoded, entry.decoded, got)
		}
	}
}

func TestValidLargeEncodeDecode(t *testing.T) {
	// As the RFC examples are very slim, there is no chance
	// that the full alphabet gets tested, so we process 1mb
	// of random data to gain some confidence that no alphabet
	// errors are present during a encode/decode cycle.
	expected := make([]byte, 1048576)
	rand.Read(expected)

	enc := Encode(expected)
	got, err := Decode(enc)

	if err != nil {
		t.Errorf("Failed to decode the large set with %v", err)
	}

	if !bytes.Equal(got, expected) {
		t.Errorf("Decoded large set not equal to expected large set")
	}
}

func TestInvalidInputLengthDecode(t *testing.T) {
	_, err := Decode([]byte("ABCD"))

	if err != ErrInvalidLength {
		t.Errorf("Expected ErrInvalidLength, got \"%s\"", err)
	}
}

func TestInvalidEncodedInputAlphabet(t *testing.T) {
	_, err := Decode([]byte("aa"))

	if err != ErrInvalidEncodingCharacters {
		t.Errorf("Expected ErrInvalidEncodingCharacters, got \"%s\"", err)
	}
}

func TestInvalidOverflow(t *testing.T) {
	/*
		[1] Chapter 6:

		Even though a Base45-encoded string contains only characters from the
		alphabet in Table 1, cases like the following have to be considered:
		The string "FGW" represents 65535 (FFFF in base 16), which is a valid
		encoding of 16 bits.  A slightly different encoded string of the same
		length, "GGW", would represent 65536 (10000 in base 16), which is
		represented by more than 16 bits.  Implementations MUST also reject
		the encoded data if it contains a triplet of characters that, when
		decoded, results in an unsigned integer that is greater than 65535
		(FFFF in base 16).
	*/

	valid, err := Decode([]byte("FGW"))

	if err != nil {
		t.Errorf("Expected value, got error %s", err)
	}

	if !bytes.Equal(valid, []byte{255, 255}) {
		t.Errorf("Expected the valid decoded value to be 0xffff")
	}

	// Test 3 byte overflows
	_, err = Decode([]byte("GGW"))

	if err != ErrInvalidEncodedDataOverflow {
		t.Errorf("Expected ErrInvalidEncodedDataOverflow, got \"%s\"", err)
	}

	// Test 2 byte overflows
	_, err = Decode([]byte("::"))

	if err != ErrInvalidEncodedDataOverflow {
		t.Errorf("Expected ErrInvalidEncodedDataOverflow, got \"%s\"", err)
	}
}

func TestEmptyEncode(t *testing.T) {
	expected := []byte("")
	got := Encode([]byte(""))

	if !bytes.Equal(got, expected) {
		t.Errorf("Expected empty encode input to lead to %v, got %v", expected, got)
	}
}

func TestEmptyDecode(t *testing.T) {
	_, err := Decode([]byte{})

	if err != ErrEmptyInput {
		t.Errorf("Expected error on decode of empty value, got \"%v\"", err)
	}
}

func TestEncodeURLSafe(t *testing.T) {
	got := EncodeURLSafe([]byte("Hello!!"))

	expected := "%2569%20VD92EX0"
	if !strings.EqualFold(got, expected) {
		t.Errorf("Expected \"%s\", got \"%s\"", expected, got)
	}
}

func TestEmptyEncodeURLSafe(t *testing.T) {
	expected := []byte("")
	got := EncodeURLSafe([]byte(""))

	if !bytes.Equal([]byte(got), expected) {
		t.Errorf("Expected empty encode input to lead to %v, got %v", expected, got)
	}
}

func TestDecodeURLSafe(t *testing.T) {
	got, err := DecodeURLSafe("%2569%20VD92EX0")

	if err != nil {
		t.Errorf("Expected url decoded data, got error \"%s\"", err)
	}

	expected := []byte("Hello!!")

	if !bytes.Equal(got, expected) {
		t.Errorf("Expected \"%s\", got \"%s\"", expected, got)
	}
}

func TestEmptyDecodeURLSafe(t *testing.T) {
	_, err := DecodeURLSafe("")

	if err != ErrEmptyInput {
		t.Errorf("Expected error on decode of empty value, got \"%v\"", err)
	}
}

func TestInvalidDecodeURLSafe(t *testing.T) {
	_, err := DecodeURLSafe("%20%")

	if err != ErrInvalidURLSafeEscaping {
		t.Errorf("Expected url decode error, got %v", err)
	}
}
