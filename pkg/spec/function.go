package spec

type Application struct {
	Function Function `hcl:"function"`
}

type Function struct {
	Name string `hcl:"name"`
}
