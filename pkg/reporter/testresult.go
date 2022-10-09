package reporter

import "time"

type TestResult struct {
	Area     string
	Feature  string
	Suite    string
	File     string
	Total    int
	Passes   int
	Pending  int
	Failures int
	Skipped  int
	Uuid     string
	TestRun  time.Time
}
