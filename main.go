package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	defaultRpcAddr         string = "localhost:26657"
	defaultUpgradeDateTime string = ""
	defaultMinutesUntil    uint   = 120
	defaultBlockWindow     uint   = 10000
)

var (
	rpcAddr         string
	upgradeDateTime string
	minutesUntil    uint
	blockWindow     uint
)

func init() {
	flag.StringVar(&rpcAddr, "rpc-addr", defaultRpcAddr, "Vega cometBFT RPC address.")
	flag.StringVar(&upgradeDateTime, "upgrade-date-time", defaultUpgradeDateTime, "The desired dateTime of the upgrade in the format: '2006-01-02 15:04:05'")
	flag.UintVar(&minutesUntil, "minutes-until-upgrade", defaultMinutesUntil, "The desired number of minutes until the upgrade.")
	flag.UintVar(&blockWindow, "block-window", defaultBlockWindow, "the desired number of block over which to calculate.")
	flag.Parse()
}

func main() {
	validateRpcAddr()
	var upgradeTime time.Time
	var err error
	if upgradeDateTime != "" {
		upgradeTime, err = time.Parse(time.DateTime, upgradeDateTime)
		if err != nil {
			log.Fatalf("failed to parse upgrade dateTime: %s, must be of format '2006-01-02 15:04:05'", err)
		}
	} else {
		log.Printf("Upgrade dataTime not provided calculating block height for %d minutes from now\n", minutesUntil)
		upgradeTime = time.Now().Add(time.Minute * time.Duration(minutesUntil))
	}

	client := http.Client{Timeout: 5 * time.Second}

	recentHeight, blockRate := getBlockRate(client, blockWindow)

	log.Printf("Block rate over past %d blocks: %f blocks/s\n", blockWindow, blockRate)

	secondsUntilUpgrade := upgradeTime.Unix() - time.Now().Unix()

	log.Printf("Seconds until upgrade: %d\n", secondsUntilUpgrade)

	blocksUntilUpgrade := int64(blockRate * float64(secondsUntilUpgrade))

	log.Printf("Blocks until upgrade: %d\n", blocksUntilUpgrade)

	upgradeHeight := recentHeight + blocksUntilUpgrade

	log.Printf("Estimated upgrade height: %d\n", upgradeHeight)
}

func getBlockRate(client http.Client, window uint) (int64, float64) {
	// Get most recent block
	bodyBytes := getBlock(client, -1)

	recentBlockHeight, recentBlockTime := parseResponse(bodyBytes)

	// Get historical block
	bodyBytes = getBlock(client, recentBlockHeight-int64(window))

	blockHeight, blockTime := parseResponse(bodyBytes)

	return recentBlockHeight, float64(recentBlockHeight-blockHeight) / float64(recentBlockTime.Unix()-blockTime.Unix())
}

func getBlock(client http.Client, h int64) []byte {
	var res *http.Response
	var err error
	if h == -1 {
		res, err = client.Get(fmt.Sprintf("%s/%s", rpcAddr, "block"))
		if err != nil {
			log.Fatalf("error calling cometBFT rpc: %s. If no rpc is running locally then provide one with the rpc-addr flag", err)
		}
	} else {
		res, err = client.Get(fmt.Sprintf("%s/%s?height=%d", rpcAddr, "block", h))
		if err != nil {
			log.Fatalf("error calling cometBFT rpc: %s. If no rpc is running locally then provide one with the rpc-addr flag", err)
		}
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("could not read response body: %s", err)
	}

	return bodyBytes
}

func parseResponse(b []byte) (int64, time.Time) {
	decoded := map[string]interface{}{}

	err := json.Unmarshal(b, &decoded)
	if err != nil {
		log.Fatalf("failed to unmarshal response body: %s", err)
	}

	result, ok := decoded["result"].(map[string]interface{})
	if !ok {
		log.Fatal("Error parsing JSON, chosen RPC likely has truncated block store. Use a different RPC or use a smaller value for block-window.")
	}

	block, ok := result["block"].(map[string]interface{})
	if !ok {
		log.Fatal("Error parsing JSON. Use a different RPC address.")
	}

	header, ok := block["header"].(map[string]interface{})
	if !ok {
		log.Fatal("Error parsing JSON. Use a different RPC address.")
	}

	blockHeight, ok := header["height"].(string)
	if !ok {
		log.Fatal("Error parsing JSON. Use a different RPC address.")
	}

	blockTime, ok := header["time"].(string)
	if !ok {
		log.Fatal("Error parsing JSON. Use a different RPC address.")
	}

	t, err := time.Parse(time.RFC3339Nano, blockTime)
	if err != nil {
		log.Fatalf("failed to parse block time (%s): %s", blockTime, err)
	}

	bh, err := strconv.Atoi(blockHeight)
	if err != nil {
		log.Fatalf("failed to convert block height to int: %s", err)
	}

	return int64(bh), t
}

func validateRpcAddr() {
	if !strings.HasPrefix(rpcAddr, "http://") && !strings.HasPrefix(rpcAddr, "https://") {
		rpcAddr = "http://" + rpcAddr
	}
}
