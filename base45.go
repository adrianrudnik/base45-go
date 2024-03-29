// Package base45 implements encoding and decoding of base 45 data by
// https://datatracker.ietf.org/doc/rfc9285/
package base45

import (
	"bytes"
	"encoding/binary"
	"math"
	"net/url"
)

/*
	Chapter references:

	[1] https://datatracker.ietf.org/doc/rfc9285/
        2022-08-11 rfc9285
*/

/*
	[1] Chapter 4:

    A 45-character subset of US-ASCII is used; the 45 characters usable
    in a QR code in Alphanumeric mode (see Section 7.3.4 and Table 2 of
    ISO18004).  Base45 encodes 2 bytes in 3 characters, compared to
    Base64, which encodes 3 bytes in 4 characters.

	[1] Chapter 4.2:

    The Alphanumeric mode is defined to use 45 characters as specified in
    this alphabet.
*/

// Alphabet defines the 45 usable characters for the base 45 encoding.
var Alphabet = []byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B',
	'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
	'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	' ', '$', '%', '*', '+', '-', '.', '/', ':',
}

// encodeSingleByte takes in a byte and converts it to base 45.
func encodeSingleByte(in byte) []byte {
	/*
		[1] Chapter 4:

		For encoding a single byte [a], it MUST be interpreted as a base 256
		number, i.e. as an unsigned integer over 8 bits.  That integer MUST
		be converted to base 45 [c d] so that a = c + (45*d).  The values c
		and d are then looked up in Table 1 to produce a two character
		string.
	*/
	a := int(in)
	c := Alphabet[a%45]
	d := Alphabet[a/45%45]

	return []byte{c, d}
}

// encodeTwoBytes takes two bytes and converts it to base 45.
func encodeTwoBytes(in []byte) []byte {
	/*
		[1] Chapter 4:

		For encoding, two bytes [a, b] MUST be interpreted as a number n in
		base 256, i.e. as an unsigned integer over 16 bits so that the number
		n = (a * 256) + b.
	*/
	n := binary.BigEndian.Uint16(in)

	/*
		[1] Chapter 4:

		This number n is converted to base 45 [c, d, e] so that n = c + (d *
		45) + (e * 45 * 45).  Note the order of c, d and e which are chosen
		so that the left-most [c] is the least significant.

		The values c, d, and e are then looked up in Table 1 to produce a
		three character string.  The process is reversed when decoding.
	*/
	c := Alphabet[n%45]
	d := Alphabet[n/45%45]
	e := Alphabet[n/(45*45)%45]

	return []byte{c, d, e}
}

// Encode encodes the given byte to base 45.
// If an empty input is given, an empty result will be returned.
func Encode(in []byte) []byte {
	// Instead of analysing the possible output length, we
	// create a byte array with the estimated capacity of two
	// output bytes per one input byte, which is a bit more
	// than we need, but it keeps the code clean.
	out := make([]byte, 0, len(in)*2)

	// Next up we consume chunks up to two bytes of decoded date
	// and encode it to base 45.
	buf := make([]byte, 2)
	reader := bytes.NewReader(in)

	for {
		n, _ := reader.Read(buf)

		if n == 2 {
			out = append(out, encodeTwoBytes(buf)...)
		} else if n == 1 {
			out = append(out, encodeSingleByte(buf[0])...)
		} else {
			// on EOF or error
			break
		}
	}

	return out
}

// EncodeURLSafe encodes the given bytes to a query safe string.
// If an empty input is given, an empty result will be returned.
func EncodeURLSafe(in []byte) string {
	/*
		[1] Chapter 6:

		It should be noted that the resulting string after encoding to Base45
		might include non-URL-safe characters so if the URL including the
		Base45 encoded data has to be URL-safe, one has to use percent-
		encoding.
	*/
	parts := &url.URL{Path: string(Encode(in))}

	return parts.String()
}

// decodeTwoBytes decodes two base 45 encoded bytes to one decoded byte.
// This will be used for very short or trailing base 45 encoded data.
func decodeTwoBytes(dst, src []byte) error {
	/*
		[1] Chapter 4:

		For encoding a single byte [a], it MUST be interpreted as a base 256
		number, i.e. as an unsigned integer over 8 bits.  That integer MUST
		be converted to base 45 [c d] so that a = c + (45 * d).  The values c
		and d are then looked up in Table 1 to produce a two-character
		string.

		For decoding a Base45 encoded string the inverse operations are
		performed.
	*/
	c := bytes.IndexByte(Alphabet, src[0])
	d := bytes.IndexByte(Alphabet, src[1])

	val := c + (d * 45)

	// Detect possible overflow attack
	if val > math.MaxUint8 {
		return ErrInvalidEncodedDataOverflow
	}

	copy(dst, []byte{byte(val)})

	return nil
}

