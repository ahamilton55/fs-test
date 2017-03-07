package packager

type Packager interface {
	Build() (string, error)
	FindPackage(string) (string, error)
	Push(string) error
	Cleanup(string) error
}
