package data

type TriggerType uint

const (
	LARGER_THAN_TARGET TriggerType = iota + 1
	EQUAL_TO_TARGET
)
