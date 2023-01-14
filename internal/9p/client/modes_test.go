package client

import "testing"

func TestModeStringSimple(t *testing.T) {
	want := "-rwxr-xr-x"
	input := uint32(0755)
	got := ModeString(input)
	if got != want {
		t.Errorf("got ModeString(%#o) = %q, want %q", input, got, want)
	}
}

func TestModeString(t *testing.T) {
	for mode, want := range map[uint32]string{
		ModeDir:        "d---------",
		ModeUserRead:   "-r--------",
		ModeUserWrite:  "--w-------",
		ModeUserExec:   "---x------",
		ModeGroupRead:  "----r-----",
		ModeGroupWrite: "-----w----",
		ModeGroupExec:  "------x---",
		ModeOtherRead:  "-------r--",
		ModeOtherWrite: "--------w-",
		ModeOtherExec:  "---------x",
	} {
		got := ModeString(mode)
		if got != want {
			t.Errorf("bad mode string for mode %#o, got: %v, want: %v", mode, got, want)
		}
	}
}
