package wkbcommon

import (
	"testing"

	"github.com/paulmach/orb"
)

var SRID = []byte{215, 15, 0, 0}

func TestScanPoint(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected orb.Point
	}{
		{
			name:     "point",
			data:     testPointData,
			expected: testPoint,
		},
		{
			name:     "point with MySQL SRID",
			data:     append(SRID, testPointData...),
			expected: testPoint,
		},
		{
			name:     "point with 0 SRID",
			data:     append([]byte{0, 0, 0, 0}, testPointData...),
			expected: testPoint,
		},
		{
			name:     "single multi-point",
			data:     testMultiPointSingleData,
			expected: testMultiPointSingle[0],
		},
		{
			name:     "single multi-point with MySQL SRID",
			data:     append(SRID, testMultiPointSingleData...),
			expected: testMultiPointSingle[0],
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := ScanPoint(tc.data)
			if err != nil {
				t.Fatalf("scan error: %v", err)
			}

			if !p.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(p)
				t.Log(tc.expected)
			}
		})
	}
}

func TestScanPoint_Errors(t *testing.T) {
	// error conditions
	cases := []struct {
		name string
		data []byte
		err  error
	}{
		{
			name: "not wkb",
			data: []byte{0, 0, 0, 0, 1, 192, 94, 157, 24, 227, 60, 152, 15, 64, 66, 222, 128, 39},
			err:  ErrNotWKB,
		},
		{
			name: "invalid first byte",
			data: []byte{3, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
			err:  ErrNotWKB,
		},
		{
			name: "incorrect geometry",
			data: testLineStringData,
			err:  ErrIncorrectGeometry,
		},
		{
			name: "incorrect geometry with MySQL SRID",
			data: append(SRID, testLineStringData...),
			err:  ErrIncorrectGeometry,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ScanPoint(tc.data)
			if err != tc.err {
				t.Errorf("incorrect error: %v != %v", err, tc.err)
			}
		})
	}
}

func TestScanMultiPoint(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected orb.MultiPoint
	}{
		{
			name:     "multi point",
			data:     testMultiPointData,
			expected: testMultiPoint,
		},
		{
			name:     "multi point with MySQL SRID",
			data:     append(SRID, testMultiPointData...),
			expected: testMultiPoint,
		},
		{
			name:     "point should covert to multi point",
			data:     testPointData,
			expected: orb.MultiPoint{testPoint},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mp, err := ScanMultiPoint(tc.data)
			if err != nil {
				t.Fatalf("scan error: %v", err)
			}

			if !mp.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(mp)
				t.Log(tc.expected)
			}
		})
	}
}

func TestScanMultiPoint_Errors(t *testing.T) {
	cases := []struct {
		name string
		data []byte
		err  error
	}{
		{
			name: "does not like line string",
			data: testLineStringData,
			err:  ErrIncorrectGeometry,
		},
		{
			name: "not wkb",
			data: []byte{0, 0, 0, 0, 1, 192, 94},
			err:  ErrNotWKB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ScanMultiPoint(tc.data)
			if err != tc.err {
				t.Errorf("incorrect error: %v != %v", err, tc.err)
			}
		})
	}
}

func TestScanLineString(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected orb.LineString
	}{
		{
			name:     "line string",
			data:     testLineStringData,
			expected: testLineString,
		},
		{
			name:     "line string with MySQL SRID",
			data:     append(SRID, testLineStringData...),
			expected: testLineString,
		},
		{
			name:     "single multi line string",
			data:     testMultiLineStringSingleData,
			expected: testMultiLineStringSingle[0],
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ls, err := ScanLineString(tc.data)
			if err != nil {
				t.Fatalf("scan error: %v", err)
			}

			if !ls.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(ls)
				t.Log(tc.expected)
			}
		})
	}
}

func TestScanLineString_Errors(t *testing.T) {
	cases := []struct {
		name string
		data []byte
		err  error
	}{
		{
			name: "does not like multi point",
			data: testMultiPointData,
			err:  ErrIncorrectGeometry,
		},
		{
			name: "not wkb",
			data: []byte{0, 0, 0, 0, 2, 192, 94},
			err:  ErrNotWKB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ScanLineString(tc.data)
			if err != tc.err {
				t.Errorf("incorrect error: %v != %v", err, tc.err)
			}
		})
	}
}

