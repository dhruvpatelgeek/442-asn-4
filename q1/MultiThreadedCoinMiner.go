package main

import (
	"bytes"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
)

const NUM_THREADS = 8 // number of threads you want to run this on

const url string = "http://cpen442coin.ece.ubc.ca/verify_example_coin"
const hash_of_prev_coin string = "a9c1ae3f4fc29d0be9113a42090a5ef9fdef93f5ec4777a008873972e60bb532"
const student_number = "43586999"
const prefix_a = "CPEN 442 Coin"
const prefix_b = "2021"

var wg sync.WaitGroup

func NewSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func worker(id int, prefix_a string, prefix_b string, id_of_miner string) {
	fmt.Printf("Worker %d starting\n", id)
	token := make([]byte, 360)
	var coin_blob string
	var coin_blob_64 string
	for {
		rand.Read(token)                     // randomize byes
		coin_blob = fmt.Sprintf("%x", token) // string print to blob in hex

		//UNCOMMENT TO TEST WITH THE ANSWER
		//coin_blob = "19b69abba0b1bcbfb6584e253ced51c2e0fd454cc78578a82b98f1121e8e65f0ec20ef88"

		res := prefix_a + prefix_b + hash_of_prev_coin + coin_blob + id_of_miner
		sha256input := NewSHA256([]byte(res))       // take hash
		sha256Str := fmt.Sprintf("%x", sha256input) // string pprint hash

		if strings.HasPrefix(sha256Str, "00000000") { // if valid print and break
			coin_blob_64 = b64.StdEncoding.EncodeToString([]byte(coin_blob))
			print("\n ANS:" + coin_blob_64 + "|" + id_of_miner + "\n")
			break
		}
	}
	// send the data to the rest endpoint

	values := map[string]string{"coin_blob": coin_blob_64, "id_of_miner": id_of_miner}
	jsonValue, _ := json.Marshal(values)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	fmt.Println("URL:>", url)
	var id_of_miner string
	id_of_miner = fmt.Sprintf("%x", NewSHA256([]byte(student_number)))
	token := make([]byte, 36)
	rand.Read(token)

	for i := 0; i < NUM_THREADS; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			worker(i, prefix_a, prefix_b, id_of_miner)
		}()
	}
	wg.Wait()
}
