package command

import "cook-robot-controller-go/instruction"

const (
	COOK         = "cook"         // multiple
	WASH         = "wash"         // multiple
	RESET        = "reset"        // multiple
	UNLOCK       = "unlock"       // single
	DISH_OUT     = "dish_out"     // multiple
	RESUME       = "resume"       // single
	PAUSE_TO_ADD = "pause_to_add" // single
)

const (
	MULTIPLE = "multiple" // 不可在其他命令执行过程中执行
	SINGLE   = "single"   // 可在其他命令执行过程中执行
)

type Command struct {
	CommandType  string                      `json:"commandType"`
	CommandName  string                      `json:"commandName"`
	Instructions []instruction.Instructioner `json:"instructions"`
}
