package main

import (
	"bytes"

	alldrol "github.com/alldroll/cdb"
	"os"
)

type DbAlldrol struct {
	handle *alldrol.CDB
	writer alldrol.Writer
	reader alldrol.Reader
	file   *os.File
}

func NewDbAlldrolReaderFromBytes(dbb []byte) (*DbAlldrol, error) {
	db := new(DbAlldrol)
	db.handle = alldrol.New()
	var err error
	db.reader, err = db.handle.GetReader(bytes.NewReader(dbBytes))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewDbAlldrolReader(filename string) (*DbAlldrol, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	db := new(DbAlldrol)
	db.file = f
	db.handle = alldrol.New()
	db.reader, err = db.handle.GetReader(db.file)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewDbAlldrolWriter(filename string) (*DbAlldrol, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	db := new(DbAlldrol)
	db.file = f
	db.handle = alldrol.New()
	db.writer, err = db.handle.GetWriter(db.file)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DbAlldrol) Put(key []byte, value []byte) error {
	return db.writer.Put(key, value)
}

func (db *DbAlldrol) Get(key []byte) ([]byte, error) {
	return db.reader.Get(key)
}

func (db *DbAlldrol) Close() error {
	err := db.writer.Close()
	if err != nil {
		return err
	}
	err = db.file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (db *DbAlldrol) File() *os.File {
	return db.file
}
