package repo

import (
	"gophermart/internal/domain/models"
	"reflect"
	"testing"
)

func TestUnmarshalBalance(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    *models.Balance
		wantErr bool
	}{
		{"good json",
			"{\n\t\"current\": 500.5,\n\t\"withdrawn\": 42\n}",
			&models.Balance{
				Current:   500.5,
				Withdrawn: 42},
			false},
		{"wrong json",
			"\n\t\"current\": 500.5,\n\t\"withdrawn\": 42\n}",
			nil,
			true},
		{"empty",
			"",
			nil,
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalBalance(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalBalance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalOrder(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    *models.Order
		wantErr bool
	}{
		{"good json",
			"{\n\"order\":\"9278923470\",\n\"status\":\"PROCESSED\",\n\"accrual\":500,\n\"uploaded_at\":\"2020-12-10T15:15:45+03:00\"\n}",
			&models.Order{
				Number:     "9278923470",
				Status:     "PROCESSED",
				Accrual:    500,
				UploadedAt: "2020-12-10T15:15:45+03:00",
			},
			false},
		{"wrong json",
			"\n\"order\":\"9278923470\",\n\"status\":\"PROCESSED\",\n\"accrual\":500,\n\"uploaded_at\":\"2020-12-10T15:15:45+03:00\"\n}",
			nil,
			true},
		{"empty",
			"",
			nil,
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalOrder(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalOrder() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalUser(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    *models.User
		wantErr bool
	}{
		{"good json",
			"{\n\t\"login\": \"<login>\",\n\t\"password\": \"<password>\"\n}",
			&models.User{
				Login:    "<login>",
				Password: "<password>",
			},
			false},
		{"wrong json",
			"\n\t\"login\": \"<login>\",\n\t\"password\": \"<password>\"\n}",
			nil,
			true},
		{"empty",
			"",
			nil,
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalUser(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalWithdraw(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    *models.Withdraw
		wantErr bool
	}{
		{"good json",
			"{\n\t\"order\": \"2377225624\",\n    \"sum\": 751\n}",
			&models.Withdraw{
				Order: "2377225624",
				Sum:   751,
			},
			false},
		{"wrong json",
			"\n\t\"order\": \"2377225624\",\n    \"sum\": 751\n}",
			nil,
			true},
		{"empty",
			"",
			nil,
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalWithdraw(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalWithdraw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalWithdraw() got = %v, want %v", got, tt.want)
			}
		})
	}
}
