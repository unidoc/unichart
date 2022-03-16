package chart

import "github.com/unidoc/unichart/render"

const (
	// DefaultChartHeight is the default chart height.
	DefaultChartHeight = 400

	// DefaultChartWidth is the default chart width.
	DefaultChartWidth = 1024

	// DefaultCanvasStrokeWidth is the default stroke on the chart canvas.
	DefaultCanvasStrokeWidth = 0.0

	// DefaultBackgroundStrokeWidth is the default chart background stroke.
	DefaultBackgroundStrokeWidth = 0.0

	// DefaultSeriesLineWidth is the default line width.
	DefaultSeriesLineWidth = 1.0

	// DefaultDPI is the default dots per inch for the chart.
	DefaultDPI = 72.0

	// DefaultTitleFontSize is the default title font size.
	DefaultTitleFontSize = 18.0

	// DefaultTitleTop is the default distance from the top of the chart to put the title.
	DefaultTitleTop = 10

	// DefaultAxisLineWidth is the line width of the axis lines.
	DefaultAxisLineWidth = 1.0

	// DefaultAxisFontSize is the font size of the axis labels.
	DefaultAxisFontSize = 10.0

	// DefaultYAxisMargin is the default distance from the right
	// of the canvas to the y axis labels.
	DefaultYAxisMargin = 10

	// DefaultXAxisMargin is the default distance from bottom of
	// the canvas to the x axis labels.
	DefaultXAxisMargin = 10

	//DefaultVerticalTickHeight is half the margin.
	DefaultVerticalTickHeight = DefaultXAxisMargin >> 1

	//DefaultHorizontalTickWidth is half the margin.
	DefaultHorizontalTickWidth = DefaultYAxisMargin >> 1

	// DefaultTickCount is the default number of ticks to show
	DefaultTickCount = 10

	// DefaultTickCountSanityCheck is a hard limit on number of
	// ticks to prevent infinite loops.
	DefaultTickCountSanityCheck = 1 << 10 //1024

	// DefaultMinimumTickHorizontalSpacing is the minimum distance
	// between horizontal ticks.
	DefaultMinimumTickHorizontalSpacing = 20

	// DefaultMinimumTickVerticalSpacing is the minimum distance
	// between vertical ticks.
	DefaultMinimumTickVerticalSpacing = 20

	// DefaultBarSpacing is the default pixel spacing between bars.
	DefaultBarSpacing = 100

	// DefaultBarWidth is the default pixel width of bars in a bar chart.
	DefaultBarWidth = 50
)

var (
	// DashArrayDots is a dash array that
	// represents '....' style stroke dashes.
	DashArrayDots = []int{1, 1}

	// DashArrayDashesSmall is a dash array that
	// represents '- - -' style stroke dashes.
	DashArrayDashesSmall = []int{3, 3}

	// DashArrayDashesMedium is a dash array that
	// represents '-- -- --' style stroke dashes.
	DashArrayDashesMedium = []int{5, 5}

	// DashArrayDashesLarge is a dash array that
	// represents '----- ----- -----' style stroke dashes.
	DashArrayDashesLarge = []int{10, 10}
)

var (
	// DefaultBackgroundPadding is the default canvas padding config.
	DefaultBackgroundPadding = render.Box{
		Top:    5,
		Left:   5,
		Right:  5,
		Bottom: 5,
	}
)
