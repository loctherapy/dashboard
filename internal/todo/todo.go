package todo

type ToDo struct {
    Line string
}

type FileToDos struct {
    FilePath string
    ToDos    []ToDo
    Context  string
    Gravity  int
}
