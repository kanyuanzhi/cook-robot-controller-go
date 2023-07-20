package instruction

import (
	"cook-robot-controller-go/action"
	"cook-robot-controller-go/core"
	"cook-robot-controller-go/data"
)

func (d DoorUnlockInstruction) ExecuteImmediately(controller *core.Controller) {
	doorUnlockControlAction := action.NewDoorUnlockControlAction(data.DOOR_UNLOCK_CONTROL_WORD_ADDRESS,
		data.DOOR_UNLOCK_STATUS_WORD_ADDRESS)
	controller.ExecuteImmediately(doorUnlockControlAction)
}

func (p PauseToAddInstruction) ExecuteImmediately(controller *core.Controller) {
	controller.IsPausing = true
	controller.IsPausingWithMovingFinished = false
	controller.Pause()
	yLocateControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_STIR_FRY_3_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	controller.ExecuteImmediately(yLocateControlAction)
	triggerAction := action.NewTriggerAction(data.NewAddressValue(data.Y_LOCATE_STATUS_WORD_ADDRESS, 100), data.EQUAL_TO_TARGET)
	controller.ExecuteImmediately(triggerAction)
	controller.IsPausingWithMovingFinished = true
	doorUnlockControlAction := action.NewDoorUnlockControlAction(data.DOOR_UNLOCK_CONTROL_WORD_ADDRESS,
		data.DOOR_UNLOCK_STATUS_WORD_ADDRESS)
	controller.ExecuteImmediately(doorUnlockControlAction)
}

func (p ResumeInstruction) ExecuteImmediately(controller *core.Controller) {
	controller.IsPausing = false
	controller.Resume()
}
