# git-changed-files

A CLI tool that detects changed files in a specific directory via git diff.

It prints one path per line for feeding subsequent commands, in local scripts and CI jobs alike.

## Usage

```bash
git-changed-files -dir configs/app -ext .json -ref origin/main
```

| Flag | Required | Default | Description |
| --- | --- | --- | --- |
| `-dir` | yes | - | Directory to detect changes in. Use `.` for the whole repository |
| `-ext` | no | - | Filter by file extension (e.g. `.json`) |
| `-max` | no | `0` (unlimited) | Maximum number of changed files allowed. Fails if exceeded |
| `-ref` | yes | - | Git ref to compare against. An all-zeros SHA is treated as the empty tree |

## Output

Detected file paths are printed to stdout, one per line.

```
configs/app/prod.json
```

Deleted files are not reported. Only files that exist at HEAD appear in the output. When no files changed, nothing is printed and the exit code is 0.

An all-zeros `-ref`, which git hooks and workflow engines pass as the old rev when a ref is newly created, is diffed against the empty tree. Every tracked file under the directory at HEAD is then listed.

On error, a message is printed to stderr and the process exits with code 1.

## Exit Codes

| Code | Meaning |
| --- | --- |
| 0 | Success (including no changed files) |
| 1 | Error (limit exceeded, git failure, etc.) |

## Examples

Lint only the files a branch touched, instead of the whole repository:

```bash
$ git-changed-files -dir configs -ext .yaml -ref origin/main | xargs -r yamllint
```

Record what a pipeline run is about to process. The file feeds the next step and stays inspectable afterwards:

```bash
$ git-changed-files -dir configs -ref "$BEFORE_SHA" > changed.txt
```

Skip an expensive step when nothing under the directory changed:

```bash
if [ -n "$(git-changed-files -dir docs -ref origin/main)" ]; then
    make html
fi
```

Enforce that a change touches at most one config, and pick it up. More than one is an error, none is a no-op:

```bash
config="$(git-changed-files -dir configs/app -ext .json -max 1 -ref origin/main)"
if [ -z "$config" ]; then
    echo "no config changed"
    exit 0
fi
```

## License

This project is licensed under the [MIT License](./LICENSE).
