package domain

import (
	"errors"
)

var (
	// This error throw when internal service fails to fulfill the request
	ErrInternal = errors.New("internal error")
	// this error throw when product that being requested is not found
	ErrProductNotFound = errors.New("product not found")
	// this error throw when product stock can't fulfill the request
	ErrInsufficientStock = errors.New("product stock is not enough")
)
