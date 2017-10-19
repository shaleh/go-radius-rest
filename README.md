# Simple RADIUS user management micro service

The service exposes the following endpoints.

    GET /users expects a list of users that have a `name` and a `disabled` boolean.
    PUT /users/<name>/disable will disable a user in RADIUS
    PUT /users/<name>/reactivate will reactivate a user in RADIUS

That is it. Yes, it is disable/reactivate not disable/enable. This was done to
prevent accidental calls.

Right now it assumes a MariaDB/MySQL server as the backend. There is
no reason that a different SQL server could not be supported in the
long run.

## Configuration
From `example.json`

    {
        "Database": {
            "User": "web",
            "Password": "web",
            "Host": "dbserver.example.com",
            "Name": "radius"
        },
        "Host": "",
        "Port": 8000
    }

On Linux system use the `example-systemd.service` as a start to making
the service run on start up.

## NOTE
The change in user state will hit the RADIUS DB immediately. However,
the client may not be affected for some time. This will be based on
RADIUS timeouts. For instance, if the user is set to have access from
8am to 8pm and it is 4pm currently the user's client may not check for
several hours. This can be adjusted in the database.
