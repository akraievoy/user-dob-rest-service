package main

import (
	"reflect"
	"testing"

	"github.com/akraievoy/tsv_load/proto"
	_ "github.com/lib/pq"
)

func Test_internationalizePhoneNumber(t *testing.T) {
	type args struct {
		phoneNumberLocalFormat string
	}
	tests := []struct {
		name                         string
		args                         args
		wantCountryCode              string
		wantPhoneNumberInternational string
	}{
		{
			"first dataset line",
			args{"056 5247 1778"},
			"+44", "56 5247 1778",
		},
		//	https://www.area-codes.org.uk/formatting.php#overseas
		{
			"area-codes example",
			args{"(020) 7946 0018"},
			"+44", "20 7946 0018",
		},
		{
			"another seemingly tricky case from dataset",
			args{"(01451) 68984"},
			"+44", "1451 68984",
		},
		{
			"this one is just used often",
			args{"0800 1111"},
			"+44", "800 1111",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCountryCode, gotPhoneNumberInternational := internationalizePhoneNumber(tt.args.phoneNumberLocalFormat)
			if gotCountryCode != tt.wantCountryCode {
				t.Errorf("internationalizePhoneNumber() gotCountryCode = %v, want %v", gotCountryCode, tt.wantCountryCode)
			}
			if gotPhoneNumberInternational != tt.wantPhoneNumberInternational {
				t.Errorf("internationalizePhoneNumber() gotPhoneNumberInternational = %v, want %v", gotPhoneNumberInternational, tt.wantPhoneNumberInternational)
			}
		})
	}
}

func Test_userToUpsertArgs(t *testing.T) {
	type args struct {
		u *proto.User
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			"simple test",
			args{
				createUser(1234, "Name", "e@ma.il", "(020) 7946 0018"),
			},
			[]interface{}{uint32(1234), "Name", "e@ma.il", "+44", "20 7946 0018"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := userToUpsertArgs(tt.args.u); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userToUpsertArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createUser(id int32, name string, email string, phoneNumber string) *proto.User {
	var user proto.User
	user.Id = id
	user.Name = name
	user.Email = email
	user.PhoneNumber = phoneNumber
	return &user
}
