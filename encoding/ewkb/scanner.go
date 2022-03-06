package ewkb

import (
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/internal/wkbcommon"
)

var (
	_ sql.Scanner  = &GeometryScanner{}
	_ driver.Value = value{}
)

// GeometryScanner is a thing that can scan in sql query results.
// It can be used as a scan destination:
//
//	var s wkb.GeometryScanner
//	err := db.QueryRow("SELECT latlon FROM foo WHERE id=?", id).Scan(&s)
//	...
//	if s.Valid {
//	  // use s.Geometry
//	} else {
//	  // NULL value
//	}
type GeometryScanner struct {
	g        interface{}
	SRID     int
	Geometry orb.Geometry
	Valid    bool // Valid is true if the geometry is not NULL
}

// Scanner will return a GeometryScanner that can scan sql query results.
// The geometryScanner.Geometry attribute will be set to the value.
// If g is non-nil, it MUST be a pointer to an orb.Geometry
// type like a Point or LineString. In that case the value will be written to
// g and the Geometry attribute.
//
//	var p orb.Point
//	err := db.QueryRow("SELECT latlon FROM foo WHERE id=?", id).Scan(wkb.Scanner(&p))
//	...
//	// use p
//
// If the value may be null check Valid first:
//
//	var point orb.Point
//	s := wkb.Scanner(&point)
//	err := db.QueryRow("SELECT latlon FROM foo WHERE id=?", id).Scan(&s)
//	...
//	if s.Valid {
//	  // use p
//	} else {
//	  // NULL value
//	}
//
// Scanning directly from MySQL columns is supported. By default MySQL returns geometry
// data as WKB but prefixed with a 4 byte SRID. To support this, if the data is not
// valid WKB, the code will strip the first 4 bytes and try again.
// This works for most use cases.
func Scanner(g interface{}) *GeometryScanner {
	return &GeometryScanner{g: g}
}

// Scan will scan the input []byte data into a geometry.
// This could be into the orb geometry type pointer or, if nil,
// the scanner.Geometry attribute.
func (s *GeometryScanner) Scan(d interface{}) error {
	s.Geometry = nil
	s.Valid = false

	if d == nil {
		return nil
	}

	data, ok := d.([]byte)
	if !ok {
		return ErrUnsupportedDataType
	}

	if data == nil {
		return nil
	}

	// go-pg will return ST_AsBinary(*) data as `\xhexencoded` which
	// needs to be converted to true binary for further decoding.
	// Code detects the \x prefix and then converts the rest from Hex to binary.
	if len(data) > 2 && data[0] == byte('\\') && data[1] == byte('x') {
		n, err := hex.Decode(data, data[2:])
		if err != nil {
			return fmt.Errorf("thought the data was hex, but it is not: %v", err)
		}
		data = data[:n]
	}

	switch g := s.g.(type) {
	case nil:
		m, srid, err := Unmarshal(data)
		if err != nil {
			return err
		}

		s.SRID = srid
		s.Geometry = m
		s.Valid = true
		return nil
	case *orb.Point:
		p, srid, err := wkbcommon.ScanPoint(data)
		if err != nil {
			return mapCommonError(err)
		}

		*g = p
		s.SRID = srid
		s.Geometry = p
		s.Valid = true
		return nil
	case *orb.MultiPoint:
		m, srid, err := wkbcommon.ScanMultiPoint(data)
		if err != nil {
			return mapCommonError(err)
		}

		*g = m
		s.SRID = srid
		s.Geometry = m
		s.Valid = true
		return nil
	case *orb.LineString:
		l, srid, err := wkbcommon.ScanLineString(data)
		if err != nil {
			return mapCommonError(err)
		}

		*g = l
		s.SRID = srid
		s.Geometry = l
		s.Valid = true
		return nil
	case *orb.MultiLineString:
		m, srid, err := wkbcommon.ScanMultiLineString(data)
		if err != nil {
			return mapCommonError(err)
		}

		*g = m
		s.SRID = srid
		s.Geometry = m
		s.Valid = true
		return nil
	case *orb.Ring:
		m, srid, err := Unmarshal(data)
		if err != nil {
			return err
		}

		if p, ok := m.(orb.Polygon); ok && len(p) == 1 {
			*g = p[0]
			s.SRID = srid
			s.Geometry = p[0]
			s.Valid = true
			return nil
		}

		return ErrIncorrectGeometry
	case *orb.Polygon:
		p, srid, err := wkbcommon.ScanPolygon(data)
		if err != nil {
			return mapCommonError(err)
		}

		*g = p
		s.SRID = srid
		s.Geometry = p
		s.Valid = true
		return nil
	case *orb.MultiPolygon:
		m, srid, err := wkbcommon.ScanMultiPolygon(data)
		if err != nil {
			return mapCommonError(err)
		}

		*g = m
		s.SRID = srid
		s.Geometry = m
		s.Valid = true
		return nil
	case *orb.Collection:
		c, srid, err := wkbcommon.ScanCollection(data)
		if err != nil {
			return mapCommonError(err)
		}

		*g = c
		s.SRID = srid
		s.Geometry = c
		s.Valid = true
		return nil
	case *orb.Bound:
		m, srid, err := Unmarshal(data)
		if err != nil {
			return err
		}

		b := m.Bound()
		*g = b
		s.SRID = srid
		s.Geometry = b
		s.Valid = true
		return nil
	}

	return ErrIncorrectGeometry
}

type value struct {
	srid int
	v    orb.Geometry
}

// Value will create a driver.Valuer that will WKB the geometry
// into the database query.
func Value(g orb.Geometry, srid int) driver.Valuer {
	return value{srid: srid, v: g}

}

func (v value) Value() (driver.Value, error) {
	val, err := Marshal(v.v, v.srid)
	if val == nil {
		return nil, err
	}
	return val, err
}
