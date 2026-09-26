# Importing Toyota's current catalog

The catalog importer scrapes Toyota's vehicle listing page through Firecrawl and
upserts one `cars` row per model and trim. It imports records for the current
calendar year and the following model year; records outside that range or
otherwise invalid abort the import.

## Configure and run

Set `DATABASE_URL` and `FIRECRAWL_API_KEY` in the environment or in the
repository's `.env` file. Apply the project migrations, then run:

```sh
go run ./cmd/import-cars
```

The command reports how many catalog records it imported or updated. It ensures
the catalog identity index exists and is safe to rerun: records are matched
case-insensitively by make, model, model year, and trim. Cargo-volume values are
not changed by this import.

If the identity-index migration fails because matching duplicate rows already
exist, resolve those duplicates before applying the migration; the migration
does not delete existing data.

The scraper requires a working Firecrawl API key and aborts without changing the
database if Firecrawl fails, the catalog is empty, or any record is invalid or
outside the current/next-year range.
