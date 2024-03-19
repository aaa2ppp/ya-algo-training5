.PHONY: profclear

mem.out:
	go test -run XXX -bench 100_200 -benchmem -memprofile=mem.out -memprofilerate=1 -o mem.bin .

cpu.out:
	go test -run XXX -bench 100_200 -cpuprofile=cpu.out -o cpu.bin .

profmem: mem.out
	go tool pprof mem.bin mem.out

profcpu: cpu.out
	go tool pprof cpu.bin cpu.out

profclear:
	rm mem.out mem.bin cpu.out cpu.bin
