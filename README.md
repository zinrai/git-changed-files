# git-changed-files

A CLI tool that detects changed files in a specific directory via git diff.

Designed to be embedded in workflow engines such as GitHub Actions, passing changed file paths to subsequent steps.

## Usage

```bash
git-changed-files -dir configs/app -ext .json -ref origin/main
```

| Flag | Required | Default | Description |
| --- | --- | --- | --- |
| `-dir` | no | `.` (current directory) | Directory to detect changes in |
| `-ext` | no | - | Filter by file extension (e.g. `.json`) |
| `-max` | no | `0` (unlimited) | Maximum number of changed files allowed. Fails if exceeded |
| `-ref` | yes | - | Git ref to compare against |

## Output

Detected file paths are printed to stdout, one per line.

```
configs/app/prod.json
```

On error, a message is printed to stderr and the process exits with code 1.

## Exit Codes

| Code | Meaning |
| --- | --- |
| 0 | Success |
| 1 | Error (no changed files, limit exceeded, etc.) |

## Usage in GitHub Actions

```yaml
- uses: actions/checkout@v4
  with:
    fetch-depth: 0

- name: Detect changed config
  id: detect
  run: |
    if [ "${{ github.event_name }}" = "workflow_dispatch" ]; then
      CONFIG="${{ inputs.config-file }}"
      if [ ! -f "$CONFIG" ]; then
        echo "file does not exist: $CONFIG"
        exit 1
      fi
    else
      CONFIG=$(git-changed-files -dir configs/app -ext .json -max 1 -ref "${{ github.event.pull_request.base.sha }}")
    fi
    echo "config-file=$CONFIG" >> "$GITHUB_OUTPUT"
```

## License

This project is licensed under the [MIT License](./LICENSE).
