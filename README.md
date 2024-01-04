# example

This a sample Go REST API implementation.

## TODOs

- [ ] Add validation using [validator](https://github.com/go-playground/validator);
- [ ] Add JWT instead of stateful tokens, since it's what's most people use;
- [ ] Add a *Posts* resource - users should be able to create new posts, like and comment them;
- [ ] Add a React client.


## Requirements

- [go 1.21](https://go.dev) - the go programming language
- [goose](https://github.com/pressly/goose) - database migration tool
- [reflex](https://github.com/cespare/reflex) - file watcher (optional)
- [make](https://www.gnu.org/software/make/) - tool for generating executables (optional)

## Settings

The project configuration is done using environment variables or a `.env` file. 
All settings are defined in the `.env.example` file at the root of the project.

Enviroment variables validation is done using the awesome [env](https://github.com/caarlos0/env)
package.

The available settings are described below:

#### APP_PORT 
Defines which port the http server listens to. Defaults to `4000`.

#### APP_ENV 
Defines which environemnt the application is currently running. 
Should be set to `development` or `production`. Defaults to `development`.

#### SECRET_KEY
**Required**. Defines the secret key used to sign JWT tokens.

If you have [openssl](https://www.openssl.org/) installed, run the command below to generate
a new secret key:

```bash
openssl rand -base64 32
```

#### DATABASE_PATH

Sets the path to the sqlite database. Defaults to `example.db`.

All of the settings are defined in the `internal/config` package. If a variable is sensitive, like `SECRET_KEY`, it should be set using the `config.Secret` type.

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
