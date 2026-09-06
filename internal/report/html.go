// internal/report/html.go
package report

import (
	"html"
	"html/template"
	"io"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
)

// HTMLReporter renders conformance reports as standalone HTML.
type HTMLReporter struct{}

// Name returns the reporter identifier.
func (r *HTMLReporter) Name() string {
	return "html"
}

// FileExtension returns the recommended file extension.
func (r *HTMLReporter) FileExtension() string {
	return ".html"
}

// Render writes the report as standalone HTML with embedded CSS.
func (r *HTMLReporter) Render(report *engine.ConformanceReport, w io.Writer) error {
	tmpl, err := template.New("html").Funcs(template.FuncMap{
		"outcomeClass": r.outcomeClass,
		"outcomeSymbol": r.outcomeSymbol,
		"escapeHTML": html.EscapeString,
		"formatDuration": func(d interface{}) string {
			return d.(string)
		},
		"add": func(a, b int) int {
			return a + b
		},
	}).Parse(htmlTemplate)
	if err != nil {
		return err
	}

	// Group results by domain for the template
	domainResults := make(map[string][]engine.Result)
	for _, result := range report.Results {
		domainResults[result.Domain] = append(domainResults[result.Domain], result)
	}

	data := struct {
		Report        *engine.ConformanceReport
		DomainResults map[string][]engine.Result
	}{
		Report:        report,
		DomainResults: domainResults,
	}

	return tmpl.Execute(w, data)
}

// outcomeClass returns a CSS class name for the outcome.
func (r *HTMLReporter) outcomeClass(outcome engine.Outcome) string {
	switch outcome {
	case engine.Supports:
		return "pass"
	case engine.PartiallySupports:
		return "partial"
	case engine.DoesNotSupport, engine.OutcomeError, engine.OutcomeTimeout:
		return "fail"
	case engine.NotApplicable:
		return "na"
	case engine.NotEvaluated:
		return "not-eval"
	default:
		return "unknown"
	}
}

