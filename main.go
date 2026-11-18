package main

// China CO2 emissions 1970-2022: load, clean, explore, one linear regression, save.
// Standard library only. Run with: go run .   Test with: go test -v

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// =============================================================================
// PART 1: LOAD AND CLEAN THE DATA
// =============================================================================

// Record is one cleaned row: one year of China's CO2 emissions.
type Record struct {
	Year int
	CO2  float64 // million tonnes of CO2 (Mt CO2e)
}

// QualityReport lists what the data check found, so it can be printed and saved.
type QualityReport struct {
	RawRows        int
	EmptyRows      int
	MissingValues  int
	BadNumbers     int
	DuplicateYears int
	NegativeValues int
	MissingColumns []string // required columns not found in the header
	MissingYears   []int
	Columns        []string
}

// LoadRaw reads the CSV file and returns the header and all rows as text.
func LoadRaw(path string) ([]string, [][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1 // allow rows with a different number of fields (e.g. a blank last line)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	if len(rows) == 0 {
		return nil, nil, fmt.Errorf("file is empty")
	}
	return rows[0], rows[1:], nil
}

// columnIndex finds a column by name (ignoring upper/lower case and spaces).
func columnIndex(header []string, name string) int {
	for i, h := range header {
		if strings.EqualFold(strings.TrimSpace(h), name) {
			return i
		}
	}
	return -1
}

// Clean turns the raw text rows into Records and reports every data-quality problem it finds.
func Clean(header []string, rows [][]string) ([]Record, QualityReport) {
	report := QualityReport{RawRows: len(rows), Columns: header}

	yearCol := columnIndex(header, "Year")
	valueCol := columnIndex(header, "CO2_Mt_CO2e")
	// a missing column would make row[-1] crash, so stop and report it instead
	for name, col := range map[string]int{"Year": yearCol, "CO2_Mt_CO2e": valueCol} {
		if col < 0 {
			report.MissingColumns = append(report.MissingColumns, name)
		}
	}
	if len(report.MissingColumns) > 0 {
		sort.Strings(report.MissingColumns)
		return nil, report
	}

	seen := map[int]bool{}
	var records []Record

	for _, row := range rows {
		// skip blank lines
		if len(row) == 0 || (len(row) == 1 && strings.TrimSpace(row[0]) == "") {
			report.EmptyRows++
			continue
		}
		// a short row would crash on row[yearCol] or row[valueCol]
		if len(row) <= yearCol || len(row) <= valueCol {
			report.MissingValues++
			continue
		}
		yearText := strings.TrimSpace(row[yearCol])
		valueText := strings.TrimSpace(row[valueCol])
		if yearText == "" || valueText == "" {
			report.MissingValues++
			continue
		}
		year, err1 := strconv.Atoi(yearText)
		value, err2 := strconv.ParseFloat(valueText, 64)
		// ParseFloat accepts "NaN" and "Inf" without an error, so reject them here too
		// a year far outside the data range would flood MissingYears with fake gaps
		if err1 != nil || err2 != nil || math.IsNaN(value) || math.IsInf(value, 0) || year < 1900 || year > 2100 {
			report.BadNumbers++
			continue
		}
		if seen[year] { // checked before negatives so a repeated year is always reported as a duplicate
			report.DuplicateYears++
			continue
		}
		if value < 0 {
			report.NegativeValues++
			continue
		}
		seen[year] = true
		records = append(records, Record{Year: year, CO2: value})
	}

	// sort by year and look for gaps
	sort.Slice(records, func(i, j int) bool { return records[i].Year < records[j].Year })
	for i := 1; i < len(records); i++ {
		for y := records[i-1].Year + 1; y < records[i].Year; y++ {
			report.MissingYears = append(report.MissingYears, y)
		}
	}
	return records, report
}

