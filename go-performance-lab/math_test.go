package main

import "testing"

func TestAdd(t *testing.T) {
	got := Add(2, 4)
	want := 6
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestMultiply(t *testing.T) {
	got := Multiply(2, 4)
	want := 8
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func BenchmarkMultiply(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Multiply(2, 4)
	}
}
func BenchmarkSlowFunction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SlowFunction(1000000)
	}
}

func BenchmarkAllocateMemory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AllocateMemory()
	}

}
