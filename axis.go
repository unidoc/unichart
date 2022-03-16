package chart

import (
	"github.com/unidoc/unichart/data"
	"github.com/unidoc/unichart/render"
)

// Axis is a chart feature detailing what values happen where.
type Axis interface {
	GetName() string
	SetName(name string)

	GetStyle() render.Style
	SetStyle(style render.Style)

	GetTicks() []Tick
	GenerateTicks(r render.Renderer, ra data.Range, vf data.ValueFormatter) []Tick

	// GenerateGridLines returns the gridlines for the axis.
	GetGridLines(ticks []Tick) []GridLine

	// Measure should return an absolute box for the axis.
	// This is used when auto-fitting the canvas to the background.
	Measure(r render.Renderer, canvasBox render.Box, ra data.Range, style render.Style, ticks []Tick) render.Box

	// Render renders the axis.
	Render(r render.Renderer, canvasBox render.Box, ra data.Range, style render.Style, ticks []Tick)
}
