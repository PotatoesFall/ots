package domain

type Secret struct {
	ID      int
	Message string
	Content string
	Hash    string
}

type NewSecret struct {
	Message string
	Content string
	Hash    string
}
