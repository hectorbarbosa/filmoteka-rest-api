This is educatinal REST-API project `filmoteka` in Golang and PostgreSQL

## Installation 
```bash
git clone https://github.com/hectorbarbosa/filmoteka-rest-api.git
```
```bash
cd filmoteka-rest-api 
```
## Config 
- Server port by default :8080. Edit .env file to change it if needed. 
- `postgreSQL` credentials are located in .env file and in Makefile (`postgres`, `postgres`, 5432 by default)

## Create Database and required tables 
```bash
make createdb
```
```bash
make createtables 
```
##  Compile project
```bash
make
```
##  Run project
```bash
make run