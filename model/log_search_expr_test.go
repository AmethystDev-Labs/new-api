package model

import "testing"

func TestParseLogSearchInput_LegacyInput(t *testing.T) {
	conditions, err := parseLogSearchInput("plain request", "request_id")
	if err != nil {
		t.Fatalf("parseLogSearchInput failed: %v", err)
	}
	if len(conditions) != 1 {
		t.Fatalf("expected 1 condition, got %d", len(conditions))
	}
	if conditions[0].field != "request_id" || conditions[0].useLike {
		t.Fatalf("unexpected condition: %+v", conditions[0])
	}
	if conditions[0].value != "plain request" {
		t.Fatalf("unexpected value: %s", conditions[0].value)
	}
}

func TestParseLogSearchInput_Expression(t *testing.T) {
	conditions, err := parseLogSearchInput("request_id:req-* username:ep-*", "request_id")
	if err != nil {
		t.Fatalf("parseLogSearchInput failed: %v", err)
	}
	if len(conditions) != 2 {
		t.Fatalf("expected 2 conditions, got %d", len(conditions))
	}
	if conditions[0].field != "request_id" || !conditions[0].useLike || conditions[0].value != "req-%" {
		t.Fatalf("unexpected first condition: %+v", conditions[0])
	}
	if conditions[1].field != "username" || !conditions[1].useLike || conditions[1].value != "ep-%" {
		t.Fatalf("unexpected second condition: %+v", conditions[1])
	}
}

func TestBuildLogSearchPattern_Escaping(t *testing.T) {
	pattern, useLike, err := buildLogSearchPattern("abc!_x*")
	if err != nil {
		t.Fatalf("buildLogSearchPattern failed: %v", err)
	}
	if !useLike {
		t.Fatalf("expected useLike=true")
	}
	if pattern != "abc!!!_x%" {
		t.Fatalf("unexpected pattern: %s", pattern)
	}
}

