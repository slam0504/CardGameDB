# CardGameDB

This project provides a small card database service written in Go. It exposes HTTP endpoints for searching, creating and updating cards. The service uses MySQL for persistence.

## Running locally with Docker Compose

The repository includes a `docker-compose.yml` which starts the application together with a MySQL database:

```bash
docker-compose up --build
```

The API will then be available on `http://localhost:8080` and MySQL on port `3306`.

## Development container

To hack on the service using VS Code you can use the included development container:

1. Install the [Remote Development](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.vscode-remote-extensionpack) extension pack.
2. Open the repository in VS Code and choose **Reopen in Container**.

The container is based on the same `docker-compose.yml` and will download Go modules and set up the Go tools automatically.

## Microservice deployment

`Dockerfile` defines how to build the service container. The `docker-compose.yml` composes the `app` service with a MySQL database, which is suitable for running the service as a small microservice.

Environment variable `DB_DSN` specifies the MySQL connection string. It defaults to `user:password@tcp(localhost:3306)/carddb` when unset.

## Event sourcing

Every domain action publishes an event which is also stored in the `events` table for auditing. The table can be created with:

```sql
CREATE TABLE events (
  id INT AUTO_INCREMENT PRIMARY KEY,
  type VARCHAR(255) NOT NULL,
  payload TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL
);
```

Events are appended whenever the service receives a search, create or update request.
