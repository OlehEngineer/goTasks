**REST API implementation.**

API works with next HTTP methods:

    1. GET
       1. GET ***/api/users?page=<value>&limit=<value>***
            This endpoint must contain page number and one page limitation value. Pagination implemented
       2. GET ***/api/users/:id***
            This endpoint must contain user id. Basic Authentication implemented.
       3. GET method return JSON data:
            {
            "id": 0,
            "nickname": "NickName"",
            "name": "Name",
            "lastname": "LastName",
            "email": "Email",
            "status": "active/passive",
            "created_at": time.Time,
            "updated_at": time.Time,
            "deleted_at": time.Time,
            "likes": 00
        }

    2. POST ***/api/users***
            JSON data required:{
            **"nickname"**: "NickName",
            **"name"**: "Name",
            **"lastname"**: "LastName",
            **"password"**: "Password"
        }
        Basic Authentication implemented.

    3. PUT ***/api/users/:id***
            User id require. Basic Authentication implemented.
            Return the same JSON data as GET method

    4. DELETE ***/api/users/:id***
            User id require. Basic Authentication implemented.
            