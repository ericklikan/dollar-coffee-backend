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

#### `POST /auth/register`

This endpoint will register a new user given

##### Request Body

```javascript
{
    "firstName": string (required),
    "lastName" : string (required),
    "email"    : string (required),
    "password" : string (required),
    "phone"    : string (optional)
}
```

##### Response

```javascript
{
    "message": string,
    "token"  : string (on success)
}
```

Returns following status codes:

| Status Code | Description             |
| :---------- | :---------------------- |
| 201         | `CREATED`               |
| 400         | `BAD REQUEST`           |
| 500         | `INTERNAL SERVER ERROR` |

#### `POST /auth/login`

This endpoint will issue a JWT with user id and user role

##### Request Body

```javascript
{
    "email"   : string (required),
    "password": string (required)
}
```

##### Response

```javascript
{
    "message": string,
    "token"  : string (on success)
}
```

Returns following status codes:

| Status Code | Description             |
| :---------- | :---------------------- |
| 200         | `OK`                    |
| 400         | `BAD REQUEST`           |
| 401         | `UNAUTHORIZED`          |
| 500         | `INTERNAL SERVER ERROR` |

### `/menu`

This module is responsible for retrieving the coffees and items that are available in the store

#### `GET /menu`

Parameters:

| Parameter | Description          |
| :-------- | :------------------- |
| `page`    | Optional page number |

##### Response

```javascript
{
    "coffees": [
        {
            "ID"         : number,
            "Name"       : string,
            "Description": string,
            "Price"      : float,
            "InStock"    : boolean
        },
    ]
}
```

### `/purchases/`

This module contains logic for submitting purchases by user, retrieving purchase history. These requests will require a valid JWT.

Required Headers:

| Header          | Description           |
| :-------------- | :-------------------- |
| `Authorization` | `Bearer {Issued JWT}` |

#### `POST /purchases/purchase`

Creates a purchase record for a user using information from the JWT token

##### Request Body

```javascript
{
    "items":[
        {
            "coffeeId": uint (required),
            "options" : string
        },
    ]
}
```

##### Response

```javascript
{
    "message": string
}
```

Returns following status codes:

| Status Code | Description             |
| :---------- | :---------------------- |
| 200         | `OK`                    |
| 400         | `BAD REQUEST`           |
| 401         | `UNAUTHORIZED`          |
| 500         | `INTERNAL SERVER ERROR` |

#### `GET /purchases/{userId}`

Retrieves a page from purchase history for userId

Parameters:

| Parameter | Description          |
| :-------- | :------------------- |
| `page`    | Optional page number |

##### Response

```javascript
{
    "message": string,
    "purchases": [
        {
            "transactionId": uint,
            "amountPaid"   : float32,
            "purchaseDate" : string,
            "items": [
                {
                    "CoffeeId"  : uint,
                    "TypeOption": string
                },
            ]
        },
    ]
}
```

Returns following status codes:

| Status Code | Description             |
| :---------- | :---------------------- |
| 200         | `OK`                    |
| 401         | `UNAUTHORIZED`          |
| 500         | `INTERNAL SERVER ERROR` |

### `/internal/`

This module is responsible for all admin tasks such as updating purchase amount paid, and creating/updating/deleting new coffees available.

Required Headers:

| Header          | Description                    |
| :-------------- | :----------------------------- |
| `Authorization` | `Bearer {Issued JWT as admin}` |

#### `POST /internal/coffee`

Creates a new coffee available in store

##### Request Body

```javascript
{
    "name"       : string,
    "price"      : float,
    "description": string
}
```

##### Response

```javascript
{
    "message": string
}
```

Returns following status codes:

| Status Code | Description             |
| :---------- | :---------------------- |
| 200         | `OK`                    |
| 400         | `BAD REQUEST`           |
| 401         | `UNAUTHORIZED`          |
| 403         | `FORBIDDEN`             |
| 500         | `INTERNAL SERVER ERROR` |

#### `PATCH /internal/coffee/{coffeeId}`

Updates given coffee attributes

##### Request Body

```javascript
{
    "name"       : string,
    "price"      : float,
    "description": string,
    "inStock"    : boolean
}
```

##### Response

```javascript
{
    "message": string
}
```

Returns following status codes:

| Status Code | Description             |
| :---------- | :---------------------- |
| 200         | `OK`                    |
| 400         | `BAD REQUEST`           |
| 401         | `UNAUTHORIZED`          |
| 403         | `FORBIDDEN`             |
| 500         | `INTERNAL SERVER ERROR` |

#### `DELETE /internal/coffee/{coffeeId}`

Deletes coffee where id = coffeeId

##### Response

```javascript
{
    "message": string
}
```

Returns following status codes:

| Status Code | Description             |
| :---------- | :---------------------- |
| 200         | `OK`                    |
| 400         | `BAD REQUEST`           |
| 401         | `UNAUTHORIZED`          |
| 403         | `FORBIDDEN`             |
| 500         | `INTERNAL SERVER ERROR` |

#### `PATCH /internal/purchase/{purchase}`

Updates amountPaid for purchase

##### Request Body

```javascript
{
    "amountPaid": float,
}
```

##### Response

```javascript
{
    "message": string
}
```

Returns following status codes:

| Status Code | Description             |
| :---------- | :---------------------- |
| 200         | `OK`                    |
| 400         | `BAD REQUEST`           |
| 401         | `UNAUTHORIZED`          |
| 403         | `FORBIDDEN`             |
| 500         | `INTERNAL SERVER ERROR` |

## TODO

- Add DB migration scripts
- Add Unit tests
- Add internal endpoint to get information for all purchases
- Add internal endpoint to get all users
- Refactor data layer
