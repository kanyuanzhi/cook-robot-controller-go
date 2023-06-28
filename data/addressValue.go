package data

type AddressValue struct {
	address string
	value   uint32
}

func NewAddressValue(address string, value uint32) *AddressValue {
	return &AddressValue{
		address: address,
		value:   value,
	}
}
