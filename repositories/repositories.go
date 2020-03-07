package repositories

type GetParam struct {
	Tablename string
	Filter    map[string]interface{}
	Result    interface{}
	Order     map[string]bool
	Offset    int
	Limit     int
}

type StoreParam struct {
	Tablename string
	Data      interface{}
}

type UpdateParam struct {
	Tablename string
	Filter    map[string]interface{}
	Data      interface{}
}

type DeleteParam struct {
	Tablename string
	Filter    map[string]interface{}
}

type NewsRepository interface {
	GetBy(param GetParam) error
	Store(param StoreParam) error
	Update(param UpdateParam) error
	Delete(param DeleteParam) error
}
