package contracts

type View interface {
	View(resource Resource) string
}
