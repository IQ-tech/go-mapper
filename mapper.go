package mapper

// MapperTag standards tag name for view mapping
const MapperTag string = "mapper"
const timeType string = "time.Time"

// Result holds result of mapper
type Result interface {
	Merge(src interface{}) Result
	To(tgr interface{}) error
}

// Mapper holds mapping operations
type Mapper interface {
	From(src interface{}) (retVal Result)
}

// New returns an instance of mapper
func New() Mapper {
	return &mapper{}
}
