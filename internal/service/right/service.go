package right

const (
	RightRead       string = "read"
	RightWrite      string = "write"
	RightReadValue  int    = 1
	RightWriteValue int    = 2
)

var rightsValues = map[string]int{
	RightRead:  RightReadValue,
	RightWrite: RightWriteValue,
}

type Service struct{}

func NewService() Service {
	return Service{}
}

func (s Service) SumRights(rights ...string) int {
	var sum int
	for _, r := range rights {
		v, _ := rightsValues[r]
		sum += v
	}

	return sum
}
