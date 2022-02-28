package ewkb

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/internal/wkbcommon"
)

func TestMarshal(t *testing.T) {
	for _, g := range orb.AllGeometries {
		Marshal(g, 0, binary.BigEndian)
	}
}

func TestMustMarshal(t *testing.T) {
	for _, g := range orb.AllGeometries {
		MustMarshal(g, 0, binary.BigEndian)
	}
}

func BenchmarkEncode_Point(b *testing.B) {
	g := orb.Point{1, 2}
	e := NewEncoder(ioutil.Discard)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Encode(g)
	}
}

func BenchmarkEncode_LineString(b *testing.B) {
	g := orb.LineString{
		{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5},
		{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5},
	}
	e := NewEncoder(ioutil.Discard)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Encode(g)
	}
}

func compare(t testing.TB, e orb.Geometry, b []byte) {
	t.Helper()

	// Decoder
	g, _, err := NewDecoder(bytes.NewReader(b)).Decode() // TODO
	if err != nil {
		t.Fatalf("decoder: read error: %v", err)
	}

	if !orb.Equal(g, e) {
		t.Errorf("decoder: incorrect geometry: %v != %v", g, e)
	}

	// Umarshal
	g, _, err = Unmarshal(b) // TODO
	if err != nil {
		t.Fatalf("unmarshal: read error: %v", err)
	}

	if !orb.Equal(g, e) {
		t.Errorf("unmarshal: incorrect geometry: %v != %v", g, e)
	}

	var data []byte
	if b[0] == 0 {
		data, err = Marshal(g, 0, binary.BigEndian) // TODO
	} else {
		data, err = Marshal(g, 0, binary.LittleEndian) // TODO
	}
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if !bytes.Equal(data, b) {
		t.Logf("%v", data)
		t.Logf("%v", b)
		t.Errorf("marshal: incorrent encoding")
	}

	// preallocation
	if len(data) != wkbcommon.GeomLength(e) {
		t.Errorf("prealloc length: %v != %v", len(data), wkbcommon.GeomLength(e))
	}
}
