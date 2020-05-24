<!-- omit in toc -->
# brick: Endpoints

- [Project README](../README.md)

## Available Endpoints

Below is a static listing of the available endpoints that may be used with
this application. Unlike the
[atc0005/bounce](https://github.com/atc0005/bounce) project, this application
intentionally does not expose available endpoints via an index page.

| Name                | Pattern                 | Description                                                   | Allowed Methods | Supported Request content types | Expected Response content type |
| ------------------- | ----------------------- | ------------------------------------------------------------- | --------------- | ------------------------------- | ------------------------------ |
| `frontpageEndpoint` | `/`                     | Fallback for unspecified routes.                              | `GET`           | `text/plain`                    | `text/plain`                   |
| `disable`           | `/api/v1/users/disable` | Disable user accounts associated with incoming JSON payloads. | `POST`          | `application/json`              | `text/plain`                   |

Other endpoints are stubbed out, but not yet implemented as of this writing
and likely will not be available until after the v0.1.0 launch.
