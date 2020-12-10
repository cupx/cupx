# CupX

[![PkgGoDev](https://pkg.go.dev/badge/cupx.github.io)](https://pkg.go.dev/cupx.github.io)

This module [`cupx.github.io`](https://pkg.go.dev/mod/cupx.github.io) contains a Cup of eXtensible libraries.

## Installation

```
go get -u cupx.github.io
```

## Package index

Summary of the packages provided by this module:

- [xlog](https://pkg.go.dev/cupx.github.io/xlog): Package `xlog` provides an extensible log library.
- [xdns](https://pkg.go.dev/cupx.github.io/xdns): Package `xdns` provides an extensible dns library.
- [xacme](https://pkg.go.dev/cupx.github.io/xacme): Package `xacme` provides an extensible acme library.

## FAQ

### Why not use [github.com/cupx/cupx](https://github.com/cupx/cupx) as module path?

The path [cupx.github.io](https://cupx.github.io) is shorter than [github.com/cupx/cupx](https://github.com/cupx/cupx).

### Why not use [pkg.cupx.net](https://pkg.cupx.net) as Module path?

I currently cannot ensure that the `cupx.net` domain name can be renewed in time. 
If I confirm that the `cupx.net` domain name can be renewed all the time in the future, the module 
path of cupx may be switched to [pkg.cupx.net](https://pkg.cupx.net) from a certain version.
