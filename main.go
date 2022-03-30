package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

	database["1"] = Student{
		ID:   1,
		Name: "Budi",
		Age:  5,
	}

	// initial server - handler
	http.HandleFunc("/", func(wr http.ResponseWriter, rq *http.Request) {
		pesan := []byte(`{"pesan":"server dijalankan"}`)
		SettingJSONWR(wr, pesan, http.StatusOK)
	})
	// Get Request Student
	http.HandleFunc("/get-students", func(wr http.ResponseWriter, rq *http.Request) {
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
	// POST Student (Soal 1) = localhost:8080/post-students
	http.HandleFunc("/post-students", func(wr http.ResponseWriter, rq *http.Request) {
		if rq.Method != "POST" {
			pesan := []byte(`{"pesan":"invalid http method"}`)
			SettingJSONWR(wr, pesan, http.StatusMethodNotAllowed)
			return
		}
		var student Student
		Payload := rq.Body
		defer rq.Body.Close()
		err := json.NewDecoder(Payload).Decode(&student)
		if err != nil {
			pesan := []byte(`{"pesan":"error ketika parsing data"}`)
			SettingJSONWR(wr, pesan, http.StatusInternalServerError)
			return
		}
		studentID := strconv.Itoa(student.ID)
		database[studentID] = student

		pesan := []byte(`{"pesan":"success create data student"}`)
		SettingJSONWR(wr, pesan, http.StatusCreated)
	})
	// Update Data Soal 2 = localhost:8080/put-student?id=2
	http.HandleFunc("/put-students", func(wr http.ResponseWriter, rq *http.Request) {
		if rq.Method != "PUT" {
			pesan := []byte(`{"pesan":"invalid error http method"}`)
			SettingJSONWR(wr, pesan, http.StatusMethodNotAllowed)
			return
		}
		id := rq.URL.Query()["id"][0]
		student, ok := database[id]
		if !ok {
			pesan := []byte(`{"pesan":"data student tak ditemukan"}`)
			SettingJSONWR(wr, pesan, http.StatusOK)
			return
		}

		var newstudent Student

		payload := rq.Body

		defer rq.Body.Close()

		err := json.NewDecoder(payload).Decode(&newstudent)
		if err != nil {
			pesan := []byte(`{"pesan":"error ketika parsing data"}`)
			SettingJSONWR(wr, pesan, http.StatusInternalServerError)
			return
		}

		student.Name = newstudent.Name
		student.Age = newstudent.Age

		studentID := strconv.Itoa(student.ID)
		database[studentID] = student

		studentJSON, err := json.Marshal(&student)
		if err != nil {
			pesan := []byte(`{"pesan":"error ketika parsing data"}`)
			SettingJSONWR(wr, pesan, http.StatusInternalServerError)
			return
		}
		SettingJSONWR(wr, studentJSON, http.StatusOK)

	})

	// Get Data Soal 3 = localhost:8080/student?id=2
	http.HandleFunc("/get-student", func(wr http.ResponseWriter, rq *http.Request) {
		if rq.Method != "GET" {
			pesan := []byte(`{"pesan":"invalid error http method"}`)
			SettingJSONWR(wr, pesan, http.StatusMethodNotAllowed)
		}

		if _, ok := rq.URL.Query()["id"]; !ok {
			pesan := []byte(`{"pesan":"membutuhkan id student"}`)
			SettingJSONWR(wr, pesan, http.StatusBadRequest)
			return
		}

		id := rq.URL.Query()["id"][0]
		student, ok := database[id]
		if !ok {
			pesan := []byte(`{"pesan":"data student tak ditemukan"}`)
			SettingJSONWR(wr, pesan, http.StatusOK)
			return
		}

		studentJSON, err := json.Marshal(&student)
		if err != nil {
			pesan := []byte(`{"pesan":"error ketika parsing data"}`)
			SettingJSONWR(wr, pesan, http.StatusInternalServerError)
			return
		}
		SettingJSONWR(wr, studentJSON, http.StatusOK)
	})

	// Delete Data Soal 4 localhost:8080/delete-student?id=2
	http.HandleFunc("/delete-student", func(wr http.ResponseWriter, rq *http.Request) {
		if rq.Method != "DELETE" {
			pesan := []byte(`{"pesan":"invalid error http method"}`)
			SettingJSONWR(wr, pesan, http.StatusMethodNotAllowed)

		}

		if _, ok := rq.URL.Query()["id"]; !ok {
			pesan := []byte(`{"pesan":"membutuhkan id student"}`)
			SettingJSONWR(wr, pesan, http.StatusBadRequest)
			return
		}

		id := rq.URL.Query()["id"][0]
		student, ok := database[id]
		if !ok {
			pesan := []byte(`{"pesan":"data student tak ditemukan"}`)
			SettingJSONWR(wr, pesan, http.StatusOK)
			return
		}
		delete(database, id)

		studentJSON, err := json.Marshal(&student)
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
