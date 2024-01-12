SELECT 
    users.id as UserID,
    users.firstname as Firstname,
    users.lastname as Lastname,
    users.email as Email,
    users.created_at as CreatedAt,
    users.fullname as Fullname
FROM users 
WHERE users.id = ?
