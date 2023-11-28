package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const uploadDir = "./uploads/"
const logFileName = "log.log"

func main() {

	log.SetFormatter(&log.JSONFormatter{})

	// open log file
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		// Set the file as the output for logging
		log.SetOutput(file)
		defer file.Close()
	} else {
		log.Fatal("Failed to open the log file.", err)
	}
	//create uploads folder
	err = os.MkdirAll(uploadDir, os.ModePerm) // os.ModePerm - provide acess to read, write and execution
	if err != nil {
		log.WithFields(log.Fields{
			"main": "create uploads folder",
		}).Info("Error creating upload directory:", err)

		fmt.Println("Error creating upload directory:", err)
		return
	}

	server := mux.NewRouter()

	// adding hendling methods

	server.HandleFunc("/ping", pingHandler).Methods("GET")

	server.HandleFunc("/upload", uploadHandler).Methods("POST")

	//start the server on port 8080
	http.Handle("/", server)
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {

	log.WithFields(log.Fields{
		"func": "pingHandler",
	}).Info("Received ping request.")

	fmt.Println("Received ping request.")
	w.Write([]byte("Pong! Server is up and running."))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) //  file size to up to 10 MB (10 * 2^20 = 10 * 1 MB)

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.WithFields(log.Fields{
			"func": "uploadHandler",
		}).Info("Error getting file:", err)

		fmt.Println("Error getting file:", err)
		http.Error(w, "Error getting file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a file in the uploads directory
	f, err := os.Create(uploadDir + handler.Filename)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "uploadHandler",
		}).Info("Error creating file:", err)

		fmt.Println("Error creating file:", err)
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copy the file
	_, err = io.Copy(f, file)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "uploadHandler",
		}).Info("Error copying file:", err)

		fmt.Println("Error copying file:", err)
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{
		"func": "uploadHandler",
	}).Info("File uploaded successfully")
	w.Write([]byte("File uploaded successfully"))
}
