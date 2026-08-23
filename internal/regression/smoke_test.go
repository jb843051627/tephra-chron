package regression

import "testing"

func TestPackageLoads(t *testing.T) {
	if 2+2 != 4 {
		t.Fatal("arithmetic")
	}
}
