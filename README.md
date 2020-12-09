# CupX

[![PkgGoDev](https://pkg.go.dev/badge/cupx.github.io)](https://pkg.go.dev/cupx.github.io)

This module
([`cupx.github.io`](https://pkg.go.dev/mod/cupx.github.io))
contains a Cup of eXtensible libraries.

## Installation

```
go get cupx.github.io 
```

## Package index

Summary of the packages provided by this module:

*   [`xlog`](https://pkg.go.dev/cupx.github.io/xlog): Package
    `xlog` provides an extensible log library.
*   [`xdns`](https://pkg.go.dev/cupx.github.io/xdns): Package
    `xdns` provides an extensible dns library.
*   [`xacme`](https://pkg.go.dev/cupx.github.io/xacme): Package
    `xacme` provides an extensible acme library.

## FAQ

1. Why not use cupx.github.io as Module path?
   - dd
2. Why not use pkg.cupx.net as Module path?
   ```
   为了依赖本模块的构建永不失败, 根据的隐私政策, 有可能cupx.net域名最终因为某种原因忘记续费,而proxy.golang.org又不保证永久不删除缓存的版本
   March 22, 2016 , left-pad https://blog.npmjs.org/post/141577284765/kik-left-pad-and-npm
   ```
  - 期待 proxy.golang.org 未来能够改进模块缓存机制, 类似 npm的(unpublish policy)[https://www.npmjs.com/policies/unpublish] 对于发布超过72, 没有被其他公共包依赖且满足特定开源许可证的模块的版本承诺永久保存, 以保证依赖此模块版本的构建永不失败.
  - 等后续确认cupx.net会一直维持续费或者proxy.golang.org承诺永久保留历史版本缓存后, cupx有可能会从某个版本开始切换到 pkg.cupx.net