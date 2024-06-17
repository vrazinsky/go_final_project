package handlers

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
