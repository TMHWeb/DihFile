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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	os.MkdirAll("uploads", 0755)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
		<html><body style="font-family: sans-serif; padding: 20px;">
		<h1>Game Host</h1>
		<form action="/upload" method="post" enctype="multipart/form-data">
			<input type="file" name="file" style="padding: 10px;">
			<input type="submit" value="Upload" style="padding: 10px; cursor: pointer;">
		</form>
		<hr><h3>Available Files:</h3><ul>`)

		files, _ := filepath.Glob("uploads/*")
		for _, f := range files {
			name := filepath.Base(f)
			fmt.Fprintf(w, `<li><a href="/uploads/%s">%s</a></li>`, name, name)
		}
		fmt.Fprintf(w, `</ul></body></html>`)
	})

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

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "failed to save file", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	log.Println("server running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
