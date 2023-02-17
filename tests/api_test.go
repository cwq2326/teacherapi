package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"govtech/pkg/controllers"
	"govtech/pkg/models/request"
	"govtech/pkg/server/databases"
	"govtech/pkg/server/handlers"
	"govtech/pkg/utilities/messages"
)

var dsn string

func init() {
	err := godotenv.Load(filepath.Join("..", ".env"))
	if err != nil {
		fmt.Println("Failed to load .env file")
	}

	config := database.MySqlConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Port:     os.Getenv("DB_PORT"),
		Host:     os.Getenv("DB_HOST"),
		Name:     os.Getenv("DB_TEST_NAME"),
	}
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.User, config.Password,
		config.Host, config.Port, config.Name)
}

// Run all tests.
func TestEndPoints(t *testing.T) {
	t.Run("suspend endpoint", Suspend)
	t.Run("commonstudents endpoint", CommonStudents)
	t.Run("retrievefornotifications endpoint", RetrieveForNotification)
	t.Run("register endpoint", Register)
}

// Tests for "/api/suspend" endpoint.
func Suspend(t *testing.T) {
	// Init DB.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	database.InitTestDB(db)
	_, err = db.Exec(`INSERT INTO students VALUES ("test@gmail.com", 0)`)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Init router and middlewares.
	r := gin.Default()
	handlers.RegisterMiddlewares(r, db)
	r.POST("/api/suspend", controllers.Suspend)

	// Test for POST.

	// Positive cases.

	// Test with valid request body.
	// Should return http status 204.
	payload := request.SuspendRequest{
		Student: "test@gmail.com",
	}

	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/suspend", bytes.NewBuffer(jsonValue))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Test DB.
	var suspended int
	err = db.QueryRow(`SELECT suspended FROM students WHERE email = ?`, "test@gmail.com").Scan(&suspended)
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, 1, suspended)

	// Negative test cases.

	// Invalid student query param.
	// Should return status code 400 and error message.
	payload = request.SuspendRequest{
		Student: "",
	}

	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", "/api/suspend", bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"message":{"`+messages.MESSAGE_MISSING_PARAMS+`":{"student":"required"}}}`, rr.Body.String())

	// Clean up DB.
	database.CleanupTestDB(db)
	db.Close()
}

// Tests for "/api/commonstudents" endpoint.
func CommonStudents(t *testing.T) {
	// Init DB.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	database.InitTestDB(db)

	_, err = db.Exec(`INSERT INTO students VALUES ("student1@gmail.com", 0)`)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.Exec(`INSERT INTO students VALUES ("student2@gmail.com", 0)`)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.Exec(`INSERT INTO students VALUES ("common@gmail.com", 0)`)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.Exec(`INSERT INTO teachers VALUES ("teacher1@gmail.com")`)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.Exec(`INSERT INTO teachers VALUES ("teacher2@gmail.com")`)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.Exec(`INSERT INTO teaches VALUES ("teacher1@gmail.com", "student1@gmail.com")`)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.Exec(`INSERT INTO teaches VALUES ("teacher2@gmail.com", "student2@gmail.com")`)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.Exec(`INSERT INTO teaches VALUES ("teacher1@gmail.com", "common@gmail.com")`)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.Exec(`INSERT INTO teaches VALUES ("teacher2@gmail.com", "common@gmail.com")`)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Init router and middleware.
	r := gin.Default()
	handlers.RegisterMiddlewares(r, db)
	r.GET("/api/commonstudents", controllers.CommonStudents)

	// Test for GET.

	// Positive cases.

	// Test for teacher1@gmail.com.
	// Should return status code 200 and students: common@gmail.com and student1@gmail.com.
	req, _ := http.NewRequest("GET", `/api/commonstudents?teacher=teacher1%40gmail.com`, nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"students":["common@gmail.com","student1@gmail.com"]}`, rr.Body.String())

	// Test for teacher1@gmail.com and teacher2@gmail.com.
	// Should return status code 200 and students: common@gmail.com.
	req, _ = http.NewRequest("GET", `/api/commonstudents?teacher=teacher1%40gmail.com&teacher=teacher2%40gmail.com`, nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"students":["common@gmail.com"]}`, rr.Body.String())

	// Test for non existent teacher.
	// Should return status code 200 and students: null.
	req, _ = http.NewRequest("GET", `/api/commonstudents?teacher=non@existent.email`, nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"students":null}`, rr.Body.String())

	// Test for wrong parameter type.
	// Should return status code 200 and students: null.
	req, _ = http.NewRequest("GET", `/api/commonstudents?teacher=123`, nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"students":null}`, rr.Body.String())

	// Negative test cases.

	// Test for missing query parameters.
	// Should return http status code 400 and error message.
	req, _ = http.NewRequest("GET", `/api/commonstudents?`, nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"message":"`+messages.MissingQueryParamsMessage([]string{"teacher"})+"\"}", rr.Body.String())

	// Clean up DB.
	database.CleanupTestDB(db)
	db.Close()
}

