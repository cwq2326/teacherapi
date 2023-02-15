package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySqlConfig struct {
	User     string
	Password string
	Port     string
	Host     string
	Name     string
}

func ConnectDB(config *MySqlConfig) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.User, config.Password,
		config.Host, config.Port, config.Name)
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Successfully connected to database")
		return db
	}
}

func ConnectTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root:123@tcp(localhost:3306)/test")

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Successfully connected to database")
		return db
	}
}

func DisconnectDB(db *sql.DB) {
	db.Close()
}

func InitDB(db *sql.DB) {
	// Create relations if not yet exist.
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS teachers (email VARCHAR(255) PRIMARY KEY)")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students
					  (email VARCHAR(255) PRIMARY KEY,
					   suspended TINYINT(1) NOT NULL DEFAULT 0)`)
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS teaches
					  (teacher VARCHAR(255), student VARCHAR(255),
					   PRIMARY KEY(teacher,student),
					   FOREIGN KEY (teacher) REFERENCES teachers(emaiL) ON DELETE CASCADE,
					   FOREIGN KEY (student) REFERENCES students(email) ON DELETE CASCADE)`)
	if err != nil {
		panic(err.Error())
	}
}

func InitTestDB(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS teachers (email VARCHAR(255) PRIMARY KEY)")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students
					  (email VARCHAR(255) PRIMARY KEY,
					   suspended TINYINT(1) NOT NULL DEFAULT 0)`)
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS teaches
					  (teacher VARCHAR(255), student VARCHAR(255),
					   PRIMARY KEY(teacher,student),
					   FOREIGN KEY (teacher) REFERENCES teachers(emaiL) ON DELETE CASCADE,
					   FOREIGN KEY (student) REFERENCES students(email) ON DELETE CASCADE)`)
	if err != nil {
		panic(err.Error())
	}
}

func CleanupTestDB(db *sql.DB) {
	_, err := db.Exec("DROP TABLE teaches")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("DROP TABLE teachers")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("DROP TABLE students")
	if err != nil {
		panic(err.Error())
	}
}
