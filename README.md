# Tags

Tag directories using rules.

# Motivation

I needed a way to easily identify the contents of a directory, primarly code
project directories, to perform some contextual actions. After writing this
logic twice in bash, I decided to turn `tags` into a single-purpose utility
program; GNU style.

Here's some examples of how I use `tags`:

**Run actions**

I use the result of `tags` to populate an [fzf] menu with only entries that
are relevant to the current directory. For example, inside a Go project
directory, `fzf` will show entries to build, run, and test the project using
the `go` command. In a React directory, it will show entries to run the dev
server, build, update dependencies, etc.

I have this bound to a tmux keymap, giving me global, contextual run actions.

**Shell prompt**

Many shell prompts use their own logic and rules for detecting VCS, languages,
tools, etc. But I wanted to unify that logic with my run actions.

Tags allows me to have a centralized config for all my rules, and makes
detection consistent between my run actions and my shell prompt.

# Installation

_TODO_: once initial version is complete

# Getting Started

Running `tags` will output all the tags that match the current directory:

```php
$ tags
make
go
git
docker
```

You may also pass a different directory to get its tags instead:

```php
$ tags some/other/path
js
react
ts
css
```

Tags are defined using rules, which can be managed using the `tags add` and
`tags rm`:

```
tags add make file_exists Makefile
tags add www in_path /var/www/
tags add react file_contains package.json react
```

Tags comes with plenty of bundled help pages. See `tags help` for more
information.

# License

This project is licensed under the [GPL-3.0 license](./LICENSE).

[fzf]: https://github.com/junegunn/fzf
