package main

import (
	"testing"
)

func TestReadNotInitilizedSetting(t *testing.T) {
	setDefaultSettingsPath("./does-not-exist")
	s, err := readSettings()

	if err != nil {
		t.Error("Failed to read settings: ", err)
		t.Fail()
	}

	if s.ActiveTenant != "" {
		t.Error("Unexpected settings: ", s)
		t.Fail()
	}
}

func TestSettingsSaveAndRead(t *testing.T) {
	setDefaultSettingsPath("./test")

	activeTenant := "foo"
	s := settings{
		ActiveTenant: activeTenant,
	}

	err := saveSettings(s)
	if err != nil {
		t.Error("Failed to save settings: ", err)
		t.Fail()
	}

	readback, err := readSettings()
	if err != nil {
		t.Error("Failed to read settings: ", err)
		t.Fail()
	}

	if readback.ActiveTenant != activeTenant {
		t.Errorf("Active tenant from settings: %s doesn't equal expected: %s", readback.ActiveTenant, activeTenant)
		t.Fail()
	}
}
