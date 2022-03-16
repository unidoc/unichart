package chart

import "github.com/unidoc/unichart/render"

const (
	// defaultChartHeight is the default chart height.
	defaultChartHeight = 400

	// defaultChartWidth is the default chart width.
	defaultChartWidth = 1024

	// defaultCanvasStrokeWidth is the default stroke on the chart canvas.
	defaultCanvasStrokeWidth = 0.0

	// defaultBackgroundStrokeWidth is the default chart background stroke.
	defaultBackgroundStrokeWidth = 0.0

	// defaultSeriesLineWidth is the default line width.
	defaultSeriesLineWidth = 1.0

	// defaultDPI is the default dots per inch for the chart.
	defaultDPI = 72.0

	// defaultTitleFontSize is the default title font size.
	defaultTitleFontSize = 18.0

	// defaultTitleTop is the default distance from the top of the chart to put the title.
	defaultTitleTop = 10

	// defaultAxisLineWidth is the line width of the axis lines.
	defaultAxisLineWidth = 1.0

	// defaultAxisFontSize is the font size of the axis labels.
	defaultAxisFontSize = 10.0

	// defaultXAxisMargin is the default distance from bottom of
	// the canvas to the x axis labels.
	defaultXAxisMargin = 10

	// defaultYAxisMargin is the default distance from the right
	// of the canvas to the y axis labels.
	defaultYAxisMargin = 10

	// defaultVerticalTickHeight is half the margin.
	defaultVerticalTickHeight = defaultXAxisMargin >> 1

	// defaultHorizontalTickWidth is half the margin.
	defaultHorizontalTickWidth = defaultYAxisMargin >> 1

	// defaultTickCountSanityCheck is a hard limit on number of
	// ticks to prevent infinite loops.
	defaultTickCountSanityCheck = 1 << 10 //1024

	// defaultMinimumTickHorizontalSpacing is the minimum distance
	// between horizontal ticks.
	defaultMinimumTickHorizontalSpacing = 20

	// defaultMinimumTickVerticalSpacing is the minimum distance
	// between vertical ticks.
	defaultMinimumTickVerticalSpacing = 20

	// defaultBarSpacing is the default pixel spacing between bars.
	defaultBarSpacing = 100

	// defaultBarWidth is the default pixel width of bars in a bar chart.
	defaultBarWidth = 50
)

var (
	// defaultBackgroundPadding is the default canvas padding config.
	defaultBackgroundPadding = render.Box{
		Top:    5,
		Left:   5,
		Right:  5,
		Bottom: 5,
	}
)
