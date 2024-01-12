CREATE TABLE IF NOT EXISTS  users (
        id varchar(50) NOT NULL,
        firstName varchar(255) NOT NULL,
        lastName varchar(255) NOT NULL,
        email varchar(255) NOT NULL,
        created_at varchar(50) DEFAULT NULL,
        fullname varchar(255) DEFAULT NULL,
        PRIMARY KEY (id)
    )