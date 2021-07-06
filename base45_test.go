package base45

import (
	"bytes"
	"math/rand"
	"strings"
	"testing"
)

var validDraftExamples = []struct {
	decoded []byte
	encoded []byte
}{
	/*
		[1] Chapter 4.3:

		Encoding example 1: The string "AB" is the byte sequence [65 66].
		The 16 bit value is 65 * 256 + 66 = 16706. 16706 equals 11 + 45 * 11
		+ 45 * 45 * 8 so the sequence in base 45 is [11 11 8].  By looking up
		these values in the Table 1 we get the encoded string "BB8".
	*/
	{[]byte("AB"), []byte("BB8")},

	/*
		[1] Chapter 4.3:

		Encoding example 2: The string "Hello!!" as ASCII is the byte
		sequence [72 101 108 108 111 33 33].  If we look at each 16 bit
		value, it is [18533 27756 28449 33].  Note the 33 for the last byte.
		When looking at the values modulo 45, we get [[38 6 9] [36 31 13] [9
		2 14] [33 0]] where the last byte is represented by two.  By looking
		up these values in the Table 1 we get the encoded string "%69
		VD92EX0".
	*/
	{[]byte("Hello!!"), []byte("%69 VD92EX0")},

	/*
		[1] Chapter 4.3:

		Encoding example 3: The string "base-45" as ASCII is the byte
		sequence [98 97 115 101 45 52 53].  If we look at each 16 bit value,
		it is [25185 29541 11572 53].  Note the 53 for the last byte.  When
		looking at the values modulo 45, we get [[30 19 12] [21 26 14] [7 32
		5] [8 1]] where the last byte is represented by two.  By looking up
		these values in the Table 1 we get the encoded string "UJCLQE7W581".
	*/
	{[]byte("base-45"), []byte("UJCLQE7W581")},

	/*
		[1] Chapter 4.4:

		Decoding example 1: The string "QED8WEX0" represents, when looked up
		in Table 1, the values [26 14 13 8 32 14 33 0].  We arrange the
		numbers in chunks of three, except for the last one which can be two,
		and get [[26 14 13] [8 32 14] [33 0]].  In base 45 we get [26981
		29798 33] where the bytes are [[105 101] [116 102] [33]].  If we look
		at the ASCII values we get the string "ietf!".
	*/
	{[]byte("ietf!"), []byte("QED8WEX0")},
}

func TestValidEncodeWithDraftExamples(t *testing.T) {
	for _, entry := range validDraftExamples {
		got := Encode(entry.decoded)

		if !bytes.Equal(got, entry.encoded) {
			t.Errorf("Unexpected encoding result for \"%s\", expected %v, got %v", entry.decoded, entry.encoded, got)
		}
	}
}

func TestValidDecodeWithDraftExamples(t *testing.T) {
	for _, entry := range validDraftExamples {
		got, err := Decode(entry.encoded)

		if err != nil {
			t.Errorf("Expected decoded string, got error \"%s\"", err)
		}

		if !bytes.Equal(got, entry.decoded) {
			t.Errorf("Unexpected decoding result for \"%s\", expected %v, got %v", entry.encoded, entry.decoded, got)
		}
	}
}

func TestInvalidInputLengthDecode(t *testing.T) {
	_, err := Decode([]byte("ABCD"))

	if err != InvalidLengthError {
		t.Errorf("Expected InvalidLengthError, got \"%s\"", err)
	}
}

func TestInvalidEncodedInputAlphabet(t *testing.T) {
	_, err := Decode([]byte("aa"))

	if err != InvalidEncodingCharactersError {
		t.Errorf("Expected InvalidEncodingCharactersError, got \"%s\"", err)
	}
}

func TestInvalidOverflow(t *testing.T) {
	/*
		[1] Chapter 6:

		Even though a Base45 encoded string contains only characters from the
		alphabet in Table 1 the following case has to be considered: The
		string "FGW" represents 65535 (FFFF in base 16), which is a valid
		encoding.  The string "GGW" would represent 65536 (10000 in base 16),
		which is represented by more than 16 bit.
	*/

	valid, err := Decode([]byte("FGW"))

	if err != nil {
		t.Errorf("Expected value, got error %s", err)
	}

	if !bytes.Equal(valid, []byte{255, 255}) {
		t.Errorf("Expected the valid decoded value to be 0xffff")
	}

	_, err = Decode([]byte("GGW"))

	if err != InvalidEncodedDataOverflowError {
		t.Errorf("Expected InvalidEncodedDataOverflowError, got \"%s\"", err)
	}
}

func TestEncodeUrlSafe(t *testing.T) {
	got := EncodeUrlSafe([]byte("Hello!!"))

	expected := "%2569%20VD92EX0"
	if !strings.EqualFold(got, expected) {
		t.Errorf("Expected \"%s\", got \"%s\"", expected, got)
	}
}

func TestDecodeUrlSafe(t *testing.T) {
	got, err := DecodeUrlSafe("%2569%20VD92EX0")

	if err != nil {
		t.Errorf("Expected url decoded data, got error \"%s\"", err)
	}

	expected := []byte("Hello!!")

	if !bytes.Equal(got, expected) {
		t.Errorf("Expected \"%s\", got \"%s\"", expected, got)
	}
}

func BenchmarkEncode(b *testing.B) {
	dec := make([]byte, b.N)
	rand.Read(dec)
	b.ResetTimer()

	Encode(dec)
}

func BenchmarkDecode(b *testing.B) {
	enc := bytes.Repeat([]byte("BB8"), b.N)
	b.ResetTimer()

	Decode(enc)
}
