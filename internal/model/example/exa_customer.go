package example

type ExampleModel struct {
	Id   int64  `json:"id"`
	Name string `xorm:"name" json:"name"`
}

func (e *ExampleModel) TableName() string {
	return "tab"
}

func (e *ExampleModel) DatabaseName() string {
	return "_t"
}
