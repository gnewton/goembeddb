package main

import (
	"os"
)

type Db interface {
	Put([]byte, []byte) error
	Get([]byte) ([]byte, error)
	File() *os.File
	Close() error
}
