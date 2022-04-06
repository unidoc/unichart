package dataset

import (
	"github.com/unidoc/unichart/dataset/sequence"
	"github.com/unidoc/unichart/render"
)

// YAxisType is a type of y-axis. It can either be primary or secondary.
type YAxisType int

const (
	// YAxisPrimary is the primary axis.
	YAxisPrimary YAxisType = 0

	// YAxisSecondary is the secondary axis.
	YAxisSecondary YAxisType = 1
)

// Series is an alias to Renderable.
type Series interface {
	GetName() string
	GetYAxis() YAxisType
	GetStyle() render.Style
	Validate() error
	Render(r render.Renderer, canvasBox render.Box, xrange, yrange sequence.Range, s render.Style)
}
