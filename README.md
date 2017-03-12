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
```

-----------------------------------------------------------
### Add false positives to `.secignore`

```
$ cat .secignore
file/that/is/not/really/a/secret/but/looks/like/one/to/diffence
these/pems/are/ok/*.pem
```

[See example in this repo](./.secignore).
