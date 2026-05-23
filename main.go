package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Railway needs the PORT variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create uploads directory if it doesn't exist
	os.MkdirAll("uploads", 0755)

	// Homepage logic
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
		<html><body style="font-family: sans-serif; padding: 20px;">
		<h1>Game Host</h1>
		<form action="/upload" method="post" enctype="multipart/form-data">
			<input type="file" name="file" style="padding: 10px;">
			<input type="submit" value="Upload Game" style="padding: 10px; cursor: pointer;">
		</form>
		<hr><h3>Available Files:</h3><ul>`)

		files, _ := filepath.Glob("uploads/*")
		for _, f := range files {
			name := filepath.Base(f)
			fmt.Fprintf(w, `<li><a href="/uploads/%s">%s</a></li>`, name, name)
		}
		fmt.Fprintf(w, `</ul></body></html>`)
	})

	// Optimized Upload logic
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Upload failed", http.StatusBadRequest)
			return
		}
		defer file.Close()

		out, err := os.Create(filepath.Join("uploads", header.Filename))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		// io.Copy is extremely fast; it streams the data directly to disk
		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// File server
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	log.Println("Server running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}