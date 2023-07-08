package data

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
