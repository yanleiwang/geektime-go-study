package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCamelToUnderline(t *testing.T) {
	testCases := []struct {
		name    string
		srcStr  string
		wantStr string
	}{
		{
			name:    "Abbreviation",
			srcStr:  "ID",
			wantStr: "id",
		},
		{
			name:    "use number",
			srcStr:  "Table1Name",
			wantStr: "table1_name",
		},
		// 我们这些用例就是为了确保
		// 在忘记 CamelToUnderline 的行为特性之后
		// 可以从这里找回来
		// 比如 userPWDuser -> user_pwduser 而不是 user_pwd_user
		{
			name:    "test 3",
			srcStr:  "userPWD",
			wantStr: "user_pwd",
		},
		{
			name:    "test 4",
			srcStr:  "userPWDuser",
			wantStr: "user_pwduser",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := CamelToUnderline(tc.srcStr)
			assert.Equal(t, tc.wantStr, res)
		})
	}
}
