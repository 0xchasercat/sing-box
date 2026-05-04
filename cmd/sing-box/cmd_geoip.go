package main

import (
	"strings"

	"github.com/sagernet/sing-box/log"
	E "github.com/sagernet/sing/common/exceptions"

	"github.com/oschwald/maxminddb-golang"
	"github.com/spf13/cobra"
)

var (
	geoipReader          *maxminddb.Reader
	commandGeoIPFlagFile string
)

var commandGeoip = &cobra.Command{
	Use:   "geoip",
	Short: "GeoIP tools",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := geoipPreRun()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	commandGeoip.PersistentFlags().StringVarP(&commandGeoIPFlagFile, "file", "f", "geoip.db", "geoip file")
	mainCommand.AddCommand(commandGeoip)
}

func geoipPreRun() error {
	reader, err := maxminddb.Open(commandGeoIPFlagFile)
	if err != nil {
		return err
	}
	dbType := reader.Metadata.DatabaseType
	if dbType != "sing-geoip" && !strings.HasPrefix(dbType, "GeoLite2") && !strings.HasPrefix(dbType, "GeoIP2") {
		reader.Close()
		return E.New("unsupported database type: ", dbType, " (expected sing-geoip, GeoLite2-*, or GeoIP2-*)")
	}
	geoipReader = reader
	return nil
}
