go-mapnik
=========

This is a fork of [omniscale/go-mapnik](https://github.com/omniscale/go-mapnik).

It adds some additional features:

- Support for creating layers and datasources. Implements [niccaluim/go-mapnik@f6bb4d9](https://github.com/niccaluim/go-mapnik/commit/f6bb4d9).
- Loading of maps, styles, routes etc from (XML) strings.
- Option to set [aspect fix mode](https://github.com/mapnik/mapnik/wiki/Aspect-Fix-Mode)

Installation
------------

This package can be installed with the go get command. The following steps need
to be performed in order. `go generate` will setup your environment and needs
to be run prior to installing the package:

    go get -d github.com/sgelb/go-mapnik
    go generate github.com/sgelb/go-mapnik
    go install github.com/sgelb/go-mapnik

This package requires [Mapnik](http://mapnik.org/) (`libmapnik-dev` on
Ubuntu/Debian, `mapnik --with-postgresql` in Homebrew). Make sure
`mapnik-config` is in your `PATH`.
