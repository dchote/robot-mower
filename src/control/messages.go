package control

type StateMessage struct {
	*MowerStateStruct
	Namespace string `json:"namespace"`
	Mutation  string `json:"mutation"`
}

type CommandMessage struct {
	Method string `json:"method"`
	Value  string `json:"value"`
}
