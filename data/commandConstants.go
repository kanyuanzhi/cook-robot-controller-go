package data

const (
	COMMAND_NAME_COOK         = "cook"         // multiple
	COMMAND_NAME_WASH         = "wash"         // multiple
	COMMAND_NAME_POUR         = "pour"         // multiple
	COMMAND_NAME_PREPARE      = "prepare"      // multiple
	COMMAND_NAME_DOOR_UNLOCK  = "door_unlock"  // single
	COMMAND_NAME_DISH_OUT     = "dish_out"     // multiple
	COMMAND_NAME_RESUME       = "resume"       // single
	COMMAND_NAME_PAUSE_TO_ADD = "pause_to_add" // single
	COMMAND_NAME_HEAT         = "heat"         // single
	COMMAND_NAME_WITHDRAW     = "withdraw"     // multiple
	COMMAND_NAME_SHUTDOWN     = "shutdown"     // single

	COMMAND_NAME_OPEN_WATER_SOURCE_VALVE  = "open_water_source_valve"  // single
	COMMAND_NAME_CLOSE_WATER_SOURCE_VALVE = "close_water_source_valve" // single
	COMMAND_NAME_OPEN_WATER_PUMP_VALVE    = "open_water_pump_valve"    // single
	COMMAND_NAME_CLOSE_WATER_PUMP_VALVE   = "close_water_pump_valve"   // single
	COMMAND_NAME_OPEN_NOZZLE_VALVE        = "open_nozzle_valve"        // single
	COMMAND_NAME_OPEN_PUMP_7_VALVE        = "open_pump_7_valve"        // single
	COMMAND_NAME_CLOSE_PUMP_7_VALVE       = "close_pump_7_valve"       // single
	COMMAND_NAME_CLOSE_NOZZLE_VALVE       = "close_nozzle_valve"       // single
)

const (
	COMMAND_TYPE_MULTIPLE = "multiple" // 不可在其他命令执行过程中执行
	COMMAND_TYPE_SINGLE   = "single"   // 可在其他命令执行过程中执行
)
