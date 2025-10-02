package version

var (
	Version   = "0.2.1"
	GitCommit = ""
)

func GetVersion() string {
	if GitCommit != "" {
		return Version + "-" + GitCommit[:min(len(GitCommit), 8)]
	}
	return Version
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
