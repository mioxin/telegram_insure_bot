package config

import (
	"strings"
	"testing"
)

type test struct {
	conf string
	user string
	want bool
}

var tests = []test{
	{"deny =all\n	allow =all\n", "user", false},
	{"deny = all\n	allow =user1,user\n", "user", true},
	{"deny =all\n	allow =user1,user\n", "user1", true},
	{"deny =all\n	allow =user1\n", "user", false},
	{"deny =user\n	allow =all\n", "user", false},
	{"deny =user1,user2\n	allow =all\n", "user", true},
	{"deny =user1,user2\n	allow =all\n", "user2", false},
	{"deny =all\n", "user", false},
	{"allow =user\n", "user", true},
	{"allow =all\n", "user", true},
}

func Test_isAccess(t *testing.T) {
	for _, tst := range tests {
		if cfg, err := NewConfig(strings.NewReader(tst.conf)); err != nil {
			t.Error(err)
		} else {
			t.Run(tst.conf, func(t *testing.T) {
				if res := cfg.IsAccess(tst.user); res != tst.want {
					t.Errorf("on user \"%s\" get %v want %v", tst.user, res, tst.want)
				}
			})
		}
	}
}