// Tests for "/api/retrievefornotifications" endpoint.
func RetrieveForNotification(t *testing.T) {
	// Init DB.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	database.InitTestDB(db)

	_, err = db.Exec(`INSERT INTO teachers VALUES ("teacher@gmail.com")`)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.Exec(`INSERT INTO students VALUES ("nottagged@gmail.com", 0)`)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.Exec(`INSERT INTO students VALUES ("tagged1@gmail.com", 0)`)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.Exec(`INSERT INTO students VALUES ("tagged2@gmail.com", 0)`)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.Exec(`INSERT INTO students VALUES ("ishouldnotappear@gmail.com", 1)`)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.Exec(`INSERT INTO teaches VALUES ("teacher@gmail.com", "nottagged@gmail.com")`)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.Exec(`INSERT INTO teaches VALUES ("teacher@gmail.com", "ishouldnotappear@gmail.com")`)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Init router and middleware.
	r := gin.Default()
	handlers.RegisterMiddlewares(r, db)
	r.POST("/api/retrievefornotifications", controllers.RetrieveForNotifications)

	// Test for POST.

	// Positive cases.

	// Test for valid request body with no tagged students.
	// Should get status code 200 and "nottagged@gmail.com"
	payload := request.ReceieveForNotificationsRequest{
		Teacher:      "teacher@gmail.com",
		Notification: "hello world",
	}
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", `/api/retrievefornotifications`, bytes.NewBuffer(jsonValue))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"recipient":["nottagged@gmail.com"]}`, rr.Body.String())

	// Test for valid request body with tagged students.
	// Should get status code 200 and "nottagged@gmail.com"
	payload = request.ReceieveForNotificationsRequest{
		Teacher:      "teacher@gmail.com",
		Notification: "hello world @tagged1@gmail.com @tagged2@gmail.com @tagged1@gmail.com",
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/retrievefornotifications`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"recipient":["nottagged@gmail.com","tagged1@gmail.com","tagged2@gmail.com"]}`, rr.Body.String())

	// Negative cases.

	// Test for wrong teacher field format.
	// Should get status code 400 and error message.
	payload = request.ReceieveForNotificationsRequest{
		Teacher:      "teacher@gmailcom",
		Notification: "hello world @tagged1@gmail.com @tagged2@gmail.com",
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/retrievefornotifications`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Test for too long teacher field > 60.
	// Should get status code 400 and error message.
	payload = request.ReceieveForNotificationsRequest{
		Teacher:      "teacherteacherteacherteacherteacherteacherteacherteacherteacher@gmail.com",
		Notification: "hello world @tagged1@gmail.com @tagged2@gmail.com",
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/retrievefornotifications`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"message":{"`+messages.MESSAGE_MISSING_PARAMS+`":{"teacher":"max=60"}}}`, rr.Body.String())

	// Test for too long notification field > 200.
	// Should get status code 400 and error message.
	payload = request.ReceieveForNotificationsRequest{
		Teacher: "teacher@gmail.com",
		Notification: `ehelloworldhelloworldhelloworldhelloworldhelloworldhelloworld
						helloworldhellowhelloworldhelloworldhelloworldhelloworldhello
						worldhelloworldhelloworldhelloworldhelloworldhelloworldhellowo
						rldhelloworldvworldhelloworldhelloworld @tagged1@gmail.com @tagged2@gmail.com"`,
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/retrievefornotifications`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"message":{"`+messages.MESSAGE_MISSING_PARAMS+`":{"notification":"max=200"}}}`, rr.Body.String())

	// Test for missing teacher field.
	// Should get status code 400 and error message.
	payload = request.ReceieveForNotificationsRequest{
		Notification: "hello world @tagged1@gmail.com @tagged2@gmail.com",
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/retrievefornotifications`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"message":{"`+messages.MESSAGE_MISSING_PARAMS+`":{"teacher":"required"}}}`, rr.Body.String())

	// Test for missing notification field.
	// Should get status code 400 and error message.
	payload = request.ReceieveForNotificationsRequest{
		Teacher: "teacher@gmail.com",
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/retrievefornotifications`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"message":{"`+messages.MESSAGE_MISSING_PARAMS+`":{"notification":"required"}}}`, rr.Body.String())

	// Test for invalid notification field format.
	// Should get status code 400 and error message.
	payload = request.ReceieveForNotificationsRequest{
		Teacher:      "teacher@gmail.com",
		Notification: "hello world @tagged1@gmail.com @tagged2@gmail.com thisshouldnotbehere",
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/retrievefornotifications`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"message":"`+messages.InvalidParamsMessage([]string{"notification"})+`"}`, rr.Body.String())

	// Clean up DB.
	database.CleanupTestDB(db)
	db.Close()
}

