# threattrack

CSIRT Incident Management Tool which combines SOAR und Project Management

## How to run the app

- Create an .env file eg:

```bash
FRONTEND_ADRESS=http://Frontend
FRONTEND_PORT=5050
BACKEND_PORT=5051
BACKEND_ADRESS=http://Backend
MYSQL_ADDR=DB
MYSQL_PORT=3306
MYSQLROOTPW=root
MYSQL_USER=contentUsr
MYSQL_PW=root
MYSQL_DBNAME=threattrack
```

- in the root directory run

```bash
docker compose up --build
```

- Note that the first build takes arround 3 Minutes to run.
