package store

const createTableQuery = `CREATE TABLE IF NOT EXISTS scheduler
(
	id  integer primary key autoincrement,
	date char(8),
	title varchar,
	comment varchar,
	repeat varchar
);`

const createIndexQuery = "CREATE INDEX IF NOT EXISTS scheduler_date on scheduler(date);"

const addTaskQuery = `INSERT INTO scheduler
(date, title, comment, repeat)
VALUES(:date, :title, :comment, :repeat)
RETURNING id;
`
const getTasksQuery = "SELECT id, date, title, comment, repeat from scheduler where (:filterByTitle = false or title like :searchValue) and (:filterByDate = false or date = :searchValue) order by date asc limit 50"
const getTaskQuery = "SELECT id, date, title, comment, repeat from scheduler where id = :id"
const updateTakQuery = `UPDATE scheduler
SET date=:date, 
title=:title, 
comment=:comment, 
repeat=:repeat
WHERE id=:id;`
const deleteTaskQuery = "DELETE FROM scheduler WHERE id=:id;"
