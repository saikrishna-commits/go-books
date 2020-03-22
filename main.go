package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Book struct {
	Id       int    `json:"id"`
	BookName string `json:"bookName"`
	Author   string `json:"author"`
}

var books []Book

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "get called"}`))
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post called"}`))
}

func put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "put called"}`))
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "delete called"}`))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to HomePage!")
	fmt.Println("Endpoint Hit: HomePage")
}

func createNewBook(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var newBook Book
	err := decoder.Decode(&newBook)

	if err != nil {
		panic(err)
	}
	books = append(books, newBook)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newBook)

	fmt.Println("Endpoint Hit: Creating New Booking")

}
func getAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Println("Endpoint Hit: returnAllBooks")
	json.NewEncoder(w).Encode(books)
}

func returnSingleBookById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, book := range books {
		// string to int
		s, err := strconv.Atoi(key)
		if err == nil {
			if book.Id == s {
				json.NewEncoder(w).Encode(book)
			}
		}
	}
}

func deleteBooks(w http.ResponseWriter, r *http.Request) {
	bookId := mux.Vars(r)["id"]

	bookIdInt, _ := strconv.Atoi(bookId)

	for i, book := range books {
		if book.Id == bookIdInt {
			books = append(books[:i], books[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", bookIdInt)
		}
	}
}

// this function is used before main
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	csvfile, err := os.Open("books.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	fileReader := csv.NewReader(csvfile)

	for {
		// Read each record from csv
		record, err := fileReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		id, err := strconv.Atoi(record[0])
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}

		books = append(books, Book{Id: id, BookName: record[1], Author: record[2]})
	}

}

func main() {

	// Get the GITHUB_USERNAME environment variable
	githubUsername, exists := os.LookupEnv("GITHUB_USERNAME")
	if exists == true {
		fmt.Println(githubUsername)
	}
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/newBook", createNewBook).Methods("POST")
	myRouter.HandleFunc("/books", getAllBooks).Methods("GET")
	myRouter.HandleFunc("/books/{id}", returnSingleBookById)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
