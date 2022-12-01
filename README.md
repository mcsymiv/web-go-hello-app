# web-go-hello-app

### Postgres
Start psql
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
Enter postgres container
```
docker exec -it postgres bash
```
Connect to psql
```
psql -U postgres -h localhost -p 5432
```
