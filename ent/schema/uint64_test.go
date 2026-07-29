package schema

import "testing"

func TestUint64ValueScanner(t *testing.T) {
	maxUint64 := ^uint64(0)

	driverValue, err := uint64ValueScanner.Value(maxUint64)
	if err != nil {
		t.Fatalf("Value() failed: %v", err)
	}
	if got, want := driverValue, "18446744073709551615"; got != want {
		t.Fatalf("Value() = %v, want %v", got, want)
	}

	scanValue := uint64ValueScanner.ScanValue()
	if err := scanValue.Scan(driverValue); err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}
	got, err := uint64ValueScanner.FromValue(scanValue)
	if err != nil {
		t.Fatalf("FromValue() failed: %v", err)
	}
	if got != maxUint64 {
		t.Fatalf("FromValue() = %d, want %d", got, maxUint64)
	}
}
