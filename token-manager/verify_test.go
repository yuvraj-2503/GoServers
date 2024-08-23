package token_manager

import (
	"reflect"
	"testing"
)

func TestJwtTokenManager_Verify(t *testing.T) {
	type fields struct {
		secretKey string
	}
	type args struct {
		apiToken string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TokenClaims
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:   "Verify Token",
			fields: fields{secretKey: "uNxxCyuO5i1YvM6QpwTsnq5njXwYjnun7k8zMT7vZsO6YD9CiZlOBcJwxQpUfrh5ZMZ6BgDXn4NK2vwweMTaT0rkCJuGnray"},
			args: args{
				apiToken: "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmhjSEFpT2lKVVJWTlVJaXdpWVhWMGFHOXlhWHBsWkNJNmRISjFaU3dpWlcxaGFXd2lPaUowWlhOMExXVnRZV2xzSWl3aVpYaHdJam94TnpJME5EYzNPVFV6TENKcFlYUWlPakUzTWpRek9URTFOVE1zSW1wMGFTSTZJbFJGVTFRaUxDSnJhVzVrSWpvaVZFVlRWQ0lzSW0xaFkyaHBibVVpT2lKMFpYTjBMVzFoWTJocGJtVWlMQ0p0WVdOb2FXNWxYMmxrSWpvaWRHVnpkQzF0WVdOb2FXNWxJaXdpYzNWaUlqb2llWFYyY21GcUxuTnBibWRvUUhwcGNtOW9MbU52YlNJc0luVnpaWElpT2lKMFpYTjBMWFZ6WlhJaUxDSjFjMlZ5U1dRaU9pSjBaWE4wTFhWelpYSWlmUS5vSWRMcGNacTJzdmc2MW15ZkVGYU9UVk0zRFQ2TzQ2WWxkRzZGeWN1a1BnanN2SXcxdXJkV2lrNlpqN2tfWlJtdzZ0LTlzeXVSeTNBWjl5Nkl3RnRIQQ", // put your jwt token here
			},
			want: &TokenClaims{
				UserId:    "test-user",
				EmailId:   "test-email",
				MachineId: "test-machine",
				App:       "TEST",
				IAT:       &IAT,
				EXP:       &EXP,
				Sub:       "yuvraj.singh@ziroh.com",
				Kind:      "TEST",
				JTI:       "TEST",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := NewJwtTokenManager(tt.fields.secretKey)
			got, err := j.Verify(tt.args.apiToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Verify() got = %v, want %v", got, tt.want)
			}
		})
	}
}
