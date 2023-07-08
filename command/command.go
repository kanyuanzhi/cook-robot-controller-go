package command

import "cook-robot-controller-go/instruction"

type CommandType string

const (
	Single   = CommandType("single")
	Multiple = CommandType("multiple")
)

type Command struct {
	CommandType  CommandType                 `json:"command_type"`
	Instructions []instruction.Instructioner `json:"instructions"`
}
