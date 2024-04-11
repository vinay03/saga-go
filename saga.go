package saga

func NewSaga()

type Saga struct {
	ID     string
	Stages []Stage
}
