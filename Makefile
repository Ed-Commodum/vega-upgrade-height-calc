build:
	go build -o ./bin/hcalc

run: build
	./bin/hcalc
