package fixer


type TimeSeriesResp struct {
	Success 	bool 	`json:"success"`
	Timeseries 	bool 	`json:"timeseries"`
	StartDate 	string 	`json:"start_date"`
	EndDate 	string 	`json:"end_date"`
	Base 		string 	`json:"base"`
	Rates 		RatesByDate	`json:"rates"`
}

type RatesByDate map[string]Rates