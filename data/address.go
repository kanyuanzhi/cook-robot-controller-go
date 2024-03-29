package data

const (
	X_RESET_CONTROL_WORD_ADDRESS  = "DD0"
	X_RESET_STATUS_WORD_ADDRESS   = "DD2900"
	Y_RESET_CONTROL_WORD_ADDRESS  = "DD2"
	Y_RESET_STATUS_WORD_ADDRESS   = "DD2902"
	Z_RESET_CONTROL_WORD_ADDRESS  = "DD4"
	Z_RESET_STATUS_WORD_ADDRESS   = "DD2904"
	R1_RESET_CONTROL_WORD_ADDRESS = "DD6"
	R1_RESET_STATUS_WORD_ADDRESS  = "DD2906"
	R2_RESET_CONTROL_WORD_ADDRESS = "DD8"
	R2_RESET_STATUS_WORD_ADDRESS  = "DD2908"

	X_LOCATE_CONTROL_WORD_ADDRESS  = "DD20"
	X_LOCATE_STATUS_WORD_ADDRESS   = "DD2910"
	X_LOCATE_POSITION_ADDRESS      = "DD22"
	X_LOCATE_SPEED_ADDRESS         = "DD24"
	Y_LOCATE_CONTROL_WORD_ADDRESS  = "DD30"
	Y_LOCATE_STATUS_WORD_ADDRESS   = "DD2912"
	Y_LOCATE_POSITION_ADDRESS      = "DD32"
	Y_LOCATE_SPEED_ADDRESS         = "DD34"
	Z_LOCATE_CONTROL_WORD_ADDRESS  = "DD40"
	Z_LOCATE_STATUS_WORD_ADDRESS   = "DD2914"
	Z_LOCATE_POSITION_ADDRESS      = "DD42"
	Z_LOCATE_SPEED_ADDRESS         = "DD44"
	R1_LOCATE_CONTROL_WORD_ADDRESS = "DD50"
	R1_LOCATE_STATUS_WORD_ADDRESS  = "DD2916"
	R1_LOCATE_POSITION_ADDRESS     = "DD52"
	R1_LOCATE_SPEED_ADDRESS        = "DD54"
	R2_LOCATE_CONTROL_WORD_ADDRESS = "DD60"
	R2_LOCATE_STATUS_WORD_ADDRESS  = "DD2918"
	R2_LOCATE_POSITION_ADDRESS     = "DD62"
	R2_LOCATE_SPEED_ADDRESS        = "DD64"

	R1_ROTATE_CONTROL_WORD_ADDRESS = "DD80"
	R1_ROTATE_STATUS_WORD_ADDRESS  = "DD2920"
	R1_ROTATE_MODE_ADDRESS         = "DD82"
	R1_ROTATE_SPEED_ADDRESS        = "DD84"
	R1_ROTATE_AMOUNT_ADDRESS       = "DD86"

	DISH_OUT_CONTROL_WORD_ADDRESS      = "DD90"
	DISH_OUT_STATUS_WORD_ADDRESS       = "DD2922"
	DISH_OUT_AMOUNT_ADDRESS            = "DD92"
	DISH_OUT_UPWARD_SPEED_ADDRESS      = "DD94"
	DISH_OUT_DOWNWARD_SPEED_ADDRESS    = "DD96"
	DISH_OUT_UPWARD_POSITION_ADDRESS   = "DD98"
	DISH_OUT_DOWNWARD_POSITION_ADDRESS = "DD100"

	SHAKE_CONTROL_WORD_ADDRESS   = "DD110"
	SHAKE_STATUS_WORD_ADDRESS    = "DD2924"
	SHAKE_AMOUNT_ADDRESS         = "DD112"
	SHAKE_UPWARD_SPEED_ADDRESS   = "DD114"
	SHAKE_DOWNWARD_SPEED_ADDRESS = "DD116"
	SHAKE_DISTANCE_ADDRESS       = "DD118"

	PUMP_1_CONTROL_WORD_ADDRESS = "DD160"
	PUMP_1_STATUS_WORD_ADDRESS  = "DD2964"
	PUMP_1_DURATION_ADDRESS     = "DD162"
	PUMP_1_LIQUID_WARINIG       = "DD2986"

	PUMP_2_CONTROL_WORD_ADDRESS = "DD164"
	PUMP_2_STATUS_WORD_ADDRESS  = "DD2966"
	PUMP_2_DURATION_ADDRESS     = "DD166"
	PUMP_2_LIQUID_WARINIG       = "DD2988"

	PUMP_3_CONTROL_WORD_ADDRESS = "DD168"
	PUMP_3_STATUS_WORD_ADDRESS  = "DD2968"
	PUMP_3_DURATION_ADDRESS     = "DD170"
	PUMP_3_LIQUID_WARINIG       = "DD2990"

	PUMP_4_CONTROL_WORD_ADDRESS = "DD172"
	PUMP_4_STATUS_WORD_ADDRESS  = "DD2970"
	PUMP_4_DURATION_ADDRESS     = "DD174"
	PUMP_4_LIQUID_WARINIG       = "DD2992"

	PUMP_5_CONTROL_WORD_ADDRESS = "DD176"
	PUMP_5_STATUS_WORD_ADDRESS  = "DD2972"
	PUMP_5_DURATION_ADDRESS     = "DD178"
	PUMP_5_LIQUID_WARINIG       = "DD2994"

	PUMP_6_CONTROL_WORD_ADDRESS = "DD180"
	PUMP_6_STATUS_WORD_ADDRESS  = "DD2974"
	PUMP_6_DURATION_ADDRESS     = "DD182"
	PUMP_6_LIQUID_WARINIG       = "DD2996"

	PUMP_7_CONTROL_WORD_ADDRESS = "DD184"
	PUMP_7_STATUS_WORD_ADDRESS  = "DD2976"
	PUMP_7_DURATION_ADDRESS     = "DD186"

	PUMP_8_CONTROL_WORD_ADDRESS = "DD188"
	PUMP_8_STATUS_WORD_ADDRESS  = "DD2978"
	PUMP_8_DURATION_ADDRESS     = "DD190"

	PUMP_9_CONTROL_WORD_ADDRESS = "DD192"
	PUMP_9_STATUS_WORD_ADDRESS  = "DD2980"
	PUMP_9_DURATION_ADDRESS     = "DD194"

	PUMP_10_CONTROL_WORD_ADDRESS = "DD196"
	PUMP_10_STATUS_WORD_ADDRESS  = "DD2982"
	PUMP_10_DURATION_ADDRESS     = "DD198"

	LAMPBLACK_PURIFY_CONTROL_WORD_ADDRESS = "DD130"
	LAMPBLACK_PURIFY_STATUS_WORD_ADDRESS  = "DD2928"
	LAMPBLACK_PURIFY_MODE_ADDRESS         = "DD132"

	DOOR_UNLOCK_CONTROL_WORD_ADDRESS = "DD140"
	DOOR_UNLOCK_STATUS_WORD_ADDRESS  = "DD2930"

	TEMPERATURE_CONTROL_WORD_ADDRESS = "DD150"
	TEMPERATURE_STATUS_WORD_ADDRESS  = "DD2938"
	TEMPERATURE_ADDRESS              = "DD152"
	TEMPERATURE_BOTTOM_ADDRESS       = "DD2932"
	TEMPERATURE_INFRARED_ADDRESS     = "DD2934"
	TEMPERATURE_WARNING_ADDRESS      = "DD2936"

	TEMPERATURE_UPPER_VALUE_ADDRESS = "DD2950"
	TEMPERATURE_LOWER_VALUE_ADDRESS = "DD2960"

	WATER_SOURCE_VALVE_CONTROL_WORD_ADDRESS = "DD200"
	WATER_SOURCE_VALVE_STATUS_WORD_ADDRESS  = "DD2998"

	WATER_PUMP_VALVE_CONTROL_WORD_ADDRESS = "DD202"
	WATER_PUMP_VALVE_STATUS_WORD_ADDRESS  = "DD3000"

	NOZZLE_VALVE_CONTROL_WORD_ADDRESS = "DD204"
	NOZZLE_VALVE_STATUS_WORD_ADDRESS  = "DD3002"
)

