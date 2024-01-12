package util

import (
	"github.com/mxrcury/rootgo/types"
)

// Create new error
func NewError(message string, statusCode int) *types.Error {
	return &types.Error{Message: message, Status: statusCode}
}
