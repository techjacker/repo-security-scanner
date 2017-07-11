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
## Notifications
Work in progress.

### Local Testing
#### Set environment variables needed
Create `env` file and update environment variables.
```
$ cp .env{.example,}
# update .env values
$ vi .env
$ source .env
```

#### Launch containers
```
$ docker-compose up -d
```

#### Run test offenses
```
$ make test-run-offenses
```


### Debugging Elastalert
```
$ docker exec -it <elastalert_container_hash> sh
# run elastalert test rule utility within elastalert container
$ elastalert-test-rule --config $ELASTALERT_CONFIG --count-only "$RULES_DIRECTORY/new_violation.yaml"
$ elastalert-test-rule --alert --config $ELASTALERT_CONFIG "$RULES_DIRECTORY/new_violation.yaml"
# run elastalert in debug mode
$ elastalert --config "$ELASTALERT_CONFIG" --rule "$RULES_DIRECTORY/new_violation.yaml" --debug
```

#### Logs
```
$ tail -f /log/elastalert_new_violation_rule.log
```
