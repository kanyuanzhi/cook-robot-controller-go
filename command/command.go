package command

import "cook-robot-controller-go/instruction"

type CommandType string

const (
	Single   = CommandType("single")
	Multiple = CommandType("multiple")
)

type Command struct {
	CommandType  CommandType                 `json:"commandType"`
	Instructions []instruction.Instructioner `json:"instructions"`
}
