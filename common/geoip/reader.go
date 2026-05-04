package geoip

import (
	"net/netip"
	"strings"

	E "github.com/sagernet/sing/common/exceptions"

	"github.com/oschwald/maxminddb-golang"
)

type Reader struct {
	reader *maxminddb.Reader
	dbType string
}

func Open(path string) (*Reader, []string, error) {
	database, err := maxminddb.Open(path)
	if err != nil {
		return nil, nil, err
	}
	dbType := database.Metadata.DatabaseType
	if dbType != "sing-geoip" && !strings.HasPrefix(dbType, "GeoLite2") && !strings.HasPrefix(dbType, "GeoIP2") {
		database.Close()
		return nil, nil, E.New("unsupported database type: ", dbType, " (expected sing-geoip, GeoLite2-*, or GeoIP2-*)")
	}
	return &Reader{database, dbType}, database.Metadata.Languages, nil
}

func (r *Reader) Lookup(addr netip.Addr) string {
	var code string
	_ = r.reader.Lookup(addr.AsSlice(), &code)
	if code != "" {
		return code
	}
	return "unknown"
}

func (r *Reader) Close() error {
	return r.reader.Close()
}
