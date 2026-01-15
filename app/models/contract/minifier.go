package contract

type Minifier interface {
	Run() error
	Bundle(mediaType string, inputPaths []string) (string, error)
}
