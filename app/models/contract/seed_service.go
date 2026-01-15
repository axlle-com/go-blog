package contract

type SeedService interface {
	GetFiles(layout, moduleName string) (map[string]string, error)
	IsApplied(name string) (bool, error)
	MarkApplied(name string) error
}