// Print shows the data-quality report in the terminal.
func (r QualityReport) Print(records []Record) {
	fmt.Println("=== Data quality check ===")
	fmt.Printf("Columns: %v\n", r.Columns)
	if len(r.MissingColumns) > 0 {
		fmt.Printf("Missing required columns: %v\n", r.MissingColumns)
	}
	fmt.Printf("Raw rows: %d | clean rows: %d\n", r.RawRows, len(records))
	fmt.Printf("Empty rows removed: %d\n", r.EmptyRows)
	fmt.Printf("Missing values: %d | bad numbers: %d | negative values: %d | duplicate years: %d\n",
		r.MissingValues, r.BadNumbers, r.NegativeValues, r.DuplicateYears)
	fmt.Printf("Missing years in the series: %v\n", r.MissingYears)
	if len(records) > 0 {
		fmt.Printf("Years covered: %d to %d\n", records[0].Year, records[len(records)-1].Year)
	}
	fmt.Println()
}

// =============================================================================
// PART 2: EXPLORE THE DATA AND COMPARE PERIODS
// =============================================================================

// Summary holds simple descriptive statistics.
type Summary struct {
	Count             int
	Mean, Median, Std float64
	Min, Max          float64
	MinYear, MaxYear  int
}

// Describe calculates the summary statistics of the CO2 values.
func Describe(records []Record) Summary {
	s := Summary{Count: len(records), Min: math.Inf(1), Max: math.Inf(-1)}
	values := make([]float64, len(records))
	for i, r := range records {
		values[i] = r.CO2
		s.Mean += r.CO2
		if r.CO2 < s.Min {
			s.Min, s.MinYear = r.CO2, r.Year
		}
		if r.CO2 > s.Max {
			s.Max, s.MaxYear = r.CO2, r.Year
		}
	}
	s.Mean /= float64(len(values))
	for _, v := range values {
		s.Std += (v - s.Mean) * (v - s.Mean)
	}
	s.Std = math.Sqrt(s.Std / float64(len(values)-1)) // sample standard deviation

	sort.Float64s(values)
	n := len(values)
	if n%2 == 1 {
		s.Median = values[n/2]
	} else {
		s.Median = (values[n/2-1] + values[n/2]) / 2
	}
	return s
}

// YearRow is one year with the extra columns we calculate.
type YearRow struct {
	Year      int
	CO2       float64
	Change    float64 // change from the previous year (Mt)
	GrowthPct float64 // % change from the previous year
	Period    string
}

// PeriodOf groups each year into a period of China's economic history.
func PeriodOf(year int) string {
	switch {
	case year <= 1977:
		return "1970-1977 Pre-reform"
	case year <= 2000:
		return "1978-2000 Reform and opening-up"
	case year <= 2013:
		return "2001-2013 WTO and industrial boom"
	default:
		return "2014-2022 Slowdown and plateau"
	}
}

// AddFeatures calculates the yearly change, growth % and period for every year.
func AddFeatures(records []Record) []YearRow {
	rows := make([]YearRow, len(records))
	for i, r := range records {
		row := YearRow{Year: r.Year, CO2: r.CO2, Period: PeriodOf(r.Year)}
		if i > 0 {
			prev := records[i-1].CO2
			row.Change = r.CO2 - prev
			row.GrowthPct = (r.CO2 - prev) / prev * 100
		} else {
			row.Change, row.GrowthPct = math.NaN(), math.NaN() // no previous year for 1970
		}
		rows[i] = row
	}
	return rows
}

// PeriodStats summarises one period.
type PeriodStats struct {
	Name          string
	Years         int
	StartCO2      float64
	EndCO2        float64
	MeanCO2       float64
	AvgGrowthPct  float64 // average yearly growth (compound annual growth rate)
	TotalIncrease float64
}

