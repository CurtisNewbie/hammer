package hammer

import "testing"

func TestCompressImage(t *testing.T) {
	in := "../test.jpg"
	out := "../compressed.jpg"

	if e := CompressImage(in, out); e != nil {
		t.Fatal(e)
	}
}
