package spec

type Application struct {
	Function Function `hcl:"function"`
}

type Function struct {
	Name    string   `hcl:"name"`
	Call    string   `hcl:"call"`
	Files   []string `hcl:"files"`
	Trigger string   `hcl:"trigger"`
	Method  string   `hcl:"method"`
}