var AxisControlWordAddressToAxis = map[string]string{
	X_RESET_CONTROL_WORD_ADDRESS:  "X",
	Y_RESET_CONTROL_WORD_ADDRESS:  "Y",
	Z_RESET_CONTROL_WORD_ADDRESS:  "Z",
	R1_RESET_CONTROL_WORD_ADDRESS: "R1",
	R2_RESET_CONTROL_WORD_ADDRESS: "R2",

	X_LOCATE_CONTROL_WORD_ADDRESS:  "X",
	Y_LOCATE_CONTROL_WORD_ADDRESS:  "Y",
	Z_LOCATE_CONTROL_WORD_ADDRESS:  "Z",
	R1_LOCATE_CONTROL_WORD_ADDRESS: "R1",
	R2_LOCATE_CONTROL_WORD_ADDRESS: "R2",
}

var AxisToResetControlWordAddress = map[string]string{
	"X":  X_RESET_CONTROL_WORD_ADDRESS,
	"Y":  Y_RESET_CONTROL_WORD_ADDRESS,
	"Z":  Z_RESET_CONTROL_WORD_ADDRESS,
	"R1": R1_RESET_CONTROL_WORD_ADDRESS,
	"R2": R2_RESET_CONTROL_WORD_ADDRESS,
}

var AxisToResetStatusWordAddress = map[string]string{
	"X":  X_RESET_STATUS_WORD_ADDRESS,
	"Y":  Y_RESET_STATUS_WORD_ADDRESS,
	"Z":  Z_RESET_STATUS_WORD_ADDRESS,
	"R1": R1_RESET_STATUS_WORD_ADDRESS,
	"R2": R2_RESET_STATUS_WORD_ADDRESS,
}

var AxisToLocateControlWordAddress = map[string]string{
	"X":  X_LOCATE_CONTROL_WORD_ADDRESS,
	"Y":  Y_LOCATE_CONTROL_WORD_ADDRESS,
	"Z":  Z_LOCATE_CONTROL_WORD_ADDRESS,
	"R1": R1_LOCATE_CONTROL_WORD_ADDRESS,
	"R2": R2_LOCATE_CONTROL_WORD_ADDRESS,
}

