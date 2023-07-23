package data

type InstructionInfo struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Index        int    `json:"index"`
	ActionNumber int    `json:"actionNumber"`
}

func NewInstructionInfo(Type string, name string, actionNumber int) *InstructionInfo {
	return &InstructionInfo{
		Type:         Type,
		Name:         name,
		Index:        0,
		ActionNumber: actionNumber,
	}
}
