package command

import "cook-robot-controller-go/instruction"

type Command struct {
	CommandType  string                      `json:"commandType"`
	CommandName  string                      `json:"commandName"`
	DishUUID     string                      `json:"dishUUID"` //如果是炒制命令，则会携带菜品的uuid
	Instructions []instruction.Instructioner `json:"instructions"`
}
