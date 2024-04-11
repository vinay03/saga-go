package saga

func NewSaga(id string) *Saga {
	return &Saga{
		ID: id,
	}
}

type StagesList []*Stage

type Saga struct {
	ID     string
	Stages StagesList
}

func (sg *Saga) AddStage(st *Stage) error {
	err := st.Verify()
	if err != nil {
		return err
	}
	return nil
}
func (sg *Saga) VerifyStageConfig(st *Stage) {}
func (sg *Saga) Start()                      {}
