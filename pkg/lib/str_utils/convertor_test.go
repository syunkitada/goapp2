package str_utils

import "testing"

func TestConvertToCamelFormat(t *testing.T) {
	got := ConvertToCamelFormat("hoge.piyo")
	want := "HogePiyo"
	if got != want {
		t.Errorf("got=%s, want=%s", got, want)
	}
}

func TestConvertToLowerFormat(t *testing.T) {
	got := ConvertToLowerFormat("HogePiyo")
	want := "hoge.piyo"
	if got != want {
		t.Errorf("got=%s, want=%s", got, want)
	}
}

func TestSplitActionDataName(t *testing.T) {
	gotAction, gotData := SplitActionDataName("CreateData")
	wantAction := "Create"
	wantData := "Data"
	if gotAction != wantAction {
		t.Errorf("got=%s, want=%s", gotAction, wantAction)
	}

	if gotData != wantData {
		t.Errorf("got=%s, want=%s", gotData, wantData)
	}
}

func TestSplitSpace(t *testing.T) {
	// TODO
}
