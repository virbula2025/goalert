In GoAlert, database migrations follow a structured, sequential process. Unlike the GraphQL API layer, migration scripts are **not auto-generated**; they are hand-coded to ensure data safety and to support complex deployment scenarios like zero-downtime upgrades.

### How Migrations Work

GoAlert uses a custom migration engine (found in `/migrate`) that treats the database schema as a versioned asset.

1.  **Migration Files**: Scripts are stored in `/migrate/migrations`. They can be plain `.sql` files for structural changes (like `CREATE TABLE`) or `.go` files for complex data transformations that require logic.
2.  **Embedded Assets**: Using Go's `embed` package, these migration scripts are compiled directly into the GoAlert binary. This means the application carries its own database "blueprints" with it.
3.  **The `schema_migrations` Table**: When GoAlert connects to a database, it checks a special table called `schema_migrations` to see which version (timestamped ID) the database is currently on.
4.  **Execution**: If the binary contains migrations with a higher ID than the database's current version, it executes the "Up" scripts in order until the database is up to date.



### Are Migration Scripts Generated?

**No.** Migration scripts are **manually written** by developers. 

While tools exist in the industry to "diff" a database and generate a migration, GoAlert avoids this for several reasons:
* **Precision**: Automatic generators often fail at complex tasks like renaming columns or splitting tables without losing data.
* **Zero-Downtime (SWO)**: GoAlert supports a "Switchover" process. Manual scripts allow developers to write migrations that are compatible with both the old and new versions of the code simultaneously.
* **Performance**: For large datasets (like `alert_log`), a developer might need to write a specific SQL query that avoids locking the table for a long period, which an auto-generator wouldn't know how to do.

### The Migration Lifecycle

| Phase | Responsibility |
| :--- | :--- |
| **Authoring** | Developer writes a new `.sql` or `.go` file in `/migrate/migrations`. |
| **Bundling** | The file is automatically embedded into the binary during the build process. |
| **Deployment** | The admin runs the GoAlert binary (or the `goalert-migrate` tool). |
| **Verification** | The engine checks the `schema_migrations` table and applies only the new scripts. |

### Summary of Differences

* **GraphQL Layer**: Generated from a schema (`schema.graphql`) because the mapping between a schema and Go code is predictable and repetitive.
* **Migration Layer**: Hand-written because changing a live database schema is a high-risk operation that requires human judgment to prevent data loss or downtime.