var AxisToLocateStatusWordAddress = map[string]string{
	"X":  X_LOCATE_STATUS_WORD_ADDRESS,
	"Y":  Y_LOCATE_STATUS_WORD_ADDRESS,
	"Z":  Z_LOCATE_STATUS_WORD_ADDRESS,
	"R1": R1_LOCATE_STATUS_WORD_ADDRESS,
	"R2": R2_LOCATE_STATUS_WORD_ADDRESS,
}

var AxisToLocatePositionAddress = map[string]string{
	"X":  X_LOCATE_POSITION_ADDRESS,
	"Y":  Y_LOCATE_POSITION_ADDRESS,
	"Z":  Z_LOCATE_POSITION_ADDRESS,
	"R1": R1_LOCATE_POSITION_ADDRESS,
	"R2": R2_LOCATE_POSITION_ADDRESS,
}

var AxisToLocateSpeedAddress = map[string]string{
	"X":  X_LOCATE_SPEED_ADDRESS,
	"Y":  Y_LOCATE_SPEED_ADDRESS,
	"Z":  Z_LOCATE_SPEED_ADDRESS,
	"R1": R1_LOCATE_SPEED_ADDRESS,
	"R2": R2_LOCATE_SPEED_ADDRESS,
}

var RotateModeToString = map[uint32]string{
	1: "正转",
	2: "反转",
	3: "正反转",
}

var PumpControlWordAddressToPumpNumber = map[string]uint32{
	PUMP_1_CONTROL_WORD_ADDRESS:  1,
	PUMP_2_CONTROL_WORD_ADDRESS:  2,
	PUMP_3_CONTROL_WORD_ADDRESS:  3,
	PUMP_4_CONTROL_WORD_ADDRESS:  4,
	PUMP_5_CONTROL_WORD_ADDRESS:  5,
	PUMP_6_CONTROL_WORD_ADDRESS:  6,
	PUMP_7_CONTROL_WORD_ADDRESS:  7,
	PUMP_8_CONTROL_WORD_ADDRESS:  8,
	PUMP_9_CONTROL_WORD_ADDRESS:  9,
	PUMP_10_CONTROL_WORD_ADDRESS: 10,
}

var PumpNumberToPumpControlWordAddress = map[uint32]string{
	1:  PUMP_1_CONTROL_WORD_ADDRESS,
	2:  PUMP_2_CONTROL_WORD_ADDRESS,
	3:  PUMP_3_CONTROL_WORD_ADDRESS,
	4:  PUMP_4_CONTROL_WORD_ADDRESS,
	5:  PUMP_5_CONTROL_WORD_ADDRESS,
	6:  PUMP_6_CONTROL_WORD_ADDRESS,
	7:  PUMP_7_CONTROL_WORD_ADDRESS,
	8:  PUMP_8_CONTROL_WORD_ADDRESS,
	9:  PUMP_9_CONTROL_WORD_ADDRESS,
	10: PUMP_10_CONTROL_WORD_ADDRESS,
}

var PumpNumberToPumpStatusWordAddress = map[uint32]string{
	1:  PUMP_1_STATUS_WORD_ADDRESS,
	2:  PUMP_2_STATUS_WORD_ADDRESS,
	3:  PUMP_3_STATUS_WORD_ADDRESS,
	4:  PUMP_4_STATUS_WORD_ADDRESS,
	5:  PUMP_5_STATUS_WORD_ADDRESS,
	6:  PUMP_6_STATUS_WORD_ADDRESS,
	7:  PUMP_7_STATUS_WORD_ADDRESS,
	8:  PUMP_8_STATUS_WORD_ADDRESS,
	9:  PUMP_9_STATUS_WORD_ADDRESS,
	10: PUMP_10_STATUS_WORD_ADDRESS,
}

var PumpNumberToPumpDurationAddress = map[uint32]string{
	1:  PUMP_1_DURATION_ADDRESS,
	2:  PUMP_2_DURATION_ADDRESS,
	3:  PUMP_3_DURATION_ADDRESS,
	4:  PUMP_4_DURATION_ADDRESS,
	5:  PUMP_5_DURATION_ADDRESS,
	6:  PUMP_6_DURATION_ADDRESS,
	7:  PUMP_7_DURATION_ADDRESS,
	8:  PUMP_8_DURATION_ADDRESS,
	9:  PUMP_9_DURATION_ADDRESS,
	10: PUMP_10_DURATION_ADDRESS,
}

var PumpNumberToPumpType = map[uint32]string{
	1:  "液体泵",
	2:  "液体泵",
	3:  "液体泵",
	4:  "液体泵",
	5:  "液体泵",
	6:  "液体泵",
	7:  "抽水泵",
	8:  "抽水泵",
	9:  "固体泵",
	10: "固体泵",
}

var LampblackPurifyModeToString = map[uint32]string{
	1: "排气",
	2: "排气+净化",
}
