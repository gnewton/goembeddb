package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"
)

//go:embed test.cdb
var dbBytes []byte

var readSequentialFlag bool
var readSingleRandomRecordFlag bool
var readRandomFlag bool
var quietFlag bool
var numRecords int

const defaultDbSize = 20000000

//const defaultDbSize = 200000
const firstLast = 5
const dbName = "test.cdb"

var start time.Time

func init() {
	start = time.Now()
	flag.BoolVar(&quietFlag, "q", false, "quiet")
	flag.BoolVar(&readRandomFlag, "R", false, "Read random")
	flag.BoolVar(&readSequentialFlag, "l", false, "Read key sequential")
	flag.BoolVar(&readSingleRandomRecordFlag, "S", false, "Read key sequential")
	flag.IntVar(&numRecords, "n", defaultDbSize, "Num records to add to db")
	flag.Parse()
	if quietFlag {
		log.SetOutput(ioutil.Discard)
	}
	rand.Seed(time.Now().UnixNano())
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if readSequentialFlag || readRandomFlag || readSingleRandomRecordFlag {
		read(numRecords, dbBytes, readSequentialFlag, readRandomFlag, readSingleRandomRecordFlag)
	} else {
		write(dbName, numRecords)
	}

	elapsed := time.Since(start)
	fmt.Printf("Runtime %s\n", elapsed)
}

func read(num int, dbBytes []byte, seq, rand, randSingle bool) {
	log.Println("Opening db (read): embedded file")
	db, err := NewDbAlldrolReaderFromBytes(dbBytes)
	if err != nil {
		log.Fatal(err)
	}
	if randSingle {
		randomRead(1, db, num)
		log.Printf("  Successfully random read %d key/value pairs", 1)
	} else {
		if seq {
			log.Printf("Start seqential read, all")
			writeOrReadSeq(num, db, false)
			log.Printf("  Successfully sequential read %d key/value pairs", num)
		} else {
			// Read num records, random keys
			log.Printf("Start random read, %d records", num)
			randomRead(num, db, num)
			log.Printf("  Successfully random read %d key/value pairs", num)
		}
	}
}

func write(name string, num int) {
	log.Println("Opening db (write): " + name)
	db, err := NewDbAlldrolWriter(name)
	if err != nil {
		log.Fatal(err)
	}

	writeOrReadSeq(num, db, true)
	log.Println("Closing db")
	if err = db.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("  Successfully wrote %d key/value pairs", num)
}

func makeKeyValue(i int) (string, string) {
	return "k_" + strconv.Itoa(i), "Lorem ipsum dolor sit amet, consectetur adipiscing elit k_" + strconv.Itoa(i)
}

func makeRandomKeyValue(n int) (string, string) {
	return makeKeyValue(rand.Intn(n))
}

func randomRead(n int, db Db, keyRange int) {

	for i := 0; i < n; i++ {
		key, value := makeRandomKeyValue(keyRange)
		thisValue, err := db.Get([]byte(key))
		if err != nil {
			log.Fatal(err)
		}
		if thisValue == nil {
			log.Fatal(errors.New("No value for key:" + key))
		}
		if string(thisValue) != value {
			log.Fatal(errors.New("Wrong value: " + string(thisValue) + " != " + value))
		}

		// Print first and last firstLast k/v
		printFirstLast(key, value, i, n, firstLast)
	}
}

func writeOrReadSeq(n int, db Db, write bool) {
	for i := 0; i < n; i++ {
		key, value := makeKeyValue(i)

		if write { // Put
			if err := db.Put([]byte(key), []byte(value)); err != nil {
				fileInfo, err := db.File().Stat()
				log.Println(i, "Error Filesize=", fileInfo.Size()/1024/1024, "GB  ", fileInfo.Size())
				log.Fatal(err)
			}
		} else { // Get
			thisValue, err := db.Get([]byte(key))
			if err != nil {
				log.Fatal(err)
			}
			if thisValue == nil {
				log.Fatal(errors.New("No value for key:" + key))
			}
			if string(thisValue) != value {
				log.Fatal(errors.New("Wrong value: " + string(thisValue) + " != " + value))
			}
		}
		printFirstLast(key, value, i, n, firstLast)
	}
}

func printFirstLast(key, value string, i, n, firstLast int) {
	if i < firstLast || i > n-firstLast {
		log.Printf("%d: key=%s  value=[%s]", i, key, value)
	}
}
