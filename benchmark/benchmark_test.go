package benchmark

import (
	"runtime"
	"testing"
)

type User struct {
	ID   int
	Name string
}

func init() {
	// Limit Go to use only 1 CPU
	runtime.GOMAXPROCS(1)
}

// Using []User
func BenchmarkSliceOfStructs(b *testing.B) {
	users := make([]User, 1000)
	for i := 0; i < b.N; i++ {
		_ = users
	}
}

// Using []*User
func BenchmarkSliceOfPointers(b *testing.B) {
	users := make([]*User, 1000)
	for i := 0; i < 1000; i++ {
		users[i] = &User{ID: i, Name: "User"}
	}
	for i := 0; i < b.N; i++ {
		_ = users
	}
}
