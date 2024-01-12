CREATE TABLE IF NOT EXISTS tasks (
        id varchar(255) NOT NULL,
        title varchar(255)  NOT NULL,
        description text  NOT NULL,
        owner_id varchar(255)  NOT NULL,
        incident_id varchar(255)  NOT NULL,
        status enum(
            'Backlog',
            'Doing',
            'Review',
            'Done'
        )  NOT NULL DEFAULT 'Backlog',
        priority enum(
            'Low',
            'Medium',
            'High',
            'Critical'
        )  NOT NULL DEFAULT 'Low',
        PRIMARY KEY (id),
        KEY owner_id (owner_id),
        KEY incident_id (incident_id),
        FOREIGN KEY (owner_id) REFERENCES users (id),
        FOREIGN KEY (incident_id) REFERENCES incidents (id)
    );