// GroupByPeriod summarises the rows period by period, keeping their order.
func GroupByPeriod(rows []YearRow) []PeriodStats {
	var order []string
	groups := map[string][]YearRow{}
	for _, r := range rows {
		if _, ok := groups[r.Period]; !ok {
			order = append(order, r.Period)
		}
		groups[r.Period] = append(groups[r.Period], r)
	}

	var result []PeriodStats
	for i, name := range order {
		g := groups[name]
		// start from the year before the period (if there is one), so growth includes the first step
		start := g[0].CO2
		steps := len(g) - 1
		if i > 0 {
			prev := groups[order[i-1]]
			start = prev[len(prev)-1].CO2
			steps = len(g)
		}
		end := g[len(g)-1].CO2
		mean := 0.0
		for _, r := range g {
			mean += r.CO2
		}
		mean /= float64(len(g))
		cagr := (math.Pow(end/start, 1/float64(steps)) - 1) * 100
		result = append(result, PeriodStats{Name: name, Years: len(g), StartCO2: start, EndCO2: end,
			MeanCO2: mean, AvgGrowthPct: cagr, TotalIncrease: end - start})
	}
	return result
}

// =============================================================================
// PART 3: THE REGRESSION MODEL
// =============================================================================

//
// The model is a straight line:   CO2 = intercept + slope × t
// where t = decades since 1970 (so 1970 -> 0, 1980 -> 1, 2022 -> 5.2).
//
// Time-based split: train on 1970-2012, test on 2013-2022.
// We never shuffle the years, because the model must predict the future from the past.
//
// For comparison we also score a "naive" guess: next year = this year (no regression).

const testFromYear = 2013 // first year of the test set

// decades turns a year into "decades since 1970". Small numbers are easier to read.
func decades(year int) float64 { return float64(year-1970) / 10 }

// Line is a straight line: y = Intercept + Slope × x.
type Line struct{ Intercept, Slope float64 }

// Predict gives the y value for one x value.
func (l Line) Predict(x float64) float64 { return l.Intercept + l.Slope*x }

// FitLine finds the best straight line through the points (ordinary least squares).
//
//	slope     = sum((x - meanX) × (y - meanY)) / sum((x - meanX)²)
//	intercept = meanY - slope × meanX
func FitLine(x, y []float64) (Line, error) {
	if len(x) != len(y) || len(x) < 2 {
		return Line{}, errors.New("need at least 2 points, with the same number of x and y values")
	}
	n := float64(len(x))
	meanX, meanY := 0.0, 0.0
	for i := range x {
		meanX += x[i]
		meanY += y[i]
	}
	meanX /= n
	meanY /= n

	cross, spread := 0.0, 0.0
	for i := range x {
		cross += (x[i] - meanX) * (y[i] - meanY)
		spread += (x[i] - meanX) * (x[i] - meanX)
	}
	if spread == 0 {
		return Line{}, errors.New("all x values are the same, so no line can be fitted")
	}
	slope := cross / spread
	return Line{Intercept: meanY - slope*meanX, Slope: slope}, nil
}

// Metrics measures how close the predictions are to the actual values.
type Metrics struct {
	MAE, RMSE, MAPE, R2 float64
}

// Evaluate calculates MAE, RMSE, MAPE (%) and R².
func Evaluate(actual, predicted []float64) Metrics {
	n := float64(len(actual))
	mean := 0.0
	for _, a := range actual {
		mean += a
	}
	mean /= n

	var m Metrics
	ssRes, ssTot := 0.0, 0.0
	for i := range actual {
		e := actual[i] - predicted[i]
		m.MAE += math.Abs(e)
		m.RMSE += e * e
		m.MAPE += math.Abs(e / actual[i])
		ssRes += e * e
		ssTot += (actual[i] - mean) * (actual[i] - mean)
	}
	m.MAE /= n
	m.RMSE = math.Sqrt(m.RMSE / n)
	m.MAPE = m.MAPE / n * 100
	m.R2 = 1 - ssRes/ssTot
	return m
}

// Result stores one model's predictions and scores.
type Result struct {
	Name, Formula string
	Predictions   map[int]float64 // year -> predicted CO2
	Train, Test   Metrics
}

