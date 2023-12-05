package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var hexString = "fc4881e4f0ffffffe8d0000000415141505251564831d265488b52603e488b52183e488b52203e488b72503e480fb74a4a4d31c94831c0ac3c617c022c2041c1c90d4101c1e2ed5241513e488b52203e8b423c4801d03e8b80880000004885c0746f4801d0503e8b48183e448b40204901d0e35c48ffc93e418b34884801d64d31c94831c0ac41c1c90d4101c138e075f13e4c034c24084539d175d6583e448b40244901d0663e418b0c483e448b401c4901d03e418b04884801d0415841585e595a41584159415a4883ec204152ffe05841595a3e488b12e949ffffff5d49c7c1000000003e488d95fe0000003e4c8d85060100004831c941ba45835607ffd54831c941baf0b5a256ffd557303020573030004d657373616765426f7800"

func main() {
	http.HandleFunc("/report", handleReport)
	http.ListenAndServe(":8080", nil)
}

func handleReport(w http.ResponseWriter, r *http.Request) {
	byteMap := make(map[int]string)
	for i := 0; i < len(hexString); i += 2 {
		byteMap[i/2] = hexString[i : i+2]
	}

	keys := make([]int, 0, len(byteMap))
	for k := range byteMap {
		keys = append(keys, k)
	}

	// Shuffle the keys
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })

	fmt.Printf("Shuffled keys: %v\n", keys)

	type KeyValue struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	var keyValuePairs []KeyValue
	for _, k := range keys {
		keyValuePairs = append(keyValuePairs, KeyValue{
			Key:   strconv.Itoa(k),
			Value: byteMap[k],
		})
	}

	response := struct {
		Report []KeyValue `json:"report"`
	}{
		Report: keyValuePairs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
