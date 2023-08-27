package data

import "cook-robot-controller-go/config"

var X_MOVE_SPEED = config.Parameter.XAxis.MoveSpeed

const (
	X_READY_POSITION uint32 = iota + 1
	X_WITHDRAWER_POSITION
	X_BOX_1_POSITION
	X_BOX_2_POSITION
	X_BOX_3_POSITION
	X_BOX_4_POSITION
	X_SAFE_POSITION
)

var SlotNumberToPosition = map[uint32]uint32{
	1: X_BOX_1_POSITION,
	2: X_BOX_2_POSITION,
	3: X_BOX_3_POSITION,
	4: X_BOX_4_POSITION,
}

var XPositionToDistance = map[uint32]uint32{
	X_READY_POSITION:      config.Parameter.XAxis.ReadyPosition,      // 上菜位
	X_WITHDRAWER_POSITION: config.Parameter.XAxis.WithdrawerPosition, // 收纳位
	X_BOX_1_POSITION:      config.Parameter.XAxis.Box1Position,       // 菜仓1位
	X_BOX_2_POSITION:      config.Parameter.XAxis.Box2Position,       // 菜仓2位
	X_BOX_3_POSITION:      config.Parameter.XAxis.Box3Position,       // 菜仓3位
	X_BOX_4_POSITION:      config.Parameter.XAxis.Box4Position,       // 菜仓4位
	X_SAFE_POSITION:       config.Parameter.XAxis.SafePosition,       // 安全位
}

var SHAKE_AMOUNT = config.Parameter.XAxis.ShakeAmount
var SHAKE_UPWARD_SPEED = config.Parameter.XAxis.ShakeUpwardSpeed
var SHAKE_DOWNWARD_SPEED = config.Parameter.XAxis.ShakeDownwardSpeed
var SHAKE_DISTANCE = config.Parameter.XAxis.ShakeDistance

var Y_MOVE_SPEED = config.Parameter.YAxis.MoveSpeed

const (
	Y_INGREDIENT_POSITION uint32 = iota + 1
	Y_LIQUID_SEASONING_POSITION
	Y_SOLID_SEASONING_POSITION
	Y_STIR_FRY_1_POSITION
	Y_STIR_FRY_2_POSITION
	Y_STIR_FRY_3_POSITION
	Y_WASH_POSITION
	Y_POUR_POSITION
	Y_DISH_OUT_HIGH_POSITION
	Y_DISH_OUT_LOW_POSITION
)

var YPositionToDistance = map[uint32]uint32{
	Y_INGREDIENT_POSITION:       config.Parameter.YAxis.IngredientPosition,      // 接菜位
	Y_LIQUID_SEASONING_POSITION: config.Parameter.YAxis.LiquidSeasoningPosition, // 接液体料位
	Y_SOLID_SEASONING_POSITION:  config.Parameter.YAxis.SolidSeasoningPosition,  // 接固体料位
	Y_STIR_FRY_1_POSITION:       config.Parameter.YAxis.StirFry1Position,        // 炒菜1位
	Y_STIR_FRY_2_POSITION:       config.Parameter.YAxis.StirFry2Position,        // 炒菜2位
	Y_STIR_FRY_3_POSITION:       config.Parameter.YAxis.StirFry3Position,        // 炒菜3位
	Y_WASH_POSITION:             config.Parameter.YAxis.WashPosition,            // 洗锅位
	Y_POUR_POSITION:             config.Parameter.YAxis.PourPosition,            // 倒位水
	Y_DISH_OUT_HIGH_POSITION:    config.Parameter.YAxis.DishOutHighPosition,     // 出菜高位
	Y_DISH_OUT_LOW_POSITION:     config.Parameter.YAxis.DishOutLowPosition,      // 出菜低位
}

var R1_MAX_ROTATE_SPEED = config.Parameter.R1Axis.MaxRotateSpeed

var DISH_OUT_AMOUNT = config.Parameter.YAxis.DishOutAmount
var DISH_OUT_UPWARD_SPEED = config.Parameter.YAxis.DishOutUpwardSpeed
var DISH_OUT_DOWNWARD_SPEED = config.Parameter.YAxis.DishOutDownwardSpeed

var WASH_ROTATE_SPEED = config.Parameter.Wash.RotateSpeed
var WASH_ROTATE_CROSS_AMOUNT = config.Parameter.Wash.RotateCrossAmount
var WASH_ADD_WATER_DURATION = config.Parameter.Wash.AddWaterDuration
var WASH_TEMPERATURE = config.Parameter.Wash.Temperature
var WASH_DURATION = config.Parameter.Wash.Duration
var WASH_PUMP_NUMBER = config.Parameter.Wash.PumpNumber
var WASH_POUR_WATER_DURATION = config.Parameter.Wash.PourWaterDuration

var LAMPBLACK_PURIFY_ENABLE = config.Parameter.LampblackPurify.Enable
var LAMPBLACK_PURIFY_VENTING_MODE = config.Parameter.LampblackPurify.VentingMode
var LAMPBLACK_PURIFY_PURIFICATION_MODE = config.Parameter.LampblackPurify.PurificationMode
var LAMPBLACK_PURIFY_AUTO_START = config.Parameter.LampblackPurify.AutoStart
var LAMPBLACK_PURIFY_AUTO_START_MODE = config.Parameter.LampblackPurify.AutoStartMode
