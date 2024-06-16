package store

const createTableQuery = `create table if not exists scheduler
(
	id  integer primary key autoincrement,
	date char(8),
	title varchar,
	comment varchar,
	repeat varchar
);`

const createIndexQuery = "create index if not exists scheduler_date on scheduler(date);"
