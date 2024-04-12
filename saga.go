package saga

func NewSaga(id string) *Saga {
	return &Saga{
		ID: id,
	}
}

type StageList []*Stage

type Saga struct {
	ID     string
	Stages StageList
}

func (sg *Saga) AddStages(sl ...*Stage) error {
	for _, st := range sl {
		err := sg.AddStage(st)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sg *Saga) AddStage(st *Stage) error {
	err := st.Verify()
	if err != nil {
		return err
	}
	sg.Stages = append(sg.Stages, st)
	return nil
}
func (sg *Saga) Start() {

	for _, stg := range sg.Stages {
		stg.Action()
	}

}
