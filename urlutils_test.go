package main

import (
	"testing"
)

func TestGetRequestURL(t *testing.T) {
	if _, err := getRequestURL("http://management.azure.com/subscriptions?api-version=2015-01-01"); err == nil {
		t.Error("Should not accept http endpoint")
		t.Fail()
	}

	if _, err := getRequestURL("https://random.azure.com/subscriptions?api-version=2015-01-01"); err == nil {
		t.Error("Should not accept non-ARM endpoint")
		t.Fail()
	}

	if _, err := getRequestURL("https://management.azure.com/?api-version=2015-01-01"); err == nil {
		t.Error("Should validate request path")
		t.Fail()
	}

	if _, err := getRequestURL("https://westus.management.azure.com/subscriptions?api-version=2015-01-01"); err != nil {
		t.Error("Should accept westus.management.azure.com", err)
		t.Fail()
	}
}

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
