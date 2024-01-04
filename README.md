# example

This a sample Go REST API implementation.

## TODOs

- [ ] Add validation using [validator](https://github.com/go-playground/validator).


## Requirements

- [go 1.21](https://go.dev) - the go programming language
- [goose](https://github.com/pressly/goose) - database migration tool
- [reflex](https://github.com/cespare/reflex) - file watcher (optional)
- [make](https://www.gnu.org/software/make/) - tool for generating executables (optional)


## Quick start

You can start the application by running the command below:

```bash
go run cmd/api/*
```


## Endpoints

This example showcases a simple authentication flow.


### 1. Register a new user

To create a new user, make a `POST` request to `/auth/sign-up`, like the example below:

```http
POST /auth/sign-up

{
    "email": "user@example.com",
    "password": "123123123"
}
```  


### 2. Authenticate a user

To authenticate a user (login), make a `POST` request to `/auth/sign-in`, like the example below:

```http
POST /auth/sign-in

{
    "email": "user@example.com",
    "password": "123123123"
}
```

If the authentication is successful, it will return an authentication token.


### 3. Get the current user

To get information about the current user, make a `GET` request to `/auth/me`. The `Authorization` header must be set.

```http
GET /auth/me

Authorization: Bearer <token>
```
