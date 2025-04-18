# Go Migration Example Project

## Overview
This project demonstrates how to manage database migrations using Go. Below, you'll find instructions on creating and running migrations, as well as starting the application.

---

# Install Dependencies
```bash
go mod download
```

## Creating Migration Files
To create a new migration file, use the following command:

```bash
migrate create -ext go -dir ./migrations -seq add_age_to_users
```

This will generate a sequential migration file inside the `migrations` directory.

---

## Running Migrations
### Migrate Up
Run all pending migrations:

```bash
go run main.go -up
```

### Migrate Down
Rollback the last migration:

```bash
go run main.go -down
```

### Migrate Last Down
Rollback only the most recent migration:

```bash
go run main.go -last-down
```

### Migrate Specific Down
Rollback a specific migration by providing its filename:

```bash
go run main.go -specific-down "000001_create_users_table"
```

---

## Running the Server
To start the application, run:

```bash
go run main.go
```

This will launch the Go application, allowing it to interact with the database.

---

## Notes
- Ensure all dependencies are installed before running migrations.
- Check logs for any errors during migrations or application startup.

Happy coding! 🚀

