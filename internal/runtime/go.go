package runtime

type GoRuntime struct {
	image     string
	extension string
}

func (g *GoRuntime) Image() string {
	return g.image
}

func (g *GoRuntime) Extension() string {
	return g.extension
}
func (g *GoRuntime) Command(filename string) []string {
	return []string{"sh", "-c", "GO111MODULE=off go run /" + filename}
}
