# Upgrade Height Calculator

Presented here is a simple calculator that makes use of cometBFT RPCs to estimate the block height of a network upgrade based on the desired time of the upgrade and the historical rate of block production.
<br/><br/>
## Basic Usage

Install or update Go by following the [official install doc](https://go.dev/doc/install).

Clone the repo and build the binary.
```
git clone https://github.com/ed-commodum/vega-upgrade-height-calc
cd vega-upgrade-height-calc
make build
```

By default the tool will look for a cometBFT RPC at localhost:26657, so you can run it locally on the host of your Vega node like so:
```
./bin/hcalc
```

If you are not running it locally you will need to specify a remote cometBFT RPC address. The address in the below command is publicly avilable at the time of writing, it is also included in the Makefile so you can simply issue `make run` to run it.
```
./bin/hcalc --rpc-addr http://164.92.138.136:26657
```

You can specify a desired time in minutes until the upgrade like so:
```
./bin/hcalc --minutes-until-upgrade 100
```

Alternatively you can specify a date an time in the following format `YYYY-MM-DD HH:MM:SS`:
```
./bin/hcalc --upgrade-date-time "2024-05-20 11:00:00"
```
