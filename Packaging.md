# Packaging PokiPoki

## Dependencies

In capability notation:

```
pkgconfig(Qt5Core)
pkgconfig(Qt5Sql)
go
```

There are no external Go dependencies for pokic besides the Go standard library.

## Build System

Meson is the preferred build system for PokiPoki. CMake is provided as a courtesy,
and can only build PokiPoki. Automated tests for the C++ library can only be run
using Meson, and automated tests for pokic can only be run using `go test`.

## Testing

pokic:
```
cd poki-compiler
go test ./...
```

libpokipoki:
```
meson _test
cd _test
ninja test
```

## Splitting Packages

PokiPoki should preferably be split into the following packages:

- Library Package (libpokipoki)
- Development Package (libpokipoki headers + pokic)
