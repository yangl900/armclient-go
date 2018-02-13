package main

import (
	"testing"
)

func TestArmUrlPath(t *testing.T) {
	if isArmURLPath("") {
		t.Fail()
	}

	if isArmURLPath("/") {
		t.Fail()
	}

	if isArmURLPath("/foo") {
		t.Fail()
	}

	if !isArmURLPath("/subscriptions?api-version=2015-01-01") {
		t.Error("Failed to match subscriptions request.")
		t.Fail()
	}

	if !isArmURLPath("/SUBSCRIPTIONS?api-version=2015-01-01") {
		t.Error("Failed to match subscriptions request.")
		t.Fail()
	}

	if !isArmURLPath("/subscriptions/12345678/resourceGroups?api-version=2015-01-01") {
		t.Error("Failed to match subscriptions request.")
		t.Fail()
	}

	if !isArmURLPath("/tenants?api-version=2015-01-01") {
		t.Error("Failed to match subscriptions request.")
		t.Fail()
	}

	if !isArmURLPath("/providers?api-version=2015-01-01") {
		t.Error("Failed to match subscriptions request.")
		t.Fail()
	}
}
