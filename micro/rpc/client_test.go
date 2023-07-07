package rpc

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_setFuncField(t *testing.T) {
	testCases := []struct {
		name    string
		service Service
		wantErr error
		mock    func(controller *gomock.Controller) Proxy
	}{
		{
			name:    "nil",
			wantErr: ErrRPCNotSupportNil,
			mock: func(controller *gomock.Controller) Proxy {
				return NewMockProxy(controller)
			},
		},
		{
			name:    "no pointer",
			service: UserService{},
			wantErr: ErrRPCOnlySupportPointOfStruct,
			mock: func(controller *gomock.Controller) Proxy {
				return NewMockProxy(controller)
			},
		},
		{
			name:    "normal",
			service: &UserService{},
			mock: func(controller *gomock.Controller) Proxy {
				p := NewMockProxy(controller)
				p.EXPECT().Invoke(gomock.Any(), gomock.Any()).Return(&Response{}, nil)
				return p
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			err := setFuncField(tc.service, tc.mock(ctrl))
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			resp, err := tc.service.(*UserService).GetById(context.Background(), &GetByIdReq{Id: 123})
			assert.Equal(t, tc.wantErr, err)
			t.Log(resp)

		})
	}

}
