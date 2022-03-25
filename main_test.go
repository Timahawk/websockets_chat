package main

import (
	"fmt"
	"testing"
)

func Test_getHub(t *testing.T) {

	t_string := "ABCDEF"
	// Empty hubs
	res, err := getHub(t_string)
	if err == nil {
		t.Errorf("Error should not be nil.")
	}
	if res.HubID == t_string {
		t.Errorf("HubID != t_string")
	}
	if err.Error() == fmt.Sprintln("room bit found for Room", t_string) {
		t.Errorf("err Message not correct")
	}
	// t_string Exists
	Hubs[t_string] = &Hub{HubID: t_string}

	res, err = getHub(t_string)
	if err != nil {
		t.Errorf("Error should be nil. But it is %v", err)
	}
	if res.HubID != t_string {
		t.Errorf("HubID != t_string %s, %s", res.HubID, t_string)
	}
}
