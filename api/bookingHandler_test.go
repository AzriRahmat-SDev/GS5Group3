package api

import "testing"

func TestStartDateIsBeforeEndDate(t *testing.T) {
	result := startDateIsBeforeEndDate("2006-01-02", "2006-01-03")
	if result == false {
		t.Errorf("Result was incorrect, got: %v, want %v", result, true)
	}

	result = startDateIsBeforeEndDate("2006-01-03", "2006-01-02")
	if result == true {
		t.Errorf("Result was incorrect, got: %v, want %v", result, false)
	}

	result = startDateIsBeforeEndDate("2006-01-02", "2006-01-02")
	if result == true {
		t.Errorf("Result was incorrect, got: %v, want %v", result, false)
	}
}
