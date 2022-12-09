# web-go-hello-app (MyQ)
#### educational web-app based on and inspired by Trevor Sawler "Make reservertion" Go project on [udemy](https://www.udemy.com/course/building-modern-web-applications-with-go/)

#### App
MyQ allows you to add and search specific queries that you forget.
Maybe not you, but me definitely.
I always have a trouble to remember 'How do I cut word in vim?', 'What was that command for ARP cache?' or 'What I've typed in Google that directed me to that link on SO?'.
Well, now i can add specific query as part of my own 'search vocabulary' and provide consice answer or short description with source link, without a SMM lifestories in the blog pages. And/or without Google's understanding of my understanding of what I actually want to find or what I mean.

1. Add your query:
<img width="1440" alt="image" src="https://user-images.githubusercontent.com/45435157/206728981-8ae9a398-67c5-465e-9847-5f2ae1f4d207.png">

> Of course, look for an answer on Google first. This is trivial example, but still, for more specific searches it would required you to scroll and look through those purple pages or find something in your bookmarks.

2. Look for it later:
<img width="1440" alt="Screenshot 2022-12-09 at 16 50 35" src="https://user-images.githubusercontent.com/45435157/206728617-712e5aaf-ea9b-4edd-954a-b1d388e0c895.png">

#### Stack:

back-end: [chi](https://github.com/go-chi/chi), [scs](https://github.com/alexedwards/scs), [nosurf](https://github.com/justinas/nosurf), [pgx](https://github.com/jackc/pgx)
front-end: [http/tmpl](https://pkg.go.dev/html/template), [bootstrap](https://getbootstrap.com/), [sweatAlert2](https://sweetalert2.github.io/)
	
#### Start app:
##### Build app image:
```
docker build --tag myq .
```
```
docker run --name myq -p 8080:8080 myq
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
