SELECT     
    tasks.id as ID,
    tasks.title as Title,
    tasks.description as Description,
    users.id as UserID,
    users.firstname as Firstname,
    users.lastname as Lastname,
    users.email as Email,
    users.created_at as CreatedAt,
    users.fullname as Fullname
FROM tasks 
LEFT JOIN
	users ON tasks.id = users.id
WHERE tasks.id = ?