// Tests for "/api/register" endpoint.
func Register(t *testing.T) {
	// Init DB.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	database.InitTestDB(db)

	// Init router and middleware.
	r := gin.Default()
	handlers.RegisterMiddlewares(r, db)
	r.POST("/api/register", controllers.Register)

	// Test for POST request.

	// Positive cases.

	// Test for request body with valid pair of teacher and students.
	// Should return status code 204.
	payload := request.RegisterRequest{
		Teacher:  "teacher1@gmail.com",
		Students: []string{"student11@gmail.com", "student12@gmail.com"},
	}
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", `/api/register`, bytes.NewBuffer(jsonValue))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Test for request body with valid pair of student and teachers.
	// Should return status code 204.
	payload = request.RegisterRequest{
		Student:  "student2@gmail.com",
		Teachers: []string{"teacher21@gmail.com", "teacher22@gmail.com"},
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/register`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Test for request body with valid pair of teacher and students, and student and teachers.
	// Should return status code 204.
	payload = request.RegisterRequest{
		Teacher:  "teacher3@gmail.com",
		Students: []string{"student31@gmail.com", "student32@gmail.com"},
		Student:  "student3@gmail.com",
		Teachers: []string{"teacher31@gmai.com", "teacher32@gmail.com"},
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/register`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Negative Cases.

	// Test for invalid pair of teacher and students(missing).
	// Should return status code 400 and error response.
	payload = request.RegisterRequest{
		Teacher: "onlyteacher@gmail.com",
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/register`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"message":"Both fields teacher and students must be present and valid"}`, rr.Body.String())

	// Test for invalid pair of teacher(missing) and students.
	// Should return status code 400 and error response.
	payload = request.RegisterRequest{
		Students: []string{"onlystudent1@gmail.com", "onlystudent2@gmail.com"},
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/register`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"message":"Both fields teacher and students must be present and valid"}`, rr.Body.String())

	// Test for invalid pair of student and teachers(missing).
	// Should return status code 400 and error response.
	payload = request.RegisterRequest{
		Student: "onlystudent@gmail.com",
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/register`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"message":"Both fields student and teachers must be present and valid"}`, rr.Body.String())

	// Test for invalid pair of student(missing) and teachers.
	// Should return status code 400 and error response.
	payload = request.RegisterRequest{
		Teachers: []string{"onlystudent1@gmail.com", "onlystudent2@gmail.com"},
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/register`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"message":"Both fields student and teachers must be present and valid"}`, rr.Body.String())

	// Test for wrong teacher email format.
	// Should return status code 400 and error response.
	payload = request.RegisterRequest{
		Teacher:  "wrong@format",
		Students: []string{"s1@gmail.com", "s2@gmail.com"},
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/register`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"message":"`+messages.INVALID_TEACHER_EMAIL_FORMAT+`"}`, rr.Body.String())
	// Test for one wrong student email format.
	// Should return status code 400 and error response.
	payload = request.RegisterRequest{
		Teacher:  "t1@gmail.com",
		Students: []string{"wrong@format", "s2@gmail.com"},
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/register`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"message":"`+messages.INVALID_STUDENT_EMAIL_FORMAT+`"}`, rr.Body.String())

	// Test for wrong student email format.
	// Should return status code 400 and error response.
	payload = request.RegisterRequest{
		Student:  "wrong.format",
		Teachers: []string{"t1@gmail.com", "t2@gmail.com"},
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/register`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"message":"`+messages.INVALID_STUDENT_EMAIL_FORMAT+`"}`, rr.Body.String())

	// Test for one wrong teacher email format.
	// Should return status code 400 and error response.
	payload = request.RegisterRequest{
		Student:  "s1@gmail.com",
		Teachers: []string{"t1@gmail.com", "wrongforma.t"},
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", `/api/register`, bytes.NewBuffer(jsonValue))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"message":"`+messages.INVALID_TEACHER_EMAIL_FORMAT+`"}`, rr.Body.String())

	// Clean up DB.
	database.CleanupTestDB(db)
	db.Close()
}
