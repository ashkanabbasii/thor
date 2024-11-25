// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package abi

import (
	"bytes"
	"errors"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
)

// Method see abi.Method in go-ethereum.
type Method struct {
	id     MethodID
	method *ethabi.Method
}

// ID returns method id.
func (m *Method) ID() MethodID {
	return m.id
}

// Name returns method name.
func (m *Method) Name() string {
	return m.method.Name
}

// Const returns if the method is const.
func (m *Method) Const() bool {
	return m.method.Constant
}

// EncodeInput encode args to data, and the data is prefixed with method id.
func (m *Method) EncodeInput(args ...interface{}) ([]byte, error) {
	data, err := m.method.Inputs.Pack(args...)
	if err != nil {
		return nil, err
	}
	return append(m.id[:], data...), nil
}

// DecodeInput decode input data into args.
func (m *Method) DecodeInput(input []byte, v interface{}) error {
	if !bytes.HasPrefix(input, m.id[:]) {
		return errors.New("input has incorrect prefix")
	}
	v, err := m.method.Inputs.Unpack(input[4:])
	if err != nil {
		return err
	}
	return nil
}

// EncodeOutput encode output args to data.
func (m *Method) EncodeOutput(args ...interface{}) ([]byte, error) {
	return m.method.Outputs.Pack(args...)
}

// DecodeOutput decode output data.
func (m *Method) DecodeOutput(output []byte, v interface{}) error {
	if len(output)%32 != 0 {
		return errors.New("output has incorrect length")
	}
	v, err := m.method.Outputs.Unpack(output)
	if err != nil {
		return err
	}
	return nil
}

// MethodID method id.
type MethodID [4]byte

// ExtractMethodID extract method id from input data.
func ExtractMethodID(input []byte) (id MethodID, err error) {
	if len(input) < len(id) {
		err = errors.New("input data too short")
		return
	}
	copy(id[:], input)
	return
}
