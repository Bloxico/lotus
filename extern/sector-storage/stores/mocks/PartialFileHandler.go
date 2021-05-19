// Code generated by mockery 2.7.5. DO NOT EDIT.

package mocks

import (
	abi "github.com/filecoin-project/go-state-types/abi"
	mock "github.com/stretchr/testify/mock"

	partialfile "github.com/filecoin-project/lotus/extern/sector-storage/partialfile"

	storiface "github.com/filecoin-project/lotus/extern/sector-storage/storiface"
)

// PartialFileHandler is an autogenerated mock type for the PartialFileHandler type
type PartialFileHandler struct {
	mock.Mock
}

// HasAllocated provides a mock function with given fields: pf, offset, size
func (_m *PartialFileHandler) HasAllocated(pf *partialfile.PartialFile, offset storiface.UnpaddedByteIndex, size abi.UnpaddedPieceSize) (bool, error) {
	ret := _m.Called(pf, offset, size)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*partialfile.PartialFile, storiface.UnpaddedByteIndex, abi.UnpaddedPieceSize) bool); ok {
		r0 = rf(pf, offset, size)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*partialfile.PartialFile, storiface.UnpaddedByteIndex, abi.UnpaddedPieceSize) error); ok {
		r1 = rf(pf, offset, size)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OpenPartialFile provides a mock function with given fields: maxPieceSize, path
func (_m *PartialFileHandler) OpenPartialFile(maxPieceSize abi.PaddedPieceSize, path string) (*partialfile.PartialFile, error) {
	ret := _m.Called(maxPieceSize, path)

	var r0 *partialfile.PartialFile
	if rf, ok := ret.Get(0).(func(abi.PaddedPieceSize, string) *partialfile.PartialFile); ok {
		r0 = rf(maxPieceSize, path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*partialfile.PartialFile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(abi.PaddedPieceSize, string) error); ok {
		r1 = rf(maxPieceSize, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}