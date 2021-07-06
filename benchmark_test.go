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

func BenchmarkEncodeByte1(b *testing.B) {
	benchmarkEncode(1, b)
}

func BenchmarkEncodeByte8(b *testing.B) {
	benchmarkEncode(8, b)
}

func BenchmarkEncodeByte64(b *testing.B) {
	benchmarkEncode(64, b)
}

func BenchmarkEncodeByte512(b *testing.B) {
	benchmarkEncode(512, b)
}

func BenchmarkEncodeByte1024(b *testing.B) {
	benchmarkEncode(1024, b)
}

func BenchmarkEncodeByte8192(b *testing.B) {
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

func BenchmarkDecodeChunk1(b *testing.B) {
	benchmarkDecode(1, b)
}

func BenchmarkDecodeChunk8(b *testing.B) {
	benchmarkDecode(8, b)
}

func BenchmarkDecodeChunk64(b *testing.B) {
	benchmarkDecode(64, b)
}

func BenchmarkDecodeChunk512(b *testing.B) {
	benchmarkDecode(512, b)
}

func BenchmarkDecodeChunk1024(b *testing.B) {
	benchmarkDecode(1024, b)
}

func BenchmarkDecodeChunk8192(b *testing.B) {
	benchmarkDecode(8192, b)
}
