package main

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkSingleRandomRead(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	db, err := NewDbAlldrolReaderFromBytes(dbBytes)
	if err != nil {
		b.Error(err)
	}
	for i := 0; i < b.N; i++ {
		randomRead(1, db, defaultDbSize)
	}
}
