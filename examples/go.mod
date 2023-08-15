module github.com/unidoc/unichart/examples

go 1.16

require (
	github.com/disintegration/imaging v1.6.2
	github.com/unidoc/unichart v0.2.0
	github.com/unidoc/unipdf/v3 v3.49.0
)

replace github.com/unidoc/unichart => ../

replace github.com/unidoc/unipdf/v3 => ../../unipdf-src
