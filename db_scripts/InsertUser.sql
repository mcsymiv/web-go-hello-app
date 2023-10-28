\c db

insert into users (
		username,
		email,
		password,
		created_at,
		updated_at
) values (
		'mcs',
		'mcs@mcs.com',
		'password',
		current_timestamp,
		current_timestamp
)
