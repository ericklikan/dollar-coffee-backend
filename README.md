# Dollar Coffee Shop Backend

This repository contains the source code for the REST API for Dollar Coffee Shop.

## Installation and Start Up

### Required Libraries

- golang.org/doc/install (1.13.0)
- postgresql.org/docs/11/tutorial-start.html (11.5)

### Instructions

1. [Set up a postgresql database](https://www.postgresql.org/docs/11/tutorial-createdb.html) with a user/password. This database can be hosted on your machine as localhost or externally.

2. Change your directory to the root of this repository and create a `.env` file using `.example.env` as a template. Using the database information you set up previously, create the database URL using `postgresql://{db_user}:{db_password}@{host}:{port}/{db_name}`. If you're using your local machine as the database, use host: `localhost` and port: `5432` (default)

3. When editing the `.env` file, you should also create a secure `token_password` that should not be shared. This password will be used to sign jwts issued by the `/auth/` module.

4. Run the following command from the repository root to install dependencies:

   ```bash
   make deps
   ```

   Run the following command to run unit tests:

   ```bash
   make test
   ```

   Run the following command to run server locally:

   ```bash
   make run
   ```

Note: when running the app, if the environment variable `PORT` is not set, it will default to run on port 5000

## API Documentation

This REST API is split up into several modules:

### `/auth/`

This module is responsible for user registration, and user authentication, where it will issue a jwt that contains information about the user and user role.

### `/menu/`

This module is responsible for

### `/purchases/`

### `/internal/`
