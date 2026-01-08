package invocation

// CommandLineData describes the command line used for a bazel
// invocation.
type CommandLineData struct {
	Executable     string              `json:"executable"`
	Command        string              `json:"command"`
	Options        []CommandLineOption `json:"options"`
	StartupOptions []CommandLineOption `json:"startupOptions"`
	Residual       []string            `json:"residual"`
}

// CommandLineOption describes a single command line option (e.g. --foo=bar)
type CommandLineOption struct {
	Option string `json:"option"`
	Value  string `json:"value"`
}

// ParsedCommandLineOptions describes the options used for a bazel
// invocation before normalization.
type ParsedCommandLineOptions struct {
	ExplicitOptions []string `json:"explicitOptions"`
	Options         []string `json:"options"`
}
