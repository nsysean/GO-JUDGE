package models

type Problem struct {
	TimeLimit   int        `xml:"judging>testset>time-limit"`
	MemoryLimit int        `xml:"judging>testset>memory-limit"`
	Tests       []Test     `xml:"judging>testset>tests>test"`
	Groups      []Group    `xml:"judging>testset>groups>group"`
	RunCount    int        `xml:"judging>testset>test-count"`
	InputFile   string     `xml:"judging>testset>input-path-pattern"`
	OutputFile  string     `xml:"judging>testset>answer-path-pattern"`
	Interactor  Interactor `xml:"assets>interactor>source"`
}

type Test struct {
	Points       float32 `xml:"points,attr"`
	Group        string  `xml:"group,attr"`
	ScoredPoints float32
	Verdict      string
	Time         float32
	Memory       int
	Exit         int
}

type Group struct {
	Name           string  `xml:"name,attr"`
	Points         float32 `xml:"points,attr"`
	PointStore     float32
	PointsPolicy   string `xml:"points-policy,attr"`
	FeedbackPolicy string `xml:"feedback-policy,attr"`
	Dependencies   []dependency `xml:"dependencies>dependency"`
}

type dependency struct {
	Group string `xml:"group,attr"`
}

type Interactor struct {
	Source string `xml:"path,attr"`
}

type Verdict struct {
	Prefix  string
	Verdict string
}
