package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Book struct {
	Title string
	Pages uint16
}

var books = []Book{{"One Of Us", 100}, {"Last of Us", 80}, {"How to Go?", 777}}

func getBooks(w http.ResponseWriter, r *http.Request){
	for _, book := range books{
		fmt.Fprintf(w, "Book: %s. Pages: %d.\n", book.Title, book.Pages)
	}
}

func getBook(w http.ResponseWriter, r *http.Request){
	title := r.URL.Query().Get("title")

	var b Book

	for _, book := range books{
		if book.Title == title{
			b = book
			fmt.Fprintf(w, "You found book %s.\n This book have %d pages\n", b.Title, b.Pages)
			return
		}
	}

	http.Error(w, "Book not found!", http.StatusNotFound)
}

func createBook(w http.ResponseWriter, r *http.Request){
	var b Book

	err:= json.NewDecoder(r.Body).Decode(&b)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	books = append(books, b)
	
	fmt.Fprintf(w, "Added book %+v", b)
}

func deleteBook(w http.ResponseWriter, r *http.Request){
	title := r.URL.Query().Get("title")
	bl := len(books)

	result := []Book{}
	
	for _, book := range books{
		if book.Title != title{
			result = append(result, book)
		}
	}
	
	books = result

	if bl != len(books){
		fmt.Fprintf(w, "Deleted book %s", title)
		return
	}

	http.Error(w, "Nothing deleted", http.StatusBadRequest)
}

func main(){
	r := mux.NewRouter()

	r.HandleFunc("/book", getBook).Methods("GET")
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books", deleteBook).Methods("DELETE")

	http.ListenAndServe(":80", r)
}