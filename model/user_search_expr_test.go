package model

import "testing"

func TestParseUserSearchKeyword_DefaultKeyword(t *testing.T) {
	conditions, err := parseUserSearchKeyword("ep")
	if err != nil {
		t.Fatalf("parseUserSearchKeyword failed: %v", err)
	}
	if len(conditions) != 1 {
		t.Fatalf("expected 1 condition, got %d", len(conditions))
	}
	if conditions[0].id != nil {
		t.Fatalf("expected nil id for non-numeric keyword")
	}
	if conditions[0].pattern != "%ep%" {
		t.Fatalf("expected pattern %%ep%%, got %s", conditions[0].pattern)
	}
}

func TestParseUserSearchKeyword_WildcardAndField(t *testing.T) {
	conditions, err := parseUserSearchKeyword("username:ep-* group:default")
	if err != nil {
		t.Fatalf("parseUserSearchKeyword failed: %v", err)
	}
	if len(conditions) != 2 {
		t.Fatalf("expected 2 conditions, got %d", len(conditions))
	}
	if conditions[0].field != "username" || conditions[0].pattern != "ep-%" {
		t.Fatalf("unexpected username condition: %+v", conditions[0])
	}
	if conditions[1].field != "group" || conditions[1].pattern != "%default%" {
		t.Fatalf("unexpected group condition: %+v", conditions[1])
	}
}

func TestParseUserSearchKeyword_KeepLegacySpaceSearch(t *testing.T) {
	conditions, err := parseUserSearchKeyword("ep user")
	if err != nil {
		t.Fatalf("parseUserSearchKeyword failed: %v", err)
	}
	if len(conditions) != 1 {
		t.Fatalf("expected 1 condition for legacy search, got %d", len(conditions))
	}
	if conditions[0].pattern != "%ep user%" {
		t.Fatalf("unexpected pattern: %s", conditions[0].pattern)
	}
}

func TestBuildUserSearchLikePattern_Escape(t *testing.T) {
	pattern, hasWildcard := buildUserSearchLikePattern("abc%_\\*")
	if !hasWildcard {
		t.Fatalf("expected wildcard=true")
	}
	if pattern != "abc\\%\\_\\\\%" {
		t.Fatalf("unexpected pattern: %s", pattern)
	}
}
