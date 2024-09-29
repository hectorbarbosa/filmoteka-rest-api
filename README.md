This is an educational REST-API project, Docker containerized server `filmoteka` in Golang. Utilizes Swagger and Elasticsearch.

## Installation 
```shell
git clone https://github.com/hectorbarbosa/filmoteka-rest-api.git
```
```shell
cd filmoteka-rest-api 
```
## Create and start postgres and elasticsearch containers in background. Database and tables will be in place
```shell
docker compose up postgres elasticsearch -d
```
##  Start REST-API server
```shell
docker compose up api
```