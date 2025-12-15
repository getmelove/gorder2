package app

// CQRS
type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
}

type Queries struct {
}