// decodeThreeBytes decodes three base 45 encoded bytes to two decoded bytes.
func decodeThreeBytes(dst, src []byte) error {
	/*
		[1] Chapter 4:

		For encoding, two bytes [a, b] MUST be interpreted as a number n in
		base 256, i.e. as an unsigned integer over 16 bits so that the number
		n = (a * 256) + b.

		This number n is converted to base 45 [c, d, e] so that n = c + (d *
		45) + (e * 45 * 45).  Note the order of c, d and e which are chosen
		so that the left-most [c] is the least significant.

		The values c, d, and e are then looked up in Table 1 to produce a
		three character string.  The process is reversed when decoding.

		For decoding a Base45 encoded string the inverse operations are
		performed.
	*/

	// We skip checks if c, d, e return -1 as the exposed Decode function
	// already does an alphabet check and only allowed entries pass through here.
	c := bytes.IndexByte(Alphabet, src[0])
	d := bytes.IndexByte(Alphabet, src[1])
	e := bytes.IndexByte(Alphabet, src[2])

	val := c + (d * 45) + (e * 45 * 45)

	/*
		[1] Chapter 6:

		When implementing encoding and decoding it is important to be very
		careful so that buffer overflow or similar issues do not occur.  This
		of course includes the calculations in base 45 and lookup in the
		table of characters (Table 1).  A decoder must also be robust
		regarding input, including proper handling of any octet value 0-255,
		including the NUL character (ASCII 0).
	*/
	if val > math.MaxUint16 {
		return ErrInvalidEncodedDataOverflow
	}

	binary.BigEndian.PutUint16(dst, uint16(val))

	return nil
}

// Decode reads the base 45 encoded bytes and returns the decoded bytes.
// If an empty input is given, ErrEmptyInput is returned.
func Decode(in []byte) ([]byte, error) {
	// Calls to this function expect an input, empty calls should not happen.
	if len(in) == 0 {
		return nil, ErrEmptyInput
	}

	/*
		[1] Chapter 6:

		Implementations MUST reject any input that is not a valid encoding.
		For example, it MUST reject the input (encoded data) if it contains
		characters outside the base alphabet (in Table 1) when interpreting
		base-encoded data.
	*/
	for _, v := range in {
		if !bytes.Contains(Alphabet, []byte{v}) {
			return nil, ErrInvalidEncodingCharacters
		}
	}

	/*
		[1] Chapter 4:

		A byte string [a b c d ... x y z] with arbitrary content and
		arbitrary length MUST be encoded as follows: From left to right pairs
		of bytes MUST be encoded as described above.  If the number of bytes
		is even, then the encoded form is a string with a length that is
		evenly divisible by 3.  If the number of bytes is odd, then the last
		(rightmost) byte MUST be encoded on two characters as described
		above.

		For decoding a Base45 encoded string the inverse operations are
		performed.
	*/
	if len(in)%3 != 0 && (len(in)+1)%3 != 0 {
		return nil, ErrInvalidLength
	}

	// Instead of analysing the possible output length, we allocate
	// enough capacity to keep the code clean and readable. In this case
	// the expected output length will always be smaller than the input length.
	out := make([]byte, len(in))

	buf := make([]byte, 3)
	reader := bytes.NewReader(in)
	written := 0

	for {
		read, _ := reader.Read(buf)

		if read == 3 {
			// Three bytes go in, two come out, we copy them into the output slice
			err := decodeThreeBytes(out[written:written+2], buf)

			if err != nil {
				return nil, err
			}

			written += 2
		} else if read == 2 {
			// Two bytes go in, one comes out, we copy it into the output slice
			err := decodeTwoBytes(out[written:written+1], buf[0:2])

			if err != nil {
				return nil, err
			}

			written += 1
		} else {
			// this happens on EOF or error, as read == 0 in both cases
			break
		}
	}

	return out[:written], nil
}

// DecodeURLSafe reads the given url encoded base 45 encoded data and returns the decoded bytes.
// If an empty input is given, ErrEmptyInput is returned.
func DecodeURLSafe(in string) ([]byte, error) {
	/*
		[1] Chapter 6:

		It should be noted that the resulting string after encoding to Base45
		might include non-URL-safe characters so if the URL including the
		Base45 encoded data has to be URL-safe, one has to use percent-
		encoding.
	*/
	enc, err := url.QueryUnescape(in)

	if err != nil {
		return nil, ErrInvalidURLSafeEscaping
	}

	dec, err := Decode([]byte(enc))

	if err != nil {
		return nil, err
	}

	return dec, nil
}
