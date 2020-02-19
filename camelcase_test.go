package generator

import "testing"

func TestCamelCaseClassic(t *testing.T) {
	cc, err := CamelCase(nil)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		in, want string
	}{
		{"one", "One"},
		{"one_two", "OneTwo"},
		{"_my_field_name_2", "XMyFieldName_2"},
		{"Something_Capped", "Something_Capped"},
		{"my_Name", "My_Name"},
		{"OneTwo", "OneTwo"},
		{"_", "X"},
		{"_a_", "XA_"},
	}
	for _, tc := range tests {
		if got := cc(tc.in); got != tc.want {
			t.Errorf("CamelCase(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestCamelCaseImproved(t *testing.T) {
	cc, err := CamelCase([]string{"id"})
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		in, want string
	}{
		{"id", "ID"},
		{"new_id", "NewID"},
		{"some_id_middle", "SomeIDMiddle"},
		{"id_begins", "IDBegins"},
	}
	for _, tc := range tests {
		if got := cc(tc.in); got != tc.want {
			t.Errorf("CamelCase(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
