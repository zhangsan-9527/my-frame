package net

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func Test_handleConn(t *testing.T) {
	testCases := []struct {
		name string

		mock    func(ctrl *gomock.Controller) net.Conn
		wantErr error
	}{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			err := handleConn(tc.mock(ctrl))
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
