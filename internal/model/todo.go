package model

type ToDo struct {
	Line string
}

type FileToDos struct {
	FilePath       string
	ToDos          []ToDo
	Context        string
	ContextGravity int
	Gravity        int
}
