# mod

![ci](https://github.com/rumpl/mod/workflows/ci/badge.svg)
[![codecov.io](https://codecov.io/github/rumpl/mod/coverage.svg?branch=master)](https://codecov.io/github/rumpl/mod?branch=master)

A better `go mod init`.

This little tool will help you initialize a go project with go modules. This is what it does:

* initializes go with modules
* creates a simple Makefile
* creates a "Hello World" main.go file
* creates a README file
* initializes a git repository

## Intallation

```
go get -u github.com/rumpl/mod
```

## Usage

```
$ mod init awesome-project
All done, happy hacking!

$ ls -l awesome-project
total 32
-rw-r--r--  1 djordjelukic  staff  41 May 20 20:38 Makefile
-rw-r--r--  1 djordjelukic  staff   7 May 20 20:38 README.md
-rw-r--r--  1 djordjelukic  staff  38 May 20 20:38 go.mod
-rw-r--r--  1 djordjelukic  staff  72 May 20 20:38 main.go
```

## License

[MIT](https://rumpl.mit-license.org)
