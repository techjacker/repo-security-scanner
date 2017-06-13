# repo-security-scanner

- CLI tool that finds secrets accidentally committed to a git repo, eg passwords, private keys
- Run it against your entire repo's history by piping the output from `git log -p`

-----------------------------------------------------------

## Installation
1. [Download](../../releases) the latest stable release of the CLI tool for your architecture
2. Extract the tar and move the ```scanrepo``` binary to somewhere in your `$PATH`, eg `/usr/bin`

-----------------------------------------------------------

## Usage

Check the entire history of the current branch for secrets.

```
$ git log -p | scanrepo

------------------
Violation 1
Commit: 4cc087a1b4731d1017844cc86323df43068b0409
File: web/src/db/seed.sql
Reason: "SQL dump file"

------------------
Violation 2
Commit: 142e6019248c0d53a5240242ed1a75c0cc110a0b
File: config/passwords.ini
Reason: "Contains word: password"

...
```

-----------------------------------------------------------
### Add false positives to `.secignore`

```
$ cat .secignore
file/that/is/not/really/a/secret/but/looks/like/one/to/diffence
these/pems/are/ok/*.pem
```

[See example in this repo](./.secignore).


-----------------------------------------------------------
### Notifications
Work in progress.

#### Local testing of notifications
```
# launch server and elasticsearch
. .env && docker-compose up -d server elasticsearch

# Elastalert expects the elasticsearch index it is monitoring to exist otherwise it will error. The server creates the index in elasticsearch the first time it writes a log.
make test-run-offenses

# launch elastalert
docker-compose up -d elastalert

# run elastalert test rule utility within elastalert container
docker exec -it <elastlaert_container_hash> sh
elastalert-test-rule --config $ELASTALERT_CONFIG --count-only "$RULES_DIRECTORY/new_violation.yaml"
```
