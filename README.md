# Introduction

[![Build Status](https://travis-ci.org/bradleyfalzon/abicheck.svg?branch=master)](https://travis-ci.org/bradleyfalzon/abicheck) [![Coverage Status](https://coveralls.io/repos/github/bradleyfalzon/abicheck/badge.svg?branch=master)](https://coveralls.io/github/bradleyfalzon/abicheck?branch=master) [![GoDoc](https://godoc.org/github.com/bradleyfalzon/abicheck?status.svg)](https://godoc.org/github.com/bradleyfalzon/abicheck)

`abicheck` is a tool to check for the introduction of backwards incompatible changes.

`abicheck`:
- Guarantees that all consumers of a library will still build without failure
- Only checks exported declarations
- There are no false positives (if there are, it's a bug)
- Not every backwards incompatible change can be detected, swapping argument parameters and other changes still need to
    be considered by the library author
- Can be simply consumed as a library
- Is in its infancy, see [Status](#status), feedback and review appreciated

Secondary tasks could include:
- Detecting current semver and suggesting an appropriate increase
- Listing all changes for help in writing release notes/commit messages.

Try at [abicheck.bradleyf.id.au](https://abicheck.bradleyf.id.au/) or via CLI:

```
go get -u github.com/bradleyfalzon/abicheck/cmd/abicheck
cd /your/project/dir/with/committed/changes
abicheck
```

# Proposed Arguments

`abicheck` also comes with a command line tool, as well as being used as a library, the following are the proposed flags
and arguments for the command line tool.

```
-vcs (auto|git|svn|hg|bzr|etc)  - Version control system to use (default: auto)
-rev FROM...TO                  - Revisions to check as before and after (default: if unstaged changes, check those, else check last two commits)
-vcsDir path                    - Path to root VCS directory (default: let VCS tool search)
-nonBreaking                    - Show non-breaking changes as well as breaking (default: false)

abicheck        # current package only
abicheck ./...  # check subdirectory packages
```

Another tool, called `abichanges` may also be included which will list all detected changes to assist in producing
release notes.

# Status

`abicheck` is currently under heavy development and refactoring. This initial version was a proof of concept and shortcuts were taken. The current tasks are focused on (but not limited to):

- Add type checking to analyse inferred types
- Choosing of import paths as the first argument, similar to other tools (no argument means just current directory, else
    support `./...` and specifying)
- Adding Mercurial, SVN and potentially other VCS systems
- Improve VCS options such as:
    - Detection of VCS and flag to overwrite
    - Choosing base VCS path to allow running for a different directory
    - Detecting if there's unstaged changes (currently only checks committed changes) and testing those
    - Choosing the versions to compare (e.g. via flag `-rev HEAD~1...HEAD` or similar, staying VCS agnostic)
    - Filtering `vendor/` directories (if this is the best place to do it, or leave it to go/type ast packages)
    - Check subdirectories if ran from a subdirectory of the VCS (currently checks all committed code)
- Add docs, flow diagram and fixing of existing docs
- Improve output formats (such as vim quickfix)
- Move these tasks to GitHub issues
- Improve test coverage and move away from golden masters (it was just a quick hack, not a long term solution)
- Once all other steps have been completed, performance will be investigated

# Testing

This uses golden masters for the tests, currently (and only due to time constraints) `testdata/` directory contains `before.go`
and `after.go`, which are before and after versions of a test package, each time `go test` is ran, the output is compared to
`testdata/exp.txt`, which should not change.

If adding new test cases, you should expect the test to fail as the code changes should create a difference with `exp.txt`.
Then, you'll need to update the golden master (see below), and commit those changes. If you add a new test case to `before.go` and
`after.go`, and the tests still pass, you've uncovered a bug within `abicheck` which will need a code change to fix, once
code has change, the tests should fail, so update the master, review all your changes and commit.

- This uses golden master `testdata/exp.txt` for the tests
- Run tests with: `go test`
- Update master with: `go test -args update`
- Alternatively to do a test run: `go install && ( cd testgit; ./make.sh && abicheck )`
