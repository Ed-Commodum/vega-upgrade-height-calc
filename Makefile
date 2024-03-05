build:
	go build -o ./bin/hcalc

run: build
	./bin/hcalc --rpc-addr http://164.92.138.136:26657