profzf
======

Profzf consists of a server and client (in the same binary). The server sits in the background and scans for your git repositories and caches them locally in an sqlite database.

The client is a simple wrapper around fzf that allows you to quickly search for your projects. Mainly intended for use with `cd`

![profzf](docs/fzf.png)

Usage
=====

Run `profzf server` to start the server. Use the `--project-dir` (can be specified multiple times) to specify the directories to scan for git repositories.

Then run `profzf cd` to output an example command that uses `fzf` and `jq` to cd into the selected project.
