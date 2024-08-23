package token_manager

import (
	"testing"
	"time"
)

var IAT = time.Now()
var EXP = IAT.Add(24 * 60 * time.Minute)

func TestJwtTokenManager_Generate(t *testing.T) {
	type fields struct {
		secretKey string
	}
	type args struct {
		claims *TokenClaims
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:   "Generate Token",
			fields: fields{secretKey: "uNxxCyuO5i1YvM6QpwTsnq5njXwYjnun7k8zMT7vZsO6YD9CiZlOBcJwxQpUfrh5ZMZ6BgDXn4NK2vwweMTaT0rkCJuGnray"},
			args: args{
				&TokenClaims{
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
			},
			want:    "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmhjSEFpT2lKVVJWTlVJaXdpWVhWMGFHOXlhWHBsWkNJNmRISjFaU3dpWlcxaGFXd2lPaUowWlhOMExXVnRZV2xzSWl3aVpYaHdJam94TnpJME5EYzNPVFV6TENKcFlYUWlPakUzTWpRek9URTFOVE1zSW1wMGFTSTZJbFJGVTFRaUxDSnJhVzVrSWpvaVZFVlRWQ0lzSW0xaFkyaHBibVVpT2lKMFpYTjBMVzFoWTJocGJtVWlMQ0p0WVdOb2FXNWxYMmxrSWpvaWRHVnpkQzF0WVdOb2FXNWxJaXdpYzNWaUlqb2llWFYyY21GcUxuTnBibWRvUUhwcGNtOW9MbU52YlNJc0luVnpaWElpT2lKMFpYTjBMWFZ6WlhJaUxDSjFjMlZ5U1dRaU9pSjBaWE4wTFhWelpYSWlmUS5vSWRMcGNacTJzdmc2MW15ZkVGYU9UVk0zRFQ2TzQ2WWxkRzZGeWN1a1BnanN2SXcxdXJkV2lrNlpqN2tfWlJtdzZ0LTlzeXVSeTNBWjl5Nkl3RnRIQQ",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := NewJwtTokenManager(tt.fields.secretKey)
			got, err := j.Generate(tt.args.claims)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Generate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
