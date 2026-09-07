// internal/engine/types.go
package engine

import (
	"time"
)

type Level int

const (
	LevelA   Level = iota
	LevelAA
	LevelAAA
)

func (l Level) String() string {
	switch l {
	case LevelA:
		return "A"
	case LevelAA:
		return "AA"
	case LevelAAA:
		return "AAA"
	default:
		return "Unknown"
	}
}

type Testability int

const (
	Auto   Testability = iota
	Semi
	Manual
)

func (t Testability) String() string {
	switch t {
	case Auto:
		return "AUTO"
	case Semi:
		return "SEMI"
	case Manual:
		return "MANUAL"
	default:
		return "Unknown"
	}
}

type Outcome int

const (
	Supports          Outcome = iota
	PartiallySupports
	DoesNotSupport
	NotApplicable
	NotEvaluated
	OutcomeError
	OutcomeTimeout
)

func (o Outcome) String() string {
	switch o {
	case Supports:
		return "Supports"
	case PartiallySupports:
		return "Partially Supports"
	case DoesNotSupport:
		return "Does Not Support"
	case NotApplicable:
		return "Not Applicable"
	case NotEvaluated:
		return "Not Evaluated"
	case OutcomeError:
		return "Error"
	case OutcomeTimeout:
		return "Timeout"
	default:
		return "Unknown"
	}
}

type Evidence struct {
	Command  string            `json:"command"`
	Env      map[string]string `json:"env,omitempty"`
	Stdout   string            `json:"stdout"`
	Stderr   string            `json:"stderr"`
	ExitCode int               `json:"exit_code"`
	Duration time.Duration     `json:"duration_ms"`
	Note     string            `json:"note,omitempty"`
}

type Result struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Domain      string        `json:"domain"`
	Level       Level         `json:"level"`
	Testability Testability   `json:"testability"`
	Outcome     Outcome       `json:"outcome"`
	Evidence    []Evidence    `json:"evidence,omitempty"`
	Remarks     string        `json:"remarks,omitempty"`
	Duration    time.Duration `json:"duration_ms"`
	NeedsReview bool          `json:"needs_review,omitempty"`
	SpecVersion string        `json:"spec_version"`
}

func (r *Result) IsFailure() bool {
	switch r.Outcome {
	case DoesNotSupport, PartiallySupports, OutcomeError, OutcomeTimeout:
		return true
	default:
		return false
	}
}

type ProductInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Version string `json:"version"`
}

type TestEnvironment struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Terminal string `json:"terminal,omitempty"`
	Shell    string `json:"shell,omitempty"`
	SuiteVer string `json:"suite_version"`
}

type DomainSummary struct {
	Domain   string `json:"domain"`
	Total    int    `json:"total"`
	Pass     int    `json:"pass"`
	Partial  int    `json:"partial"`
	Fail     int    `json:"fail"`
	NA       int    `json:"na"`
	NotEval  int    `json:"not_evaluated"`
}

type ConformanceReport struct {
	CLIACSVersion   string            `json:"cli_acs_version"`
	SuiteVersion    string            `json:"suite_version"`
	Product         ProductInfo       `json:"product"`
	ReportDate      time.Time         `json:"report_date"`
	Environment     TestEnvironment   `json:"environment"`
	Results         []Result          `json:"results"`
	DomainSummaries []DomainSummary   `json:"domain_summaries"`
	OverallLevel    string            `json:"overall_level"`
	Threshold       string            `json:"threshold"`
}

func (r *ConformanceReport) ComputeSummaries() {
	domainMap := make(map[string]*DomainSummary)
	for _, res := range r.Results {
		s, ok := domainMap[res.Domain]
		if !ok {
			s = &DomainSummary{Domain: res.Domain}
			domainMap[res.Domain] = s
		}
		s.Total++
		switch res.Outcome {
		case Supports:
			s.Pass++
		case PartiallySupports:
			s.Partial++
		case DoesNotSupport, OutcomeError, OutcomeTimeout:
			s.Fail++
		case NotApplicable:
			s.NA++
		case NotEvaluated:
			s.NotEval++
		}
	}
	r.DomainSummaries = make([]DomainSummary, 0, len(domainMap))
	for _, s := range domainMap {
		r.DomainSummaries = append(r.DomainSummaries, *s)
	}
}
