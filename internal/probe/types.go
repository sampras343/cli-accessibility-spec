package probe

type ProbeResult struct {
	BinaryPath     string
	BinaryName     string
	Version        string
	HelpText       string
	Subcommands    []Subcommand
	GlobalFlags    []Flag
	HasColor       bool
	HasHelp        bool
	HasVersion     bool
	HasSubcommands bool
	HasJSONFlag    bool
	HasQuietFlag   bool
	HasColorFlag   bool
	HasDryRunFlag  bool
	HasNoInputFlag bool
	HelpFormat     string
	SampleErrors   []ErrorSample
}

type Subcommand struct {
	Name     string
	HelpText string
	Flags    []Flag
}

type Flag struct {
	Short       string
	Long        string
	Description string
	TakesValue  bool
}

type ErrorSample struct {
	Trigger  string
	Stdout   []byte
	Stderr   []byte
	ExitCode int
}
