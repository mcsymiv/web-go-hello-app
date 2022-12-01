# web-go-hello-app

### App
#### Build app image:
```
docker build --tag web_go_hello_app .
```
#### Start app:
```
docker run --name web_go_hello_app -p 8080:8080 web_go_hello_app
```

### Postgres:
#### Start psql:
```
docker run --name postgres \
	-v pgdata:/var/lib/postgresql/data
	-p 5432:5432 \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_PASSWORD=password \
	-e POSTGRES_DB=db \
	-d \
	postgres
```
#### Enter postgres container:
```
docker exec -it postgres bash
```
#### Connect to psql:
```
psql -U postgres -h localhost -p 5432
```
