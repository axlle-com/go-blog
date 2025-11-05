package contract

type View interface {
	View(resource Resource) string
	ViewStatic(name string) string
}
