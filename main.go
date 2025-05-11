package main

import (
	"fmt"
	"net/http"
	"path"
	"text/template"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error occured", r)
		} else {
			fmt.Println("Application running perfectly")
		}
	}()

	type M map[string]any

	handleApple := func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello Apple"))
		if err != nil {
			fmt.Println(err)
		}
	}

	// string response
	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/index", handlerIndex)
	http.HandleFunc("/hello", handlerHello)
	http.HandleFunc("/apple", handleApple)

	// render an HTML
	http.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
		filePath := path.Join("views", "index.html")
		fmt.Println(filePath)
		tmpl, err := template.ParseFiles(filePath)
		if err != nil {
			fmt.Println(err.Error())
		}
		data := map[string]any{
			"title": "Basic Web",
			"name":  "Putu Arya",
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// parse all globals
	tmpl, err := template.ParseGlob("views/partial/*")
	if err != nil {
		panic(err.Error())
	}

	// render a partial HTML Template
	http.HandleFunc("/welcome", func(w http.ResponseWriter, r *http.Request) {
		data := M{
			"name": "Putu Arya in Welcome",
		}

		err = tmpl.ExecuteTemplate(w, "index", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		data := M{
			"name": "Putu Arya in About",
		}

		err = tmpl.ExecuteTemplate(w, "about", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// file access handler
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	const address = ":9000"
	server := new(http.Server)
	server.Addr = address
	fmt.Println("server running at: ", address)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	message := "Welcome"
	_, err := w.Write([]byte(message))
	if err != nil {

	}
}

func handlerHello(w http.ResponseWriter, r *http.Request) {
	message := "Hello World"
	_, err := w.Write([]byte(message))
	if err != nil {

	}
}
