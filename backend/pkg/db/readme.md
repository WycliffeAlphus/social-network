# Database Migrations (SQLite + golang-migrate)

This folder contains SQL migration files used to manage the database schema.

We use [`golang-migrate`](https://github.com/golang-migrate/migrate) to apply version-controlled schema changes to our SQLite database.

---

## Creating a New Migration

To create a new migration (e.g., to add or change a table), run:

```bash
migrate create -ext sql -dir pkg/db/migrations -seq create_users_table
```

This generates two files:
```bash
pkg/db/migrations/
├── 00000X_create_users_table.up.sql   # SQL to apply schema change
└── 00000X_create_users_table.down.sql # SQL to roll back schema change
```

Edit the .up.sql file to define what should happen when the migration is applied, and the .down.sql file to define how to undo that change.



## Fixing a Dirty Migration
If a migration fails and leaves the database in a dirty state, you can reset the version manually:
```bash
migrate -database sqlite://pkg/db/data/app.db -path pkg/db/migrations force 0
```
Make sure you're using the SQLite-enabled version of the migrate CLI.

If needed, delete the pkg/db/data/app.db file to start fresh during development.