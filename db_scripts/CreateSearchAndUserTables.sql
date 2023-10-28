\c db

create table users (
		id serial not null unique,
		username varchar ( 255 ) unique not null,
		email varchar ( 255 ) unique not null,
		password varchar ( 255) not null,
		created_at timestamp not null,
		updated_at timestamp not null,
		access_level smallint not null
);

create table searches (
		id serial not null,
		user_id integer references users (id),
		query varchar ( 255 ) not null,
		query_link varchar ( 255 ),
		description varchar ( 255 ) not null,
		created_at timestamp not null,
		updated_at timestamp not null
);



