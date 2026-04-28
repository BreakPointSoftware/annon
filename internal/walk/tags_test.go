package walk

import "testing"

func TestParseTag(t *testing.T) {
	if got := parseTag("remove"); !got.remove { t.Fatalf("unexpected tag: %+v", got) }
}
