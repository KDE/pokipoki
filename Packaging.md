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

Meson is the preferred build system for PokiPoki.

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
