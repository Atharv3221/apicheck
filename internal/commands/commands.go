package commands

type Command string

const (
	Run    Command = "run"
	List   Command = "list"
	Add    Command = "add"
	Remove Command = "remove"
	Scope  Command = "scope"
	Init   Command = "init"
)
