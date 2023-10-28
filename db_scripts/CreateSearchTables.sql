\c db

create table searches (
		id serial not null,
		user_id integer references users (id),
		query varchar ( 255 ) not null,
		query_link varchar ( 255 ),
		description varchar ( 255 ) not null,
		created_at timestamp not null,
		updated_at timestamp not null
);