// TrainAndTest fits the line on the training years and scores it and the naive guess.
func TrainAndTest(records []Record) (Line, []Result, error) {
	var x, y []float64
	for _, r := range records {
		if r.Year < testFromYear { // training rows only
			x = append(x, decades(r.Year))
			y = append(y, r.CO2)
		}
	}
	line, err := FitLine(x, y)
	if err != nil {
		return Line{}, nil, err
	}

	linear := map[int]float64{}
	naive := map[int]float64{}
	for i, r := range records {
		linear[r.Year] = line.Predict(decades(r.Year))
		if i > 0 {
			naive[r.Year] = records[i-1].CO2 // last year's value
		}
	}
	results := []Result{
		score("Naive (last year)", "CO2(t) = CO2(t-1)", naive, records),
		score("Linear trend", "CO2 = b0 + b1*t", linear, records),
	}
	return line, results, nil
}

// score compares the predictions with the actual values, for the train and test years.
func score(name, formula string, predictions map[int]float64, records []Record) Result {
	var trainA, trainP, testA, testP []float64
	for _, r := range records {
		p, ok := predictions[r.Year]
		if !ok {
			continue
		}
		if r.Year < testFromYear {
			trainA, trainP = append(trainA, r.CO2), append(trainP, p)
		} else {
			testA, testP = append(testA, r.CO2), append(testP, p)
		}
	}
	return Result{Name: name, Formula: formula, Predictions: predictions,
		Train: Evaluate(trainA, trainP), Test: Evaluate(testA, testP)}
}

// =============================================================================
// PART 4: MAIN PROGRAM (runs the six steps and saves the CSV files)
// =============================================================================

//   1. load  2. clean  3. explore  4. compare periods  5. regression model  6. save CSVs
//
// Run with:  go run .

const (
	rawFile      = "data/raw/NGDS_CO2_China.csv"
	processedDir = "data/processed"
)

func main() {
	// 1. Load
	header, rows, err := LoadRaw(rawFile)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Clean and check quality
	records, report := Clean(header, rows)
	report.Print(records)
	if len(records) < 2 {
		fmt.Println("Not enough clean rows to analyse. Check the data quality report above.")
		os.Exit(1)
	}

	// 3. Explore
	s := Describe(records)
	fmt.Println("=== Summary statistics (Mt CO2) ===")
	fmt.Printf("Count %d | mean %.0f | median %.0f | std %.0f\n", s.Count, s.Mean, s.Median, s.Std)
	fmt.Printf("Lowest %.0f (%d) | highest %.0f (%d) | growth %.1fx\n\n",
		s.Min, s.MinYear, s.Max, s.MaxYear, records[len(records)-1].CO2/records[0].CO2)

	yearRows := AddFeatures(records)
	biggestRise, biggestFall := yearRows[1], yearRows[1]
	for _, r := range yearRows[1:] {
		if r.GrowthPct > biggestRise.GrowthPct {
			biggestRise = r
		}
		if r.GrowthPct < biggestFall.GrowthPct {
			biggestFall = r
		}
	}
	fmt.Printf("Fastest yearly growth: %d (%+.1f%%) | biggest fall: %d (%+.1f%%)\n\n",
		biggestRise.Year, biggestRise.GrowthPct, biggestFall.Year, biggestFall.GrowthPct)

	// 4. Differences across periods
	periods := GroupByPeriod(yearRows)
	fmt.Println("=== Periods ===")
	for _, g := range periods {
		fmt.Printf("%-36s %2d yrs | %6.0f -> %6.0f Mt | avg growth %5.1f%%/yr\n",
			g.Name, g.Years, g.StartCO2, g.EndCO2, g.AvgGrowthPct)
	}

	// 5. Regression model: train on the early years, test on the last ten
	line, results, err := TrainAndTest(records)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n=== Regression model (train 1970-%d, test %d-2022) ===\n", testFromYear-1, testFromYear)
	for _, r := range results {
		fmt.Printf("%-20s train MAPE %5.1f%% | test MAE %6.0f  RMSE %6.0f  MAPE %5.1f%%  R2 %6.2f\n",
			r.Name, r.Train.MAPE, r.Test.MAE, r.Test.RMSE, r.Test.MAPE, r.Test.R2)
	}
	fmt.Printf("  %-20s CO2 = b0 + b1*t  coefficients %v\n", "Linear trend",
		[]float64{round(line.Intercept, 4), round(line.Slope, 4)})
	fmt.Printf("The line adds %.0f Mt per decade, about %.0f Mt every year (t = decades since 1970).\n",
		line.Slope, line.Slope/10)

	// 6. Save processed data
	must(os.MkdirAll(processedDir, 0o755))
	saveClean(yearRows)
	saveGroups(periods)
	saveMetrics(results)
	savePredictions(records, results)
	fmt.Println("\nSaved processed CSV files to", processedDir)
}

