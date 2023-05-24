package commands

import "testing"

type test struct {
	item string
	want bool
}

var iins = []test{
	{"12354", false},
	{"840629300612", false},
	{"840629300619", true},
	{"971240001315", true},
}

func Test_okBinIin(t *testing.T) {
	for _, tst := range iins {
		t.Run(tst.item, func(t *testing.T) {
			if okBinIin(tst.item) != tst.want {
				t.Errorf("error chec iin/bin %v, expected %v", tst.item, tst.want)
			}
		})
	}
}
