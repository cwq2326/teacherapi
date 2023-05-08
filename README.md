## Golang Assessment for Division Tech Assessment (GDS / DCube)
### Name: Chua Wen Quan
```
Go v1.20
MySQL v 8.0.32
```
---
### Instructions to host locally
---
#### Clone repository
* Clone this repo by running the command `git clone https://github.com/cwq2326/govtech.git`

#### Configure .env file
* Copy `.env.example` and rename it to `.env`
* Edit `.env` to contain your mysql (including database name for testing) and router configurations

#### Install dependencies
* From root directory, run the command `go get -u -d ./...`

#### Configure mysql database
* Run mysql locally.
  * eg. `mysql -u root -p`
* Create database used in `.env` file, it should be same as `DB_NAME`
  * eg. `CREATE DATABASE <DB_NAME>` (Replace <DB_NAME> with database name in mysql)

#### Run API server
* From root directory, run the command `cd cmd/main && go run .`
---
### Instructions to test
---
* You should have done the `.env` file configuration step
* Run mysql locally.
  * eg. `mysql -u root -p`
* Create database used in `.env` file, it should be same as `DB_TEST_NAME`
  * eg. `CREATE DATABASE <DB_TEST_NAME>` (Replace <DB_TEST_NAME> with database name in mysql)
* From root directory, run the command `go test -v ./tests`
