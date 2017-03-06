# repo-security-scanner


## Installation

1. [Download](../../releases) the latest stable release of the CLI tool for your architecture
2. Extract the tar and move the ```scanrepo``` binary to somewhere in your `$PATH`, eg `/usr/bin`

-----------------------------------------------------------

## Example Usage

Check the entire history of the current branch for secrets.

```
$ git log -p | scanrepo
```