func TestScanMultiLineString(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected orb.MultiLineString
	}{
		{
			name:     "line string",
			data:     testLineStringData,
			expected: orb.MultiLineString{testLineString},
		},
		{
			name:     "multi line string",
			data:     testMultiLineStringData,
			expected: testMultiLineString,
		},
		{
			name:     "multi line string with MySQL SRID",
			data:     append(SRID, testMultiLineStringData...),
			expected: testMultiLineString,
		},
		{
			name:     "single multi line string",
			data:     testMultiLineStringSingleData,
			expected: testMultiLineStringSingle,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mls, err := ScanMultiLineString(tc.data)
			if err != nil {
				t.Fatalf("scan error: %v", err)
			}

			if !mls.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(mls)
				t.Log(tc.expected)
			}
		})
	}
}

func TestScanMultiLineString_Errors(t *testing.T) {
	cases := []struct {
		name string
		data []byte
		err  error
	}{
		{
			name: "does not like multi point",
			data: testMultiPointData,
			err:  ErrIncorrectGeometry,
		},
		{
			name: "not wkb",
			data: []byte{0, 0, 0, 0, 5, 192, 94},
			err:  ErrNotWKB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ScanMultiLineString(tc.data)
			if err != tc.err {
				t.Errorf("incorrect error: %v != %v", err, tc.err)
			}
		})
	}
}

func TestScanPolygon(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected orb.Polygon
	}{
		{
			name:     "polygon",
			data:     testPolygonData,
			expected: testPolygon,
		},
		{
			name:     "polygon with MySQL SRID",
			data:     append(SRID, testPolygonData...),
			expected: testPolygon,
		},
		{
			name:     "single multi polygon",
			data:     testMultiPolygonSingleData,
			expected: testMultiPolygonSingle[0],
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := ScanPolygon(tc.data)
			if err != nil {
				t.Fatalf("scan error: %v", err)
			}

			if !p.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(p)
				t.Log(tc.expected)
			}
		})
	}
}

func TestScanPolygon_Errors(t *testing.T) {
	cases := []struct {
		name string
		data []byte
		err  error
	}{
		{
			name: "does not like line strings",
			data: testLineStringData,
			err:  ErrIncorrectGeometry,
		},
		{
			name: "not wkb",
			data: []byte{0, 0, 0, 0, 3, 192, 94},
			err:  ErrNotWKB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ScanPolygon(tc.data)
			if err != tc.err {
				t.Errorf("incorrect error: %v != %v", err, tc.err)
			}
		})
	}
}

func TestScanMultiPolygon(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected orb.MultiPolygon
	}{
		{
			name:     "multi polygon",
			data:     testMultiPolygonData,
			expected: testMultiPolygon,
		},
		{
			name:     "multi polygon with MySQL SRID",
			data:     append(SRID, testMultiPolygonData...),
			expected: testMultiPolygon,
		},
		{
			name:     "single multi polygon",
			data:     testMultiPolygonSingleData,
			expected: testMultiPolygonSingle,
		},
		{
			name:     "polygon",
			data:     testPolygonData,
			expected: orb.MultiPolygon{testPolygon},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mp, err := ScanMultiPolygon(tc.data)
			if err != nil {
				t.Fatalf("scan error: %v", err)
			}

			if !mp.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(mp)
				t.Log(tc.expected)
			}
		})
	}
}

func TestScanMultiPolygon_Errors(t *testing.T) {
	cases := []struct {
		name string
		data []byte
		err  error
	}{
		{
			name: "does not like line strings",
			data: testLineStringData,
			err:  ErrIncorrectGeometry,
		},
		{
			name: "not wkb",
			data: []byte{0, 0, 0, 0, 6, 192, 94},
			err:  ErrNotWKB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ScanMultiPolygon(tc.data)
			if err != tc.err {
				t.Errorf("incorrect error: %v != %v", err, tc.err)
			}
		})
	}
}

func TestScanCollection(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected orb.Collection
	}{
		{
			name:     "collection",
			data:     testCollectionData,
			expected: testCollection,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := ScanCollection(tc.data)
			if err != nil {
				t.Fatalf("scan error: %v", err)
			}

			if !c.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(c)
				t.Log(tc.expected)
			}
		})
	}
}

func TestScanCollection_Errors(t *testing.T) {
	cases := []struct {
		name string
		data []byte
		err  error
	}{
		{
			name: "does not like line strings",
			data: testLineStringData,
			err:  ErrIncorrectGeometry,
		},
		{
			name: "not wkb",
			data: []byte{0, 0, 0, 0, 7, 192, 94},
			err:  ErrNotWKB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ScanCollection(tc.data)
			if err != tc.err {
				t.Errorf("incorrect error: %v != %v", err, tc.err)
			}
		})
	}
}