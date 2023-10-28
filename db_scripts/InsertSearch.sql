\c db

insert into searches ( 
		id, 
		user_id, 
		query, 
		query_link, 
		description, 
		created_at, 
		updated_at
)
values ( 1, 1, 'add query to myq', '', 'enter key and description', current_timestamp, current_timestamp );
