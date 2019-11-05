package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"math"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
	"math/rand"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

// Block represents each 'item' in the blockchain
// Data represents the data each block holds - some number
// Hash is the hash of a block
// prevHash is the hash of the last block, allowing chain validation
// validator is the id of the node which won the POL round
type Block struct {
	Index     int
	Timestamp string
	Data       int
	Hash      string
	PrevHash  string
	Validator string
	Luck float64
}

// Blockchain is a series of Blocks created by a node
// accepted Blockchain is the network current chain
var acceptedBlockchain []Block
var Blockchain []Block


// announcements broadcasts winning validator to all nodes
// luck broadcasts the highest luck to all nodes
var announcements = make(chan string)
var luck = make(chan float64)

var mutex = &sync.Mutex{}

// validators keeps track of open validators and balances
var validators = make(map[string]int)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// create genesis block
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), 0, calculateBlockHash(genesisBlock), "", "", 0}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	httpPort := os.Getenv("PORT")

	// start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+httpPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("HTTP Server Listening on port :", httpPort)
	defer server.Close()

	go func() {
		for candidate := range candidateBlocks {
			mutex.Lock()
			tempBlocks = append(tempBlocks, candidate) //TODO: REWRITE
			mutex.Unlock()
		}
	}()


	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

// propogateChain is a protocol used when a node has generated a chain that is luckier than the current accepted chain
func propogateChain() {
	//calculate luck of the chain to be propogated
	newLuck := chainLuck(Blockchain)

	// let all the other nodes know
			for _ = range validators {
				announcements <- "\nLuckier chain minted by: " + address  + " Chain luck: " +  + "\n"
				luck <- newLuck
			}

}

func handleConn(conn net.Conn) {
	defer conn.Close()

	go func() {
		for {
			msg := <-announcements
			io.WriteString(conn, msg)
		}
	}()
	// validator address
	var address string

	// allow user to allocate number of tokens to stake
	// the greater the number of tokens, the greater chance to forging a new block
	io.WriteString(conn, "Enter token balance:")
	scanBalance := bufio.NewScanner(conn)
	for scanBalance.Scan() {
		balance, err := strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v not a number: %v", scanBalance.Text(), err)
			return
		}
		if balance <= 1000 {
			log.printf("not enough staked, you must stake more currency to become a validator")
			return
		}

		// validator is assigned address: a hash of the current time
		t := time.Now()
		address = calculateHash(t.String())
		validators[address] = balance
		fmt.Println(validators)
		break
	}


	io.WriteString(conn, "\nEnter some Data:")

	scanData := bufio.NewScanner(conn)

	//block proposal by nodes
	go func() {
		for {
			// take in Data from stdin and add it to blockchain after conducting necessary validation
			for scanData.Scan() {
				Data, err := strconv.Atoi(scanData.Text())
				// if malicious party tries to mutate the chain with a bad input, delete them as a validator and they lose their staked tokens
				if err != nil {
					log.Printf("%v not a number: %v", scanData.Text(), err)
					delete(validators, address)
					conn.Close()
				}

				mutex.Lock()
				oldLastIndex := Blockchain[len(Blockchain)-1]
				mutex.Unlock()

				// create newBlock for consideration to be forged
				newBlock, err := generateBlock(oldLastIndex, Data, address)
				if err != nil {
					log.Println(err)
					continue
				}
				if isBlockValid(newBlock, oldLastIndex) {
					candidateBlocks <- newBlock
				}
				io.WriteString(conn, "\nEnter a new Data:")
			}
		}
	}()

	// simulate receiving broadcast
	for {
		time.Sleep(time.Minute)
		mutex.Lock()
		output, err := json.Marshal(Blockchain)
		mutex.Unlock()
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(conn, string(output)+"\n")
	}

}

// isBlockValid makes sure block is valid by checking index
// and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateBlockHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// validates an entire chain (useful for participants)
func isChainValid(chain []Block) bool {
	for index, block := range chain {
		if index > 0 {
			if !isBlockValid(chain[index], chain[index-1]) {
				return false
			}
		}
	}
	return true
}

// calculates the luck of an entire chain (useful for participants)
func chainLuck(chain []Block) int {
	sum := 0
	for _, block := range chain {
		sum += block.luck
	}
	return sum
}

// SHA256 hashing
// calculateHash is a simple SHA256 hashing function
func calculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// calculateBlockHash returns the hash of all block information
// info hashed includes the block's index, timestamp, data, and the previous hash
func calculateBlockHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.Data) + block.PrevHash
	return calculateHash(record)
}

// generateBlock creates a new block using previous block's hash
func generateBlock(oldBlock Block, Data int, address string) (Block, error) {

	var newBlock Block

	// wait for length of roundtime
	time.Sleep(15 * time.Second)

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Data = Data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateBlockHash(newBlock)
	newBlock.Validator = address

	// luck generation with inversely proportional wait time
	randNum := rand.float64()
	newBlock.Luck = randNum
	waitTime := 10/(1 + math.Pow(e, (-9 * (randnum - 0.5))))
	time.Sleep(waitTime * time.Second)


	return newBlock, nil
}
