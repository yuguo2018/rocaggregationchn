package core

import "testing"

func TestPasswordValidation(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		pw         string
		shouldFail bool
	}{
		{"test", true},
		{"testtest\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98", true},
		{"placeOfInterest⌘", true},
		{"password\nwith\nlinebreak", true},
		{"password\twith\vtabs", true},
		// Ok passwords
		{"password WhichIsOk", false},
		{"passwordOk!@#$%^&*()", false},
		{"12301203123012301230123012", false},
	}
	for _, test := range testcases {
		err := ValidatePasswordFormat(test.pw)
		if err == nil && test.shouldFail {
			t.Errorf("password '%v' should fail validation", test.pw)
		} else if err != nil && !test.shouldFail {
			t.Errorf("password '%v' shound not fail validation, but did: %v", test.pw, err)
		}
	}
}
