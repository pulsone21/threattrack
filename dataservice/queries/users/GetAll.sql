SELECT  
    users.id as UserID,
    users.firstname as Firstname,
    users.lastname as Lastname,
    users.email as Email,
    UNIX_TIMESTAMP(users.created_at) as createdAt,
    users.fullname as Fullname
FROM users 
ORDER BY users.created_at DESC
LIMIT ? OFFSET ?;
