# auth-server

The auth-server has multiple roles, mainly to:

- provide authentication to clients and hosts on the same network
- keep tabs on all hosts on the network and their current status
- matchmake a host and a client, once the client wishes it.

## Technologies

The auth-server is built using:

- Java
- Spring Boot
- Lombok
- JPA with PostgreSQL
- Flyway for migrations
- JWT Tokens for authentication
- Docker & Docker Compos
