package main

import (
	"net/netip"
	"os"
	"strings"

	"github.com/sagernet/sing-box/common/geoip"
	"github.com/sagernet/sing-box/log"
	E "github.com/sagernet/sing/common/exceptions"
	N "github.com/sagernet/sing/common/network"

	"github.com/spf13/cobra"
)

var commandGeoipInfo = &cobra.Command{
	Use:   "info <address>",
	Short: "Lookup country, timezone, and locale for an IP address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := geoipInfo(args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	commandGeoip.AddCommand(commandGeoipInfo)
}

func geoipInfo(address string) error {
	addr, err := netip.ParseAddr(address)
	if err != nil {
		return E.Cause(err, "parse address")
	}
	if !N.IsPublicAddr(addr) {
		os.Stdout.WriteString("private\tUTC\ten-US\n")
		return nil
	}
	if geoipReader.Metadata.DatabaseType == "sing-geoip" {
		var code string
		_ = geoipReader.Lookup(addr.AsSlice(), &code)
		if code != "" {
			writeGeoipInfo(code, geoip.CountryDefaultTimezone(code), geoip.CountryLocale(code))
			return nil
		}
	} else {
		var record struct {
			Country struct {
				ISOCode string `maxminddb:"iso_code"`
			} `maxminddb:"country"`
			Location struct {
				TimeZone string `maxminddb:"time_zone"`
			} `maxminddb:"location"`
		}
		_ = geoipReader.Lookup(addr.AsSlice(), &record)
		if record.Country.ISOCode != "" {
			code := strings.ToLower(record.Country.ISOCode)
			timezone := record.Location.TimeZone
			if timezone == "" {
				timezone = geoip.CountryDefaultTimezone(code)
			}
			writeGeoipInfo(code, timezone, geoip.CountryLocale(code))
			return nil
		}
	}
	os.Stdout.WriteString("unknown\tUTC\ten-US\n")
	return nil
}

func writeGeoipInfo(country string, timezone string, locale string) {
	os.Stdout.WriteString(country + "\t" + timezone + "\t" + locale + "\n")
}
