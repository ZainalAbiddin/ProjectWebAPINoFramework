/*

Field Data Type
ID (PK) Numeric
Name String
Age Integer

Soal 1
Functions required:
Function Register Student
Endpoint POST http://{host}:{port}/student
Description Register student into the system
Payload { “id”: 1, name: “budi”, age: 5 }

Soal 2
Function Update Student
Endpoint PUT http://{host}:{port}/student/{id}
Description Update student by ID
Payload { “name: “budi kurniawan” }

Soal 3
Function Get Student
Endpoint GET http://{host}:{port}/student/{id}
Description Get student by ID

Soal 4
Function Delete Student
Endpoint DELETE http://{host}:{port}/student/{id}
Description Delete student by ID
*/

/*
kekurangan file ini
1. Belum menemukan ID menjadi primary key tanpa menggunakan tag gorm dan fitur mysql
*/

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var (
	database = make(map[string]Student)
)

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
	// Create Fake Database/ Get Student
	http.HandleFunc("/students", GetStudents)
	// Soal 1(post), Soal 2(put), Soal 3(get by id), Soal 4(delete)
	http.HandleFunc("/student/", methodStudents)

	// cek error dan rute alamat
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func methodStudents(wr http.ResponseWriter, rq *http.Request) {
	switch rq.Method {
	case "GET":
		GetStudentsByID(wr, rq)
		return
	case "POST":
		PostStudents(wr, rq)
		return
	case "PUT":
		PutStudents(wr, rq)
		return
	case "DELETE":
		DeleteStudent(wr, rq)
		return
	default:
		pesan := []byte(`{"pesan":"invalid http method"}`)
		SettingJSONWR(wr, pesan, http.StatusMethodNotAllowed)
		return
	}
}

func SettingJSONWR(wr http.ResponseWriter, pesan []byte, httpCode int) {
	wr.Header().Set("Content-type", "application/json")
	wr.WriteHeader(httpCode)
	wr.Write(pesan)
}

func GetStudents(wr http.ResponseWriter, rq *http.Request) {
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
}

func PostStudents(wr http.ResponseWriter, rq *http.Request) {
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

	studentJSON, err := json.Marshal(&student)
	if err != nil {
		pesan := []byte(`{"pesan":"error ketika parsing data"}`)
		SettingJSONWR(wr, pesan, http.StatusInternalServerError)
		return
	}
	SettingJSONWR(wr, studentJSON, http.StatusOK)
}

func PutStudents(wr http.ResponseWriter, rq *http.Request) {
	// metode params
	params := strings.Split(rq.URL.String(), "/")
	if len(params) != 3 {
		pesan := []byte(`{"pesan":"error penulisan id di URL"}`)
		SettingJSONWR(wr, pesan, http.StatusOK)
		return
	}
	id := params[2]
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
	pesan := []byte(`{"pesan":"update data berhasil"}`)
	SettingJSONWR(wr, pesan, http.StatusOK)
}

func GetStudentsByID(wr http.ResponseWriter, rq *http.Request) {
	// metode params
	params := strings.Split(rq.URL.String(), "/")
	if len(params) != 3 {
		pesan := []byte(`{"pesan":"error penulisan id di URL"}`)
		SettingJSONWR(wr, pesan, http.StatusOK)
		return
	}
	id := params[2]
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
}

func DeleteStudent(wr http.ResponseWriter, rq *http.Request) {
	// metode params
	params := strings.Split(rq.URL.String(), "/")
	if len(params) != 3 {
		pesan := []byte(`{"pesan":"error penulisan id di URL"}`)
		SettingJSONWR(wr, pesan, http.StatusOK)
		return
	}
	id := params[2]
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
	pesan := []byte(`{"pesan":"delete data berhasil"}`)
	SettingJSONWR(wr, pesan, http.StatusOK)
}

/*
1. You have 3 days to do the test.
2. You are NOT allowed to use any framework.
3. You are NOT allowed to use any DBMS.
4. Your code will be evaluated for both functional and non-functional such as code quality.
5. It’s nice if you have test your code, and if you do please include your test files in your project.
6. We are really interested in clean and maintanable codebase, so please solve the problem keeping this in mind.
*/

/*
	// metode query
	// if _, ok := rq.URL.Query()["id"]; !ok {
	// 	pesan := []byte(`{"pesan":"membutuhkan id student"}`)
	// 	SettingJSONWR(wr, pesan, http.StatusBadRequest)
	// 	return
	// }

	// metode params
	params := strings.Split(rq.URL.String(), "/")
	if len(params) != 3 {
		pesan := []byte(`{"pesan":"error penulisan id di URL"}`)
		SettingJSONWR(wr, pesan, http.StatusOK)
		return
	}
	id := params[2]
	student, ok := database[id]
	if !ok {
		pesan := []byte(`{"pesan":"data student tak ditemukan"}`)
		SettingJSONWR(wr, pesan, http.StatusOK)
		return
	}
*/
