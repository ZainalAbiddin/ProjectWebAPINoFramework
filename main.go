package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var (
	database = make(map[string]Student)
)

func SettingJSONWR(wr http.ResponseWriter, pesan []byte, httpCode int) {
	wr.Header().Set("Content-type", "application/json")
	wr.WriteHeader(httpCode)
	wr.Write(pesan)
}
func main() {

	database["01"] = Student{
		ID:   1,
		Name: "Budi",
		Age:  5,
	}

	// initial server - handler
	http.HandleFunc("/", func(wr http.ResponseWriter, rq *http.Request) {
		pesan := []byte(`{"pesan":"server dijalankan"}`)
		SettingJSONWR(wr, pesan, http.StatusOK)
	})
	http.HandleFunc("/student", func(wr http.ResponseWriter, rq *http.Request) {

		if rq.Method != "GET" {
			pesan := []byte(`{"pesan":"invalid http method"}`)
			SettingJSONWR(wr, pesan, http.StatusMethodNotAllowed)
			return
		}

		var students []Student

		for _, student := range database {
			students = append(students, student)
		}

		studentJSON, err := json.Marshal(&students)
		if err != nil {
			pesan := []byte(`{"pesan":"error ketika parsing data"}`)
			SettingJSONWR(wr, pesan, http.StatusInternalServerError)
			return
		}
		SettingJSONWR(wr, studentJSON, http.StatusOK)
	})

	// cek error dan rute alamat
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
1. You have 3 days to do the test.
2. You are NOT allowed to use any framework.
3. You are NOT allowed to use any DBMS.
4. Your code will be evaluated for both functional and non-functional such as code quality.
5. It’s nice if you have test your code, and if you do please include your test files in your project.
6. We are really interested in clean and maintanable codebase, so please solve the problem keeping this in mind.

Field Data Type
ID (PK) Numeric
Name String
Age Integer

Functions required:
Function Register Student
Endpoint POST http://{host}:{port}/student
Description Register student into the system
Payload { “id”: 1, name: “budi”, age: 5 }

Function Update Student
Endpoint PUT http://{host}:{port}/student/{id}
Description Update student by ID
Payload { “name: “budi kurniawan” }

Function Get Student
Endpoint GET http://{host}:{port}/student/{id}
Description Get student by ID

Function Delete Student
Endpoint DELETE http://{host}:{port}/student/{id}
Description Delete student by ID
*/
