# dfinder

A simple CLI application for finding duplicate files.

## How to build

```bash
go build main.go
```

## TODO

- [ ] Add subcommand for scanning directories
  - [ ] Add flag, which allows to specify root directory for scanning
- [ ] Implement results listing
  - [ ] Add subcommand "list"
  
Query for finding duplicate files in a DB:

```sql
SELECT hashes_id, file_path FROM files WHERE hashes_id IN
(SELECT hashes_id FROM files GROUP BY hashes_id HAVING count(*) > 1) ORDER BY hashes_id;
```
