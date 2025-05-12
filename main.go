package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

type Result struct {
	Name string
}

type ResultFile struct {
	Alias string
}

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	defer func() {
		if r := recover(); r != nil {
			logger.Error(r.(string))
		}
	}()

	const port = ":9000"

	// render view
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("views/index.html")
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "index", nil)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// process form
	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			tmpl, err := template.ParseFiles("views/result.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			name := r.FormValue("name")
			data := Result{Name: name}
			logger.Info("Name variable", "name", data)
			if err := tmpl.ExecuteTemplate(w, "result", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		logger.Error("Bad Request")
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {
			if err := r.ParseMultipartForm(1024); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// get alias input field
			alias := r.FormValue("alias")
			// get file input field
			file, header, err := r.FormFile("file")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			// get current working directory
			dir, err := os.Getwd()
			logger.Info("Working dir", "dir", dir)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// get original filename
			filename := header.Filename
			logger.Info("Original filename:", "filename", filename)
			if alias != "" {
				// take the alias
				// take the extension file ex: .jpg, .png, .mov, etc
				// it would be <alias>.jpg
				filename = fmt.Sprintf("%s%s", alias, filepath.Ext(filename))
			}

			// create target new file uploaded
			// create new memory for empty files (create space)
			emptySpace := filepath.Join(dir, "files", filename)
			targetFileSpace, err := os.OpenFile(emptySpace, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer targetFileSpace.Close()

			if _, err := io.Copy(targetFileSpace, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tmpl, err := template.ParseFiles("views/done_upload.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := tmpl.ExecuteTemplate(w, "done_upload", nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		logger.Error("Bad Request")
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	// test
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			x := 10
			_, err := w.Write(fmt.Appendf(nil, "You got this %d", x))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "POST":
			_, err := w.Write([]byte("You posted this"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
	})

	server := new(http.Server)
	server.Addr = port
	logger.Info("Running Server at: ", "port", port)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