// ---------- CSV helpers ----------

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func writeCSV(name string, rows [][]string) {
	f, err := os.Create(processedDir + "/" + name)
	must(err)
	defer f.Close()
	w := csv.NewWriter(f)
	must(w.WriteAll(rows))
}

// num formats a number; missing values (NaN) become an empty cell.
func num(v float64, decimals int) string {
	if math.IsNaN(v) {
		return ""
	}
	return strconv.FormatFloat(v, 'f', decimals, 64)
}

func round(v float64, decimals int) float64 {
	p := math.Pow(10, float64(decimals))
	return math.Round(v*p) / p
}

func splitOf(year int) string {
	if year >= testFromYear {
		return "test"
	}
	return "train"
}

func saveClean(rows []YearRow) {
	out := [][]string{{"year", "co2_mt", "change_mt", "growth_pct", "period", "split"}}
	for _, r := range rows {
		out = append(out, []string{strconv.Itoa(r.Year), num(r.CO2, 2), num(r.Change, 2), num(r.GrowthPct, 2),
			r.Period, splitOf(r.Year)})
	}
	writeCSV("co2_china_clean.csv", out)
}

func saveGroups(periods []PeriodStats) {
	out := [][]string{{"period", "years", "start_co2_mt", "end_co2_mt", "mean_co2_mt", "avg_growth_pct_per_year", "total_increase_mt"}}
	for _, g := range periods {
		out = append(out, []string{g.Name, strconv.Itoa(g.Years), num(g.StartCO2, 1), num(g.EndCO2, 1),
			num(g.MeanCO2, 1), num(g.AvgGrowthPct, 2), num(g.TotalIncrease, 1)})
	}
	writeCSV("period_summary.csv", out)
}

func saveMetrics(results []Result) {
	out := [][]string{{"model", "formula", "train_mae", "train_mape", "test_mae", "test_rmse", "test_mape", "test_r2"}}
	for _, r := range results {
		out = append(out, []string{r.Name, r.Formula, num(r.Train.MAE, 1), num(r.Train.MAPE, 2),
			num(r.Test.MAE, 1), num(r.Test.RMSE, 1), num(r.Test.MAPE, 2), num(r.Test.R2, 3)})
	}
	writeCSV("model_metrics.csv", out)
}

func savePredictions(records []Record, results []Result) {
	head := []string{"year", "actual_co2_mt", "split"}
	for _, r := range results {
		head = append(head, r.Name)
	}
	out := [][]string{head}
	for _, rec := range records {
		row := []string{strconv.Itoa(rec.Year), num(rec.CO2, 2), splitOf(rec.Year)}
		for _, r := range results {
			if p, ok := r.Predictions[rec.Year]; ok {
				row = append(row, num(p, 2))
			} else {
				row = append(row, "") // no prediction for the first year
			}
		}
		out = append(out, row)
	}
	writeCSV("model_predictions.csv", out)
}
