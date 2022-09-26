package domain

type Secret struct {
	ID      int
	Message string
	Content string
}

type NewSecret struct {
	Message string
	Content string
}
