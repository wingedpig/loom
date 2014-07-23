loom
============

A pure Go implementation of several SSH commands, inspired by [Python's Fabric](http://www.fabfile.org/).

With loom, you can run commands as well as put and get files from remote servers, over SSH.

For documentation, check [godoc](http://godoc.org/github.com/wingedpig/loom).

## TODOs
* Examples
* In Run(), use pipes instead of CombinedOutput so that we can show the output of commands more interactively, instead of now, which is after they're completely done executing.
* Handle wildcards in Get()
* Better error checking in Get()
