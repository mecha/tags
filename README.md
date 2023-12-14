# Tags

A command line utility that uses rules to tag directories.

# Usage

When given a path as an argument, `tags` will output the matching tags, one on each line.

```
$> tags /home/myself/projects/react-todo-list-app
js
docker
tailwind
```

The tags are determined using the rules in your [config file](#config). See the next section for more info.

Note that the tags are processed concurrently, so the output order will differ from run to run. Pipe the result through
`sort` or another similar utility if you want to have a predictable order.

# Config

The config file uses JSON syntax, and is expected to be in one of these locations, depending on your platform:

- **Linux**: `$XDG_CONFIG_HOME/tags/rules.json`
- **Windows**: `%APPDATA%/tags/rules.json`
- **MacOS**: `~/Library/Application Support/tags/rules.json` (untested)

The config should contain a JSON object that maps tag names to their rules. Each rule is also a JSON object that maps
the [rule type](#rule-types) to its config.

_Example config:_

```json
{
  "rust": {
    "file_exists": "Cargo.toml"
  },
  "docker": {
    "files_exist": ["Dockefile", "docker-compose.yml"]
  },
  "tailwindcss": {
    "file_exists": "tailwind.config.js",
    "file_contains": {
      "package.json": "\"tailwindcss\":"
    }
  }
}
```

The rules for a single tag are processed in an `OR` fashion. That is, only a single rule needs to match a directory
for the tag to be returned. A tag's rules will stop being evaluated after a match has been found.

## Rule Types

### `file_exists`

Matches directories that contain a specific file, or multiple specific files, matched by their name.

_Example:_

```json
{
  "go": {
    "file_exists": "go.mod"
  },
  "git": {
    "file_exists": ".git"
  },
  "docker": {
    "file_exists": ["Dockerfile", "docker-compose.yml"]
  }
}
```

### `file_contains`

Matches directories that contain specific files that also include a piece of text. Note that the text search is
case-sensitive.

_Example:_

```json
{
  "react": {
    "file_contains": {
      "package.json": "react",
      "index.jsx": "react",
      "index.tsx": "react"
    }
  }
}
```

### `in_path`

Checks if a directory is in a given path.

_Example:_

```json
{
  "wordpress": {
    "in_path": ["~/sites/wp", "/var/www/"]
  }
}
```

# Troubleshooting

If you're having trouble getting a tag to match correctly, use the `-v` flag to enable verbose output. All verbose
output is sent over `stderr`.

```
$> tags /some/path -v

Checking docker tag
=> [file_exists] Dockfile ... no

Checking git tag
=> [file_exists] .git ... yes
```

# License

This project is licensed under the [GPL-3.0 license](./LICENSE).

[fzf]: https://github.com/junegunn/fzf
