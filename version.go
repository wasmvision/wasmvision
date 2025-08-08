package wasmvision

var (
	version = "0.5.0-dev"
	sha     string
)

func Version() string {
	if sha != "" {
		return version + "-" + sha
	}
	return version
}
