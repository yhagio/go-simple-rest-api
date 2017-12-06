# Simple REST API with Go (Practice)
Simple CRUD REST API built with Go & PostgreSQL

![alt text](https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS1ykK7EvxhHcf3_fo1Bgkpu2tZJXFNgBuFaCmtwwfbCMTC43uVQw "Go")

### Third party libraries
- [julienschmidt/httprouter](github.com/julienschmidt/httprouter)

### Setup

```sh
> psql
> CREATE DATABASE twit;
```

```sql
DROP TABLE twit;
DROP TABLE users;

CREATE TABLE users (
   ID         SERIAL NOT NULL PRIMARY KEY,
   USERNAME   TEXT   NOT NULL,
   EMAIL      TEXT   NOT NULL,
   PASSWORD   TEXT   NOT NULL,
   CREATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE twit (
   ID         SERIAL NOT NULL PRIMARY KEY,
   USER_ID    INT    NOT NULL,
   BODY       TEXT   NOT NULL,
   CREATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   UPDATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (USERNAME, EMAIL, PASSWORD) VALUES ('test', 'test@test.com', 'test');
SELECT * FROM users;

INSERT INTO twit (USER_ID, BODY) VALUES ('1', 'Hello World! First Twit!');
SELECT * FROM twit;
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
- [x] Create index handler (fetch all twits)
- [x] Create one twit handler (fetch one twit)
- [x] Create authentication
  - [x] Create Signup
  - [x] Create Login
  - [ ] Create Logout
- [x] Create authorization middleware
- [x] Create POST twit
- [ ] Create EDIT twit
- [ ] Create DELETE twit
