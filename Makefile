
default:
	-mkdir .tmp
	export TMPDIR=./.tmp
	touch test.cdb
	@echo ""
	@echo "// Compiling with empty embedded db"
	/usr/bin/time -f "Max resident=%M Elapsed real=%E PageFaults=%R KCPU=%S UCPU=%U Elapsed=%e" go build
	@ls -lh goembeddb
	@echo ""
	@echo "// Writing db with goembeddb"
	/usr/bin/time -f "Max resident=%M Elapsed real=%E PageFaults=%R KCPU=%S UCPU=%U Elapsed=%e" ./goembeddb
	@echo ""
	@echo "// DB size"
	@ls -lh test.cdb
	@echo ""
	@echo "// Compiling with populated embedded db"
	/usr/bin/time -f "Max resident=%M Elapsed real=%E PageFaults=%R KCPU=%S UCPU=%U Elapsed=%e" go build
	@echo ""
	@echo "// Go binary size"
	@ls -lh goembeddb
	@echo ""
	@echo "// Reading from embedded db with goembeddb, sequential, all records"
	/usr/bin/time -f "Max resident=%M Elapsed real=%E PageFaults=%R KCPU=%S UCPU=%U Elapsed=%e" ./goembeddb -l -q
	@echo ""
	@echo "// Reading from embedded db with goembeddb, random, all records"
	/usr/bin/time -f "Max resident=%M Elapsed real=%E PageFaults=%R KCPU=%S UCPU=%U Elapsed=%e" ./goembeddb -R -q
	@echo ""
	@echo "// Reading from embedded db with goembeddb, random, one record"
	/usr/bin/time -f "Max resident=%M Elapsed real=%E PageFaults=%R KCPU=%S UCPU=%U Elapsed=%e" ./goembeddb -S -q
	@echo ""

clean:
	-rm goembeddb test.cdb *~
