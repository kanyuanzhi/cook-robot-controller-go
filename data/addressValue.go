package data

import "fmt"

type AddressValue struct {
	Address string
	Value   uint32
}

func NewAddressValue(address string, value uint32) *AddressValue {
	return &AddressValue{
		Address: address,
		Value:   value,
	}
}

func (a *AddressValue) ToString() string {
	return fmt.Sprintf("(%s, %d)", a.Address, a.Value)
}
