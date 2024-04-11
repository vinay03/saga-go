package saga

type Stage struct {
	Name     string
	Func     interface{}
	CompFunc interface{}
}
