# web-go-hello-app (MyQ)
#### educational web-app inspired by Trevor Sawler [udemy course on Go](https://www.udemy.com/course/building-modern-web-applications-with-go/)

#### App
MyQ allows you to add and search specific queries that you forget. 
I always have a trouble to remember 'How do I cut word in vim?', 'What was that command for ARP cache?' or 'What I've typed in Google that directed me to that link on SO?'.
Well, now i can add specific to my own 'search vocabulary' query and provide consice answer or short description with source link, without a SMM lifestories in the blog. And/or without Google understanding of my understanding of what I actually want to find.

#### Stack: 
back-end: [chi](https://github.com/go-chi/chi), [scs](https://github.com/alexedwards/scs), [nosurf](https://github.com/justinas/nosurf), [pgx](https://github.com/jackc/pgx)
front-end: [http/tmpl](https://pkg.go.dev/html/template), [bootstrap](https://getbootstrap.com/), [sweatAlert2](https://sweetalert2.github.io/)
	
#### Start app:
##### Build app image:
```
docker build --tag web_go_hello_app .
```
```
docker run --name web_go_hello_app -p 8080:8080 web_go_hello_app
```

##### Postgres:
##### Start psql:
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
