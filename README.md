# web-go-hello-app (MyQ)
#### educational web-app based on and inspired by Trevor Sawler "Bookings" Go project on [udemy](https://www.udemy.com/course/building-modern-web-applications-with-go/)

#### App
MyQ allows me to add and search specific queries that I forget.   

I always have a trouble to remember 'How do I cut word in vim?', 'What was that command for ARP cache?' or 'What I've typed in Google that directed me to that link on SO?'.  
Well, now i can add specific query as part of my own '*search vocabulary*' and provide consice answer or short description with source link. Without a SMM lifestories in the blog pages. And/or without Google's understanding of my understanding of what I actually want to find or what I mean.

1. Add your query:
<img width="1440" alt="image" src="https://user-images.githubusercontent.com/45435157/206728981-8ae9a398-67c5-465e-9847-5f2ae1f4d207.png">

> Of course, look for an answer on Google first. This is trivial example, but still, for more specific searches it would required you to scroll and look through those purple pages or find something in your bookmarks.

2. Look for it later:
<img width="1440" alt="Screenshot 2022-12-09 at 16 50 35" src="https://user-images.githubusercontent.com/45435157/206728617-712e5aaf-ea9b-4edd-954a-b1d388e0c895.png">

back-end: [chi](https://github.com/go-chi/chi), [scs](https://github.com/alexedwards/scs), [nosurf](https://github.com/justinas/nosurf), [pgx](https://github.com/jackc/pgx), [postgres](https://hub.docker.com/_/postgres) image, docker  
front-end: [http/tmpl](https://pkg.go.dev/html/template), [bootstrap](https://getbootstrap.com/), [sweatAlert2](https://sweetalert2.github.io/)  

Run appication  
Start postgres container:
```
docker run --name postgres \
  --network=myq_network \
  --network-alias=postgres \
  -v pgdata:/var/lib/postgresql/data \
  -p 5432:5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=db \
  -d \
  postgres
```  
  
> don't forget to specify postgres variables according to your environment.  

Create DB users and searches tables. I used a custom sql script against running container:  
```cat db_scripts/CreateSearchAndUserTables.sql | docker exec -i postgres psql -U postgres -h localhost -p 5432 -d db```
But you can dump .sql file to the container and run it from within postgres container:  
```
docker cp ./dump.sql postgres:/docker-entrypoint-initdb.d/dump.sql  
docker exec -u postgres postgres psql db postgres -f docker-entrypoint-initdb.d/dump.sql  
```  
General form to the last line of code is:  
```
docker exec -u postgresuser containername psql dbname postgresuser -f /container/path/file.sql
```
Users table:  
   id   |   username   |   email   |   password   |   created_at   |   updated_at   |   access_level   
--------+--------------+-----------+--------------+----------------+----------------+-----------------
Searches table:  
   id   |   user_id   |   query   |   query_link   |   description   |   created_at   |   updated_at
--------+-------------+-----------+----------------+-----------------+----------------+---------------
