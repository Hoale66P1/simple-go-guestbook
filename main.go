package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type GuestbookEntry struct {
	Name      string
	Message   string
	CreatedAt time.Time
}

type PageData struct {
	Entries []GuestbookEntry
}

var entries []GuestbookEntry

const dataFile = "data.json"

func saveEntries() {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		log.Println("Error saving data:", err)
		return
	}
	os.WriteFile(dataFile, data, 0644)
}

func loadEntries() {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Println("Error reading data:", err)
		return
	}
	json.Unmarshal(data, &entries)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		name := r.FormValue("name")
		message := r.FormValue("message")

		if name != "" && message != "" {
			entry := GuestbookEntry{
				Name:      name,
				Message:   message,
				CreatedAt: time.Now(),
			}

			entries = append([]GuestbookEntry{entry}, entries...)

			saveEntries()
		}
	}

	data := PageData{
		Entries: entries,
	}

	tmpl.Execute(w, data)
}

func main() {
	loadEntries()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)

	log.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
