package cmd

// Version of current binary tool
const Version = "0.1.0"

type (
	// Opts is struct docopt binds flag data to
	Opts struct {
		Config bool `docopt:"config"`
	}
)
