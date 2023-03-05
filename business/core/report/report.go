// Package report provides an example of a core business API.
package report

import (
	"reflect"
	"sync"

	"github.com/deliveranceTechSolutions/erp/business/data/store/user"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Core manages the set of API's for user access.
type Core struct {
	log  *zap.SugaredLogger
	user user.Store
	dash Dashboard
}

// NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger, db *sqlx.DB) Core {
	return Core{
		log:  log,
		user: user.NewStore(log, db),
		dash: generateDashboard(),
	}
}

// Interface being implemented in chart.go
//
// type report interface {
// 	LoadData()
// 	CanHaveView() bool
// }

// Dashboard creates a map which stores new charts
type Dashboard struct {
	mu     sync.Mutex
	charts map[string]report
}

// generateDashboard creates a personal mutex for each user's dash
func generateDashboard() Dashboard {
	var mu sync.Mutex
	return Dashboard{
		mu:     mu,
		charts: make(map[string]report),
	}
}

// CreateNew is a chart constructor
func (d *Dashboard) CreateNew(name string, chartType interface{}) Chart {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	d.charts[name] = selectChart(chartType)

	return Chart{}
}

func selectChart[CT collection](chartType CT) CT {
	var chart CT	

	return chart
}

type collection interface {
	GanttChart 			|
	GanttTable 			|
	HistogramChart 		|
	ColumnChart 		|
	LineChart 			|
	BarChart 			|
	PieChart 			|
	ScatterXYChart 		|
	BubbleChart 		|
	StockChart 			|
	SurfaceChart 		|
	RadarChart 			|
	TreemapChart 		|
	SunburstChart 		|
	BoxandWhiskerChart 	|
	WaterfallChart 		|
	FunnelChart 		|
	ComboChart 			|
	MapChart
}
