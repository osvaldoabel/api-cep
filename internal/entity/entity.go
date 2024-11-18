package entity

import "context"

type Address interface {
	GetAddressFields() (string, error)
}

type PostalCodeProvider interface {
	GetAddress(ctx context.Context) (Address, error)
}
