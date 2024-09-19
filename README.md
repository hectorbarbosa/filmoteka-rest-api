This is an educatinal REST-API project, Docker containerized server `filmoteka` in Golang and PostgreSQL.

## Installation 
```shell
git clone https://github.com/hectorbarbosa/filmoteka-rest-api.git
```
```shell
cd filmoteka-rest-api 
```
## Create and start postgres container in background. Database and tables will be in place
```shell
docker compose up postgres -d
```
##  Start REST-API server
```shell
docker compose up api
```