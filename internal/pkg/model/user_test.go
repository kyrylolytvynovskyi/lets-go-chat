package model

import (
	"testing"
)

func TestCreateUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		r       *CreateUserRequest
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Username len less than 4",
			r:       &CreateUserRequest{UserName: "usr"},
			wantErr: true,
		},
		{
			"Password len less than 8", &CreateUserRequest{"user", "pwd"}, true,
		},
		{
			"Valid input data", &CreateUserRequest{"user", "password"}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("CreateUserRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoginUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		r       *LoginUserRequest
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Always Succeed",
			r:       &LoginUserRequest{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("LoginUserRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
