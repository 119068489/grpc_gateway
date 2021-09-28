package main

import "testing"

func BenchmarkMain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Ranges()
	}
}

//19825	     60118 ns/op	     792 B/op	      33 allocs/op
// BenchmarkMain-8   	288256086	         4.308 ns/op	       0 B/op	       0 allocs/op
