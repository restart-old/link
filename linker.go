package link

type Linker struct {
	Token string
	Storer
}

func NewLinker(token string) *Linker {
	return &Linker{
		Token: token,
	}
}
