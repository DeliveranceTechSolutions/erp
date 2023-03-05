package report

import "time"

type report interface {
	LoadData()
	CanHaveView() bool
}

type Chart struct {
	Data []int
	Name string
}

func (c *Chart) LoadData() {
	data := []int{}

	c.Data = data
}

func (c *Chart) CanHaveView() bool {
	return true
}

type GanttChart struct {
	Table  []*GanttTable
	loaded bool
}

type GanttTable struct {
	ItemName  string
	Duration  string
	StartDate time.Time
	EndDate   time.Time
}

type HistogramChart struct {
}

type ColumnChart struct {
}
type LineChart struct {
}
type BarChart struct {
}
type PieChart struct {
}
type ScatterXYChart struct {
}
type BubbleChart struct {
}
type StockChart struct {
}
type SurfaceChart struct {
}
type RadarChart struct {
}
type TreemapChart struct {
}
type SunburstChart struct {
}

type BoxandWhiskerChart struct {
}

type WaterfallChart struct {
}

type FunnelChart struct {
}

type ComboChart struct {
}

type MapChart struct {
}
