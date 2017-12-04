# Simple REST API
Simple CRUD REST API built with Go & PostgreSQL

### Third party libraries
- [julienschmidt/httprouter](github.com/julienschmidt/httprouter)

### Setup

```sh
> psql
> CREATE DATABASE twit;
```

```sql
CREATE TABLE users (
   ID INT PRIMARY KEY     NOT NULL,
   USERNAME       TEXT    NOT NULL,
   EMAIL          TEXT    NOT NULL,
   PASSWORD       TEXT    NOT NULL,
   CREATED_AT     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE twit (
   ID INT PRIMARY KEY    NOT NULL,
   USER_ID        INT    NOT NULL,
   BODY           TEXT   NOT NULL,
   CREATED_AT     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   UPDATED_AT     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

```sh
> \l # list dbs
> \d # list tables
> \c twit # connect to twit db
> SELECT * FROM twit;
```

### TODO
- [x] Install Postgres
- [x] Pick a router standard or third party (httprouter)
- [x] Define endpoints / routes
- [x] Setup database tables
- [ ] Create index handler (fetch all twits)
- [ ] Create authentication
  - [ ] Create Signup
  - [ ] Create Login
  - [ ] Create Logout
- [ ] Create authorization middleware
- [ ] Create POST twit
- [ ] Create EDIT twit
- [ ] Create DELETE twit
