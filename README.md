# goembeddb
TLDR; a MWE embedding a key-value store file of 20M k/v pairs (1.9GB)
into a Go binary and using the db from the binary.

The same Go program that writes the db then is recompiled to embed the db.

Then the same Go program is used to read the db.



## How does it work?
### Constant database
The [constanct DB](https://en.wikipedia.org/wiki/Cdb_(software)) or cdb, is a write-once, static key-value store that stores its content in a [single file](http://cr.yp.to/cdb.html).

### cdb: Go implementation
A Go implementation of cdb is used: [github.com/alldroll/cdb](https://github.com/alldroll/cdb)

Checkout the repo and type:
```
 $ make clean; make
```

This does the following:

1. Creates an empty cdb database file `test.cdb` (using `touch`).
2. Go build compiles `main.go`, which embeds the empty `test.cdb` resulting in the Go binary `goembeddb` which is 2.1MB in size.
3. `goembeddb` is run, which creates a real cdb database, with 20m k/v pairs, resulting in a 1.9BG `test.cdb` file.
4. Go build re-compiles `main.go`, which embeds the 1.9GB `test.cdb` resulting in the Go binary `goembeddb` which is now 1.9GB in size.
5. `goembeddb` is run, in sequential read-mode, reading all 20m values by keys, in key write order.
```
Runtime 14.268457937s
Max resident=1939484 Elapsed real=0:14.30 PageFaults=44988 KCPU=0.57 UCPU=14.98 Elapsed=14.30
```
6. `goembeddb` is run, in random read-mode, reading randomly 20m values by keys.
```
Runtime 20.711046526s
Max resident=1939336 Elapsed real=0:20.74 PageFaults=52041 KCPU=0.60 UCPU=21.78 Elapsed=20.74
```
7. `goembeddb` is run, in random single record read-mode, reading 1 value by a random key.
```
Runtime 115.46µs
Max resident=2248 Elapsed real=0:00.00 PageFaults=213 KCPU=0.00 UCPU=0.00 Elapsed=0.00
```
   


# Example run

```
$make; make clean
mkdir .tmp
mkdir: cannot create directory ‘.tmp’: File exists
make: [Makefile:3: default] Error 1 (ignored)
export TMPDIR=./.tmp
touch test.cdb

// Compiling with empty embedded db
/usr/bin/time -f "Max resident=%M Elapsed real=%E PageFaults=%R KCPU=%S UCPU=%U Elapsed=%e" go build
Max resident=51128 Elapsed real=0:00.18 PageFaults=14759 KCPU=0.04 UCPU=0.26 Elapsed=0.18
-rwxrwxr-x 1 gnewton gnewton 2.1M Sep 23 18:55 goembeddb

// Writing db with goembeddb
/usr/bin/time -f "Max resident=%M Elapsed real=%E PageFaults=%R KCPU=%S UCPU=%U Elapsed=%e" ./goembeddb
2021/09/23 18:55:38 main.go:83: Opening db (write): test.cdb
2021/09/23 18:55:38 main.go:153: 0: key=k_0  value=[Lorem ipsum dolor sit amet, consectetur adipiscing elit k_0]
2021/09/23 18:55:38 main.go:153: 1: key=k_1  value=[Lorem ipsum dolor sit amet, consectetur adipiscing elit k_1]
2021/09/23 18:55:38 main.go:153: 2: key=k_2  value=[Lorem ipsum dolor sit amet, consectetur adipiscing elit k_2]
2021/09/23 18:55:38 main.go:153: 3: key=k_3  value=[Lorem ipsum dolor sit amet, consectetur adipiscing elit k_3]
2021/09/23 18:55:38 main.go:153: 4: key=k_4  value=[Lorem ipsum dolor sit amet, consectetur adipiscing elit k_4]
2021/09/23 18:55:53 main.go:153: 19999996: key=k_19999996  value=[Lorem ipsum dolor sit amet, consectetur adipiscing elit k_19999996]
2021/09/23 18:55:53 main.go:153: 19999997: key=k_19999997  value=[Lorem ipsum dolor sit amet, consectetur adipiscing elit k_19999997]
2021/09/23 18:55:53 main.go:153: 19999998: key=k_19999998  value=[Lorem ipsum dolor sit amet, consectetur adipiscing elit k_19999998]
2021/09/23 18:55:53 main.go:153: 19999999: key=k_19999999  value=[Lorem ipsum dolor sit amet, consectetur adipiscing elit k_19999999]
2021/09/23 18:55:53 main.go:90: Closing db
2021/09/23 18:57:09 main.go:94:   Successfully wrote 20000000 key/value pairs
Runtime 1m30.64634428s
Max resident=383608 Elapsed real=1:30.67 PageFaults=102501 KCPU=52.49 UCPU=38.79 Elapsed=90.67

// DB size
-rw-rw-r-- 1 gnewton gnewton 1.9G Sep 23 18:57 test.cdb

// Compiling with populated embedded db
/usr/bin/time -f "Max resident=%M Elapsed real=%E PageFaults=%R KCPU=%S UCPU=%U Elapsed=%e" go build
Max resident=3912788 Elapsed real=0:21.78 PageFaults=528400 KCPU=1.68 UCPU=19.47 Elapsed=21.78

// Go binary size
-rwxrwxr-x 1 gnewton gnewton 1.9G Sep 23 18:57 goembeddb

// Reading from embedded db with goembeddb, sequential, all records
/usr/bin/time -f "Max resident=%M Elapsed real=%E PageFaults=%R KCPU=%S UCPU=%U Elapsed=%e" ./goembeddb -l -q
Runtime 13.946731996s
Max resident=1939076 Elapsed real=0:13.98 PageFaults=43720 KCPU=0.41 UCPU=14.78 Elapsed=13.98

// Reading from embedded db with goembeddb, random, all records
/usr/bin/time -f "Max resident=%M Elapsed real=%E PageFaults=%R KCPU=%S UCPU=%U Elapsed=%e" ./goembeddb -R -q
Runtime 20.67538902s
Max resident=1939820 Elapsed real=0:20.71 PageFaults=51214 KCPU=0.74 UCPU=21.59 Elapsed=20.71

// Reading from embedded db with goembeddb, random, one record
/usr/bin/time -f "Max resident=%M Elapsed real=%E PageFaults=%R KCPU=%S UCPU=%U Elapsed=%e" ./goembeddb -S -q
Runtime 133.166µs
Max resident=2244 Elapsed real=0:00.00 PageFaults=213 KCPU=0.00 UCPU=0.00 Elapsed=0.00
$
```

# Observations



# Discussion on Google Groups golang-nuts
https://groups.google.com/g/golang-nuts/c/jFKGLbTv2XQ
