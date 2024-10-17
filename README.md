# repo

A simple repository manager.

Clone repositories using `repo URL`. The path to the cloned repository will be
output to standard out. If the path already exists, it will not be re-cloned.

Set `repo`'s prefix using `REPOPREFIX` (default `$HOME/.local/src`).

## Installation

### curl

```sh
curl lesiw.io/repo | sh
```

### go install

```sh
go install lesiw.io/repo@latest
```

## Shell function

To automatically change to the cloned repository directory, add this function to
your shell's rc file.

``` sh
repo() {
    cd "$(command repo $@)"
}
```
