package dotfiles

type Dotfile struct {
	Id        string
	Name      string
	Path      string
	Extension string
	MimeType  string
	Language  string
	Directory DotfileDir
}

type DotfileDir struct {
	Id   string
	Name string
	Path string
}