// outcomeSymbol returns a visual symbol for the outcome.
func (r *HTMLReporter) outcomeSymbol(outcome engine.Outcome) string {
	switch outcome {
	case engine.Supports:
		return "✓"
	case engine.PartiallySupports:
		return "⚠"
	case engine.DoesNotSupport, engine.OutcomeError, engine.OutcomeTimeout:
		return "✗"
	case engine.NotApplicable:
		return "N/A"
	case engine.NotEvaluated:
		return "⊘"
	default:
		return "?"
	}
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CLI Accessibility Conformance Report - {{.Report.Product.Name}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            background: #f5f5f5;
            padding: 20px;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            padding: 40px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
            border-radius: 8px;
        }

        h1 {
            color: #2c3e50;
            margin-bottom: 30px;
            padding-bottom: 15px;
            border-bottom: 3px solid #3498db;
            font-size: 2em;
        }

        h2 {
            color: #34495e;
            margin-top: 40px;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 2px solid #ecf0f1;
            font-size: 1.5em;
        }

        h3 {
            color: #7f8c8d;
            margin-top: 30px;
            margin-bottom: 15px;
            font-size: 1.25em;
        }

        h4 {
            color: #34495e;
            margin-top: 20px;
            margin-bottom: 10px;
            font-size: 1.1em;
        }

        table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 20px;
        }

        th, td {
            padding: 12px;
            text-align: left;
            border: 1px solid #ddd;
        }

        th {
            background-color: #3498db;
            color: white;
            font-weight: 600;
        }

        tr:nth-child(even) {
            background-color: #f8f9fa;
        }

        .info-table th {
            background-color: #95a5a6;
            width: 30%;
        }

        .summary-table th {
            background-color: #3498db;
        }

        .result-card {
            margin-bottom: 20px;
            padding: 15px;
            border: 1px solid #ddd;
            border-radius: 5px;
            background: #fafafa;
        }

        .result-header {
            display: flex;
            align-items: center;
            margin-bottom: 10px;
        }

        .outcome-badge {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 4px;
            font-weight: 600;
            font-size: 0.9em;
            margin-right: 10px;
        }

        .outcome-badge.pass {
            background-color: #27ae60;
            color: white;
        }

        .outcome-badge.partial {
            background-color: #f39c12;
            color: white;
        }

        .outcome-badge.fail {
            background-color: #e74c3c;
            color: white;
        }

        .outcome-badge.na {
            background-color: #95a5a6;
            color: white;
        }

        .outcome-badge.not-eval {
            background-color: #bdc3c7;
            color: white;
        }

        details {
            margin-top: 10px;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            background: white;
        }

        summary {
            cursor: pointer;
            font-weight: 600;
            padding: 5px;
            user-select: none;
        }

        summary:hover {
            background-color: #f0f0f0;
        }

        .evidence-item {
            margin-top: 10px;
            padding: 10px;
            background: #f8f9fa;
            border-left: 3px solid #3498db;
        }

        pre {
            background: #2c3e50;
            color: #ecf0f1;
            padding: 10px;
            border-radius: 4px;
            overflow-x: auto;
            font-size: 0.9em;
        }

        code {
            background: #e8e8e8;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: "Courier New", monospace;
            font-size: 0.9em;
        }

        .dashboard {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin: 20px 0;
        }

        .stat-card {
            padding: 20px;
            border-radius: 8px;
            text-align: center;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }

        .stat-card h3 {
            margin: 0 0 10px 0;
            font-size: 2em;
        }

        .stat-card p {
            margin: 0;
            color: #7f8c8d;
            font-size: 0.9em;
        }

        .stat-card.pass { background-color: #d5f4e6; border-left: 4px solid #27ae60; }
        .stat-card.partial { background-color: #fef5e7; border-left: 4px solid #f39c12; }
        .stat-card.fail { background-color: #fadbd8; border-left: 4px solid #e74c3c; }
        .stat-card.total { background-color: #e8f4f8; border-left: 4px solid #3498db; }

        .needs-review {
            background-color: #fff3cd;
            border: 1px solid #ffc107;
            padding: 8px;
            border-radius: 4px;
            margin-top: 10px;
            font-weight: 600;
        }

        @media print {
            body { background: white; }
            .container { box-shadow: none; }
            details { border: none; }
            details[open] summary { display: none; }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>CLI Accessibility Conformance Report</h1>

        <h2>Product Information</h2>
        <table class="info-table">
            <tr><th>Product Name</th><td>{{.Report.Product.Name}}</td></tr>
            <tr><th>Product Version</th><td>{{.Report.Product.Version}}</td></tr>
            {{if .Report.Product.Path}}<tr><th>Product Path</th><td>{{.Report.Product.Path}}</td></tr>{{end}}
            <tr><th>CLI-ACS Version</th><td>{{.Report.CLIACSVersion}}</td></tr>
            <tr><th>Suite Version</th><td>{{.Report.SuiteVersion}}</td></tr>
            <tr><th>Report Date</th><td>{{.Report.ReportDate.Format "2006-01-02 15:04:05 MST"}}</td></tr>
            <tr><th>Overall Level</th><td>{{.Report.OverallLevel}}</td></tr>
            <tr><th>Threshold</th><td>{{.Report.Threshold}}</td></tr>
        </table>

        <h2>Conformance Summary</h2>
        {{if .Report.DomainSummaries}}
        {{$totalPass := 0}}
        {{$totalPartial := 0}}
        {{$totalFail := 0}}
        {{$totalTests := 0}}
        {{range .Report.DomainSummaries}}
            {{$totalPass = add $totalPass .Pass}}
            {{$totalPartial = add $totalPartial .Partial}}
            {{$totalFail = add $totalFail .Fail}}
            {{$totalTests = add $totalTests .Total}}
        {{end}}

        <div class="dashboard">
            <div class="stat-card total">
                <h3>{{$totalTests}}</h3>
                <p>Total Tests</p>
            </div>
            <div class="stat-card pass">
                <h3>{{$totalPass}}</h3>
                <p>Passed</p>
            </div>
            <div class="stat-card partial">
                <h3>{{$totalPartial}}</h3>
                <p>Partial</p>
            </div>
            <div class="stat-card fail">
                <h3>{{$totalFail}}</h3>
                <p>Failed</p>
            </div>
        </div>

        <table class="summary-table">
            <thead>
                <tr>
                    <th>Domain</th>
                    <th>Total</th>
                    <th>Pass</th>
                    <th>Partial</th>
                    <th>Fail</th>
                    <th>N/A</th>
                    <th>Not Evaluated</th>
                </tr>
            </thead>
            <tbody>
                {{range .Report.DomainSummaries}}
                <tr>
                    <td>{{.Domain}}</td>
                    <td>{{.Total}}</td>
                    <td>{{.Pass}}</td>
                    <td>{{.Partial}}</td>
                    <td>{{.Fail}}</td>
                    <td>{{.NA}}</td>
                    <td>{{.NotEval}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
        {{else}}
        <p>No test results available.</p>
        {{end}}

        {{if .DomainResults}}
        <h2>Domain Results</h2>
        {{range $domain, $results := .DomainResults}}
        <h3>Domain: {{$domain}}</h3>
        {{range $results}}
        <div class="result-card">
            <div class="result-header">
                <span class="outcome-badge {{outcomeClass .Outcome}}">
                    {{outcomeSymbol .Outcome}} {{.Outcome.String}}
                </span>
                <h4>{{.ID}} - {{.Name}} (Level {{.Level.String}})</h4>
            </div>

            <table class="info-table">
                <tr><th>Testability</th><td>{{.Testability.String}}</td></tr>
                <tr><th>Spec Version</th><td>{{.SpecVersion}}</td></tr>
                {{if .Remarks}}<tr><th>Remarks</th><td>{{.Remarks}}</td></tr>{{end}}
            </table>

            {{if .NeedsReview}}
            <div class="needs-review">⚠ Manual Review Required</div>
            {{end}}

            {{if .Evidence}}
            <details>
                <summary>Evidence ({{len .Evidence}} item{{if ne (len .Evidence) 1}}s{{end}})</summary>
                {{range $i, $ev := .Evidence}}
                <div class="evidence-item">
                    <p><strong>Evidence {{add $i 1}}:</strong></p>
                    <p><strong>Command:</strong> <code>{{$ev.Command}}</code></p>
                    {{if $ev.Note}}<p><strong>Note:</strong> {{$ev.Note}}</p>{{end}}
                    <p><strong>Exit Code:</strong> {{$ev.ExitCode}}</p>
                    <p><strong>Duration:</strong> {{$ev.Duration}}</p>
                    {{if $ev.Stdout}}
                    <p><strong>Stdout:</strong></p>
                    <pre>{{$ev.Stdout}}</pre>
                    {{end}}
                    {{if $ev.Stderr}}
                    <p><strong>Stderr:</strong></p>
                    <pre>{{$ev.Stderr}}</pre>
                    {{end}}
                </div>
                {{end}}
            </details>
            {{end}}
        </div>
        {{end}}
        {{end}}
        {{end}}

        <h2>Testing Environment</h2>
        <table class="info-table">
            <tr><th>Operating System</th><td>{{.Report.Environment.OS}}</td></tr>
            <tr><th>Architecture</th><td>{{.Report.Environment.Arch}}</td></tr>
            {{if .Report.Environment.Terminal}}<tr><th>Terminal</th><td>{{.Report.Environment.Terminal}}</td></tr>{{end}}
            {{if .Report.Environment.Shell}}<tr><th>Shell</th><td>{{.Report.Environment.Shell}}</td></tr>{{end}}
        </table>
    </div>
</body>
</html>
`
