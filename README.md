unichart
========

unichart is a native Golang charting library. Its primary purpose is to
integrate with [UniPDF](https://github.com/unidoc/unipdf) and other products
in the [UniDoc](https://unidoc.io/) ecosystem in order to provide charting
capabilities.

# Installation

```bash
go get -u github.com/unidoc/unichart
```

# Examples

For usage and output samples, see the [examples](examples) directory.

# Supported chart types

### Line chart

![Sample chart output](examples/request_timing/preview.png)

### Bar chart

![Sample bar chart output](examples/bar_chart/preview.png)

### Stacked bar chart

![Sample stacked bar chart output](examples/stacked_bar/preview.png)

### Pie chart

![Sample pie chart output](examples/pie_chart/preview.png)

### Donut chart

![Sample donut chart output](examples/donut_chart/preview.png)

### Linear Progress Bar

![Sample linear progress bar](examples/progress_bar/preview_linear.png)

### Circular Progress Bar

![Sample circular progress bar](examples/progress_bar/preview_circular.png)

# License

This project is licensed under the [MIT license](https://opensource.org/licenses/MIT).
See [LICENSE](LICENSE) for more details.

A commercial license is required in order to use this project with [UniPDF](https://github.com/unidoc/unipdf).
There are two options:

- Metered License API keys.
  - Free ones can be obtained at https://cloud.unidoc.io.
- Offline Perpetual codes.
  - Can be purchased at https://unidoc.io/pricing.

# Credits

This is a modified fork of [wcharczuk/go-chart](https://github.com/wcharczuk/go-chart).
