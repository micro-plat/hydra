package aigw

import "testing"

func TestNew_Defaults(t *testing.T) {
	s := New(DefaultAddress)
	if s.GetAddress() != DefaultAddress {
		t.Fatalf("address = %s, want %s", s.GetAddress(), DefaultAddress)
	}
	if s.GetRTimeout() != DefaultRTimeOut {
		t.Fatalf("rTimeout = %d, want %d", s.GetRTimeout(), DefaultRTimeOut)
	}
	if s.GetWTimeout() != DefaultWTimeOut {
		t.Fatalf("wTimeout = %d, want %d", s.GetWTimeout(), DefaultWTimeOut)
	}
	if s.GetRHTimeout() != DefaultRHTimeOut {
		t.Fatalf("rhTimeout = %d, want %d", s.GetRHTimeout(), DefaultRHTimeOut)
	}
	if s.GetStreamTimeout() != DefaultStreamTimeOut {
		t.Fatalf("streamTimeout = %d, want %d", s.GetStreamTimeout(), DefaultStreamTimeOut)
	}
}

func TestNew_WithOptions(t *testing.T) {
	s := New("9082", WithTimeout(10, 20, 30), WithStreamTimeout(40), WithTrace())
	if s.GetAddress() != "9082" {
		t.Fatalf("address = %s, want 9082", s.GetAddress())
	}
	if s.GetRTimeout() != 10 || s.GetWTimeout() != 20 || s.GetRHTimeout() != 30 {
		t.Fatalf("timeout = %d/%d/%d, want 10/20/30", s.GetRTimeout(), s.GetWTimeout(), s.GetRHTimeout())
	}
	if s.GetStreamTimeout() != 40 {
		t.Fatalf("streamTimeout = %d, want 40", s.GetStreamTimeout())
	}
	if !s.Trace {
		t.Fatal("trace = false, want true")
	}
}

func TestNew_WithDisable(t *testing.T) {
	s := New(DefaultAddress, WithDisable())
	if s.Status != StartStop {
		t.Fatalf("status = %s, want %s", s.Status, StartStop)
	}
}
