package udm

// existing_test.go - Tests that pass on the buggy version.
// These tests verify UDM SDM behavior when NSSAI data exists.

import (
	"testing"
)

// SubscriptionData represents subscriber data with NSSAI info.
type SubscriptionData struct {
	Supi     string
	NssaiMap map[string]NssaiData
}

// NssaiData represents NSSAI-specific subscription data.
type NssaiData struct {
	Snssai string
	Dnn    string
	Qos    string
}

// TestGetExistingNssaiData verifies retrieval of existing NSSAI data.
func TestGetExistingNssaiData(t *testing.T) {
	subData := SubscriptionData{
		Supi: "imsi-208930000000001",
		NssaiMap: map[string]NssaiData{
			"1-010203": {Snssai: "1-010203", Dnn: "internet", Qos: "9"},
		},
	}

	nssai, ok := subData.NssaiMap["1-010203"]
	if !ok {
		t.Fatal("expected NSSAI data to exist")
	}
	if nssai.Dnn != "internet" {
		t.Fatalf("unexpected DNN: %s", nssai.Dnn)
	}
}

// TestSubscriptionDataCreation verifies basic subscription data creation.
func TestSubscriptionDataCreation(t *testing.T) {
	subData := SubscriptionData{
		Supi:     "imsi-208930000000001",
		NssaiMap: make(map[string]NssaiData),
	}
	if subData.Supi != "imsi-208930000000001" {
		t.Fatal("unexpected SUPI")
	}
}
