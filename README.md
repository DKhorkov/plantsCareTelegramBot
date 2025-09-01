## Usage

Before usage need to create network for correct dependencies work:

```shell
task -d scripts network -v
```

To stop all docker containers,
use next command:

```bash
task -d scripts docker_stop -v
```

To clean up all created dirs and docker containers,
use next command:

```bash
task -d scripts clean_up -v
```

### Run via docker:

To run app and it's dependencies in docker, use next command:

```bash
task -d scripts prod -v
```

### Run via source files:

To run application via source files, use next commands:

1) Run all application dependencies:

```shell
task -d scripts local -v
```

2) Run application:

```shell
go run ./cmd/main.go
```

## Linters

To run linters, use next command:

```shell
 task -d scripts linters -v
```

## Tests

To run test, use next commands. Coverage info will be
recorded to ```coverage``` folder:

```shell
task -d scripts tests -v
```

To include integration tests, add `integration` flag:

```shell
task -d scripts tests integration=true -v
```

## Benchmarks

To run benchmarks, use next command:

```shell
task -d scripts bench -v
```

## Migrations

To create migration file, use next command:

```shell
task -d scripts makemigrations NAME={{migration name}}
```

To apply all available migrations, use next command:

```shell
task -d scripts migrate
```

To migrate up to a specific version, use next command:

```shell
task -d scripts migrate_to VERSION={{migration version}}
```

To rollback migrations to a specific version, use next command:

```shell
task -d scripts downgrade_to VERSION={{migration version}}
```

To rollback all migrations (careful!), use next command:

```shell
task -d scripts downgrade_to_base
```

To print status of all migrations, use next command:

```shell
task -d scripts migrations_status
```

## Database

To connect to database container, use next command:

```shell
sudo docker exec -it database sh
```

To connect to DB, use next command:

```shell
psql -U $POSTGRES_USER
```

To create backup of database, use next command:

```shell
sudo docker exec database /scripts/backup.sh
```

To restore database from latest backup, use next command:

```shell
sudo docker exec database /scripts/restore.sh
```

To restore database from specific backup, use next command:

```shell
sudo docker exec database /scripts/restore.sh /backups/{{backup_filename}}
```
