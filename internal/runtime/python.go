package runtime

type PythonRuntime struct {
	image     string
	extension string
}

func (p *PythonRuntime) Image() string {
	return p.image
}

func (p *PythonRuntime) Extension() string {
	return p.extension
}
func (p *PythonRuntime) Command(filename string) []string {
	return []string{"python", "/", filename}
}
