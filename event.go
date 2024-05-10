package saga

type EventData struct {
	Key string

	SagaId string
	Saga   *Saga

	StageId string
	Stage   *Stage

	Action string

	Data interface{}
}

func (ed *EventData) GetKey() string {
	key := ed.SagaId
	return key
}

func GetEventData(eventKey string, data interface{}) *EventData {
	ed := EventData{}
	ed.Set(eventKey, data)
	return &ed
}

func (ed *EventData) Set(eventKey string, data interface{}) {

}
