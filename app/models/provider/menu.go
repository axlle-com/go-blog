package provider

type MenuProvider interface {
	GetMenuString(id uint, url string) (string, error)
}
