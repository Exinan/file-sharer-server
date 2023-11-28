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

	server.HandleFunc("/download/{filename}", downloadHandler).Methods("GET")

	server.HandleFunc("/delete/{filename}", deleteHandler).Methods("DELETE")

	server.HandleFunc("/shutdown", shutdownHandler).Methods("POST")

	//start the server on port 8080
	http.Handle("/", server)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.WithFields(log.Fields{
			"main": "starting server",
		}).Info("Error starting server:", err)

		fmt.Println("Error starting server:", err)
	}
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

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	filePath := uploadDir + filename

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "downloadHandler",
		}).Info("Error opening file:", err)
		fmt.Println("Error opening file:", err)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// respond file
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)

	log.WithFields(log.Fields{
		"func": "downloadHandler",
	}).Info("File downloaded successfully")
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	filePath := uploadDir + filename

	// Удаляем файл
	err := os.Remove(filePath)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "deleteHandler",
		}).Info("Error deleting file:", err)

		fmt.Println("Error deleting file:", err)
		http.Error(w, "Error deleting file", http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{
		"func": "deleteHandler",
	}).Info("File deleted successfully")
	w.Write([]byte("File deleted successfully"))
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {

	log.WithFields(log.Fields{
		"func": "shutdownHandler",
	}).Info("Received shutdown request. Shutting down...")

	fmt.Println("Received shutdown request. Shutting down...")
	w.Write([]byte("Shutting down server..."))

	os.Exit(0)
}
