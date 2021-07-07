package base45

import (
	"math/rand"
	"testing"
)

func benchmarkEncode(len int, b *testing.B) {
	dec := make([]byte, len)
	rand.Read(dec)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Encode(dec)
	}
}

func BenchmarkEncode1(b *testing.B) {
	benchmarkEncode(1, b)
}

func BenchmarkEncode128(b *testing.B) {
	benchmarkEncode(128, b)
}

func BenchmarkEncode512(b *testing.B) {
	benchmarkEncode(512, b)
}

func BenchmarkEncode1024(b *testing.B) {
	benchmarkEncode(1024, b)
}

func BenchmarkEncode8192(b *testing.B) {
	benchmarkEncode(8192, b)
}

func benchmarkDecode(len int, b *testing.B) {
	dec := make([]byte, len)
	rand.Read(dec)
	enc := Encode(dec)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Decode(enc)
	}
}

func BenchmarkDecode1(b *testing.B) {
	benchmarkDecode(1, b)
}

func BenchmarkDecode128(b *testing.B) {
	benchmarkDecode(128, b)
}

func BenchmarkDecode512(b *testing.B) {
	benchmarkDecode(512, b)
}

func BenchmarkDecode1024(b *testing.B) {
	benchmarkDecode(1024, b)
}

func BenchmarkDecode8192(b *testing.B) {
	benchmarkDecode(8192, b)
}
