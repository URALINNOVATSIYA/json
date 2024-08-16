package json

import (
	stdjson "encoding/json"
	"testing"
)

func BenchmarkDecodeString(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`"hello world"`)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := Bytes(data).Decode(); err != nil {
				b.Fatal("Decode error:", err)
			}
		}
	})
}

func BenchmarkStdJsonUnmarshalString(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`"hello world"`)
	b.RunParallel(func(pb *testing.PB) {
		var s string
		for pb.Next() {
			if err := stdjson.Unmarshal(data, &s); err != nil {
				b.Fatal("Unmarshal error:", err)
			}
		}
	})
}

func BenchmarkDecodeInt(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`1234`)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := Bytes(data).Decode(); err != nil {
				b.Fatal("Decode error:", err)
			}
		}
	})
}

func BenchmarkStdJsonUnmarshalInt(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`1234`)
	b.RunParallel(func(pb *testing.PB) {
		var s int
		for pb.Next() {
			if err := stdjson.Unmarshal(data, &s); err != nil {
				b.Fatal("Unmarshal error:", err)
			}
		}
	})
}

func BenchmarkDecodeFloat(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`1.234`)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := Bytes(data).Decode(); err != nil {
				b.Fatal("Decode error:", err)
			}
		}
	})
}

func BenchmarkStdJsonUnmarshalFloat(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`1.234`)
	b.RunParallel(func(pb *testing.PB) {
		var s float64
		for pb.Next() {
			if err := stdjson.Unmarshal(data, &s); err != nil {
				b.Fatal("Unmarshal error:", err)
			}
		}
	})
}

func BenchmarkDecodeFloatE(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`1234E+2`)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := Bytes(data).Decode(); err != nil {
				b.Fatal("Decode error:", err)
			}
		}
	})
}

func BenchmarkStdJsonUnmarshalFloatE(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`1234E+2`)
	b.RunParallel(func(pb *testing.PB) {
		var s float64
		for pb.Next() {
			if err := stdjson.Unmarshal(data, &s); err != nil {
				b.Fatal("Unmarshal error:", err)
			}
		}
	})
}

func BenchmarkDecodeBool(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`true`)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := Bytes(data).Decode(); err != nil {
				b.Fatal("Decode error:", err)
			}
		}
	})
}

func BenchmarkStdJsonUnmarshalBool(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`true`)
	b.RunParallel(func(pb *testing.PB) {
		var s bool
		for pb.Next() {
			if err := stdjson.Unmarshal(data, &s); err != nil {
				b.Fatal("Unmarshal error:", err)
			}
		}
	})
}

func BenchmarkDecodeNull(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`null`)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := Bytes(data).Decode(); err != nil {
				b.Fatal("Decode error:", err)
			}
		}
	})
}

func BenchmarkStdJsonUnmarshalNull(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`null`)
	b.RunParallel(func(pb *testing.PB) {
		var s any
		for pb.Next() {
			if err := stdjson.Unmarshal(data, &s); err != nil {
				b.Fatal("Unmarshal error:", err)
			}
		}
	})
}

func BenchmarkDecodeStringList(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`["hello world", "\u042d\u043c\u043e\u0434\u0437\u0438: \u2705, \ud83e\udd2a", "\ud83e\udd2a", "\u2705"]`)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := Bytes(data).Decode(); err != nil {
				b.Fatal("Decode error:", err)
			}
		}
	})
}

func BenchmarkStdJsonUnmarshalStringList(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`["hello world", "\u042d\u043c\u043e\u0434\u0437\u0438: \u2705, \ud83e\udd2a", "\ud83e\udd2a", "\u2705"]`)
	b.RunParallel(func(pb *testing.PB) {
		var s []string
		for pb.Next() {
			if err := stdjson.Unmarshal(data, &s); err != nil {
				b.Fatal("Unmarshal error:", err)
			}
		}
	})
}

func BenchmarkDecodeObjectWithFloat(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`{"executionTime":0.1423170566558838,"requestTime":0.1424870491027832,"activationTime":1719388801.086876,"backendReceived":1719388801.087046,"backendHandled":1719388801.229363,"memoryUsage":4194304.1}`)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := Bytes(data).Decode(); err != nil {
				b.Fatal("Decode error:", err)
			}
		}
	})
}

func BenchmarkStdJsonUnmarshalObjectWithFloat(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`{"executionTime":0.1423170566558838,"requestTime":0.1424870491027832,"activationTime":1719388801.086876,"backendReceived":1719388801.087046,"backendHandled":1719388801.229363,"memoryUsage":4194304.1}`)
	b.RunParallel(func(pb *testing.PB) {
		var s map[string]float64
		for pb.Next() {
			if err := stdjson.Unmarshal(data, &s); err != nil {
				b.Fatal("Unmarshal error:", err)
			}
		}
	})
}

func BenchmarkDecodeObjectWithStrings(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`{"executionTime":"0.1423170566558838","requestTime":"0.1424870491027832","activationTime":"1719388801.086876","backendReceived":"1719388801.087046","backendHandled":"1719388801.229363","memoryUsage":"4194304.1"}`)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := Bytes(data).Decode(); err != nil {
				b.Fatal("Decode error:", err)
			}
		}
	})
}

func BenchmarkStdJsonUnmarshalObjectWithStrings(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`{"executionTime":"0.1423170566558838","requestTime":"0.1424870491027832","activationTime":"1719388801.086876","backendReceived":"1719388801.087046","backendHandled":"1719388801.229363","memoryUsage":"4194304.1"}`)
	b.RunParallel(func(pb *testing.PB) {
		var s map[string]string
		for pb.Next() {
			if err := stdjson.Unmarshal(data, &s); err != nil {
				b.Fatal("Unmarshal error:", err)
			}
		}
	})
}
