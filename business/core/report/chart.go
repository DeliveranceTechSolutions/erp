package report

import (
	"errors"
	"time"
)

type report interface {
	LoadData() error
	Render() error
}

type Chart struct {
	Data 		any
	Title 		string
	X			Axis
	Y			Axis
	IsLoaded	bool
}

type Axis struct {
	Frequency 	int
	Magnitude	int
	Name		string
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

// ==================================================
// LineChart
type LineChart struct {
	Chart
}

func (lc *LineChart) LoadData() error {
	databaseData := []int{}
	lc.Data = databaseData

	return nil
}

func (lc *LineChart) Render() error {
	if !lc.IsLoaded {
		return errors.New("LineChart render error")
	}

	return nil
}

// ==================================================
// BarChart
type BarChart struct {
	Chart
}

func (bc *BarChart) LoadData() error {
	databaseData := []int{}
	bc.Data = databaseData

	return nil

}

func (bc *BarChart) Render() error {
	if !bc.IsLoaded {
		return errors.New("BarChart render error")
	}

	return nil
}

// ==================================================
// PieChart
type PieChart struct {
	Chart
}

func (pc *PieChart) LoadData() error {
	databaseData := []int{}
	pc.Data = databaseData

	return nil
}

func (pc *PieChart) Render() error {
	if !pc.IsLoaded {
		return errors.New("PieChart render error")
	}

	return nil
}

// ==================================================
// ScatterXYChart
type ScatterXYChart struct {
	Chart
}
func (sc *ScatterXYChart) LoadData() error {
	databaseData := []int{}
	sc.Data = databaseData

	return nil
}

func (sc *ScatterXYChart) Render() error {
	if !sc.IsLoaded {
		return errors.New("ScatterXYChart render error")
	}

	return nil
}

// ==================================================
// BubbleChart
type BubbleChart struct {
	Chart
}

func (bc *ScatterXYChart) LoadData() error {
	databaseData := []int{}
	sc.Data = databaseData

	return nil
}

func (sc *ScatterXYChart) Render() error {
	if !sc.IsLoaded {
		return errors.New("ScatterXYChart render error")
	}

	return nil
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
