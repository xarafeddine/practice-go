package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func postHandler(w http.ResponseWriter, req *http.Request) {
	// Check if method is POST
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read and decode the request body
	var lake Lake
	err := json.NewDecoder(req.Body).Decode(&lake)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Store the lake in our database
	store[lake.Id] = lake

	// Send success response
	w.WriteHeader(http.StatusOK)
}

func deleteHandler(w http.ResponseWriter, req *http.Request) {
	// Check if method is DELETE
	if req.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the ID from query parameters
	// req.PathValue("id")
	id := req.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter required", http.StatusBadRequest)
		return
	}

	// Check if lake exists
	if _, exists := store[id]; !exists {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	// Delete the lake from store
	delete(store, id)

	// Send success response
	w.WriteHeader(http.StatusOK)
}

func getHandler(w http.ResponseWriter, req *http.Request) {
	// Check if method is GET
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the ID from query parameters
	id := req.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter required", http.StatusBadRequest)
		return
	}

	// Check if lake exists
	lake, exists := store[id]
	if !exists {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	// Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Encode and send the lake data
	err := json.NewEncoder(w).Encode(lake)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func http_server() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)
	stdout, err := os.Create(os.Getenv("OUTPUTPATH"))
	checkError(err)
	defer stdout.Close()
	writer := bufio.NewWriterSize(stdout, 16*1024*1024)
	actionsCount, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/delete", deleteHandler)
	go http.ListenAndServe(portSuffix, nil)
	time.Sleep(100 * time.Millisecond)
	var actions []string
	for i := 0; i < int(actionsCount); i++ {
		actionsItem := readLine(reader)
		actions = append(actions, actionsItem)
	}
	for _, actionStr := range actions {
		var action Action
		err := json.Unmarshal([]byte(actionStr), &action)
		checkError(err)
		switch action.Type {
		case "post":
			_, err := http.Post(address+"/post", "application/json", strings.NewReader(action.Payload))
			checkError(err)
		case "delete":
			client := &http.Client{}
			req, err := http.NewRequest("DELETE", address+"/delete?id="+action.Payload, nil)
			checkError(err)
			resp, err := client.Do(req)
			checkError(err)
			if resp.StatusCode != 200 {
				fmt.Fprintf(writer, "%s\n", resp.Status)
				continue
			}
		case "get":
			resp, err := http.Get(address + "/get?id=" + action.Payload)
			checkError(err)
			if resp.StatusCode != 200 {
				fmt.Fprintf(writer, "%s\n", resp.Status)
				continue
			}
			var lake Lake
			err = json.NewDecoder(resp.Body).Decode(&lake)
			checkError(err)
			fmt.Fprintf(writer, "%s\n", lake.Name)
			fmt.Fprintf(writer, "%d\n", lake.Area)
		}
	}
	fmt.Fprintf(writer, "\n")
	writer.Flush()
}

const portSuffix = ":3333"

var address = "http://127.0.0.1" + portSuffix

type Lake struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Area int32  `json:"area"`
}
type Action struct {
	Type    string
	Payload string
}

var store = map[string]Lake{}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}
	return strings.TrimRight(string(str), "\r\n")
}
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
