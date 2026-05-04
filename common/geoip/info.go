package geoip

import (
	"net/netip"
	"strings"
)

// GeoInfo holds per-IP geographic metadata for use in browser automation.
type GeoInfo struct {
	Country  string // ISO 3166-1 alpha-2 lowercase, e.g. "us"
	TimeZone string // IANA timezone, e.g. "America/New_York"
	Locale   string // BCP 47 locale tag, e.g. "en-US"
}

// cityRecord matches the GeoLite2-City / GeoIP2-City MMDB record layout.
type cityRecord struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
	Location struct {
		TimeZone string `maxminddb:"time_zone"`
	} `maxminddb:"location"`
}

// LookupInfo returns locale and timezone for addr in addition to country.
// For GeoLite2/GeoIP2 city databases the timezone is read directly from the
// MMDB record; for sing-geoip databases it is derived from a country map.
func (r *Reader) LookupInfo(addr netip.Addr) GeoInfo {
	if r.dbType != "sing-geoip" {
		var record cityRecord
		if err := r.reader.Lookup(addr.AsSlice(), &record); err == nil && record.Country.ISOCode != "" {
			country := strings.ToLower(record.Country.ISOCode)
			tz := record.Location.TimeZone
			if tz == "" {
				tz = countryDefaultTimezone(country)
			}
			return GeoInfo{
				Country:  country,
				TimeZone: tz,
				Locale:   countryLocale(country),
			}
		}
	}

	// sing-geoip: record is a plain country-code string.
	var code string
	_ = r.reader.Lookup(addr.AsSlice(), &code)
	if code == "" {
		return GeoInfo{Country: "unknown", TimeZone: "UTC", Locale: "en-US"}
	}
	return GeoInfo{
		Country:  code,
		TimeZone: countryDefaultTimezone(code),
		Locale:   countryLocale(code),
	}
}

// countryLocale returns the primary BCP 47 locale tag for a country code.
func countryLocale(country string) string {
	if locale, ok := localeMap[strings.ToLower(country)]; ok {
		return locale
	}
	return "en-US"
}

// countryDefaultTimezone returns the primary IANA timezone for a country code.
// Used as a fallback when the MMDB record has no timezone field.
func countryDefaultTimezone(country string) string {
	if tz, ok := timezoneMap[strings.ToLower(country)]; ok {
		return tz
	}
	return "UTC"
}

var localeMap = map[string]string{
	"us": "en-US", "gb": "en-GB", "au": "en-AU", "ca": "en-CA",
	"nz": "en-NZ", "ie": "en-IE", "in": "en-IN", "za": "en-ZA",
	"sg": "en-SG", "ph": "en-PH", "ng": "en-NG", "gh": "en-GH",
	"ke": "en-KE", "ug": "en-UG", "zw": "en-ZW", "bw": "en-BW",
	"de": "de-DE", "at": "de-AT", "ch": "de-CH", "li": "de-LI",
	"fr": "fr-FR", "be": "fr-BE", "lu": "fr-LU", "mc": "fr-MC",
	"es": "es-ES", "mx": "es-MX", "ar": "es-AR", "co": "es-CO",
	"cl": "es-CL", "pe": "es-PE", "ve": "es-VE", "ec": "es-EC",
	"bo": "es-BO", "py": "es-PY", "uy": "es-UY", "cr": "es-CR",
	"gt": "es-GT", "hn": "es-HN", "sv": "es-SV", "ni": "es-NI",
	"pa": "es-PA", "do": "es-DO", "cu": "es-CU",
	"pt": "pt-PT", "br": "pt-BR", "ao": "pt-AO", "mz": "pt-MZ",
	"it": "it-IT", "sm": "it-SM", "va": "it-VA",
	"nl": "nl-NL", "sr": "nl-SR",
	"ru": "ru-RU", "by": "be-BY", "kz": "kk-KZ", "uz": "uz-UZ",
	"ua": "uk-UA",
	"cn": "zh-CN", "tw": "zh-TW", "hk": "zh-HK", "mo": "zh-MO",
	"jp": "ja-JP",
	"kr": "ko-KR",
	"pl": "pl-PL",
	"cz": "cs-CZ", "sk": "sk-SK",
	"hu": "hu-HU",
	"ro": "ro-RO",
	"bg": "bg-BG",
	"hr": "hr-HR", "si": "sl-SI", "rs": "sr-RS", "ba": "bs-BA",
	"se": "sv-SE",
	"no": "nb-NO",
	"dk": "da-DK",
	"fi": "fi-FI",
	"gr": "el-GR", "cy": "el-CY",
	"tr": "tr-TR",
	"il": "he-IL",
	"sa": "ar-SA", "ae": "ar-AE", "eg": "ar-EG", "iq": "ar-IQ",
	"jo": "ar-JO", "kw": "ar-KW", "lb": "ar-LB", "ly": "ar-LY",
	"ma": "ar-MA", "om": "ar-OM", "qa": "ar-QA", "sy": "ar-SY",
	"tn": "ar-TN", "ye": "ar-YE",
	"ir": "fa-IR", "af": "fa-AF",
	"th": "th-TH",
	"vn": "vi-VN",
	"id": "id-ID",
	"my": "ms-MY",
	"pk": "ur-PK",
	"bd": "bn-BD",
	"lk": "si-LK",
	"np": "ne-NP",
	"mm": "my-MM",
	"kh": "km-KH",
	"mn": "mn-MN",
}

var timezoneMap = map[string]string{
	"us": "America/New_York",
	"ca": "America/Toronto",
	"mx": "America/Mexico_City",
	"br": "America/Sao_Paulo",
	"ar": "America/Argentina/Buenos_Aires",
	"cl": "America/Santiago",
	"co": "America/Bogota",
	"pe": "America/Lima",
	"ve": "America/Caracas",
	"ec": "America/Guayaquil",
	"bo": "America/La_Paz",
	"py": "America/Asuncion",
	"uy": "America/Montevideo",
	"gt": "America/Guatemala",
	"cr": "America/Costa_Rica",
	"hn": "America/Tegucigalpa",
	"sv": "America/El_Salvador",
	"ni": "America/Managua",
	"pa": "America/Panama",
	"do": "America/Santo_Domingo",
	"cu": "America/Havana",
	"gb": "Europe/London",
	"ie": "Europe/Dublin",
	"pt": "Europe/Lisbon",
	"es": "Europe/Madrid",
	"fr": "Europe/Paris",
	"be": "Europe/Brussels",
	"nl": "Europe/Amsterdam",
	"lu": "Europe/Luxembourg",
	"de": "Europe/Berlin",
	"at": "Europe/Vienna",
	"ch": "Europe/Zurich",
	"li": "Europe/Vaduz",
	"it": "Europe/Rome",
	"mc": "Europe/Monaco",
	"sm": "Europe/San_Marino",
	"va": "Europe/Vatican",
	"dk": "Europe/Copenhagen",
	"se": "Europe/Stockholm",
	"no": "Europe/Oslo",
	"fi": "Europe/Helsinki",
	"ee": "Europe/Tallinn",
	"lv": "Europe/Riga",
	"lt": "Europe/Vilnius",
	"pl": "Europe/Warsaw",
	"cz": "Europe/Prague",
	"sk": "Europe/Bratislava",
	"hu": "Europe/Budapest",
	"si": "Europe/Ljubljana",
	"hr": "Europe/Zagreb",
	"ba": "Europe/Sarajevo",
	"rs": "Europe/Belgrade",
	"me": "Europe/Podgorica",
	"mk": "Europe/Skopje",
	"al": "Europe/Tirane",
	"ro": "Europe/Bucharest",
	"bg": "Europe/Sofia",
	"gr": "Europe/Athens",
	"cy": "Asia/Nicosia",
	"tr": "Europe/Istanbul",
	"ua": "Europe/Kyiv",
	"by": "Europe/Minsk",
	"md": "Europe/Chisinau",
	"ru": "Europe/Moscow",
	"kz": "Asia/Almaty",
	"uz": "Asia/Tashkent",
	"tm": "Asia/Ashgabat",
	"kg": "Asia/Bishkek",
	"tj": "Asia/Dushanbe",
	"il": "Asia/Jerusalem",
	"sa": "Asia/Riyadh",
	"ae": "Asia/Dubai",
	"qa": "Asia/Qatar",
	"kw": "Asia/Kuwait",
	"bh": "Asia/Bahrain",
	"om": "Asia/Muscat",
	"ye": "Asia/Aden",
	"iq": "Asia/Baghdad",
	"ir": "Asia/Tehran",
	"jo": "Asia/Amman",
	"lb": "Asia/Beirut",
	"sy": "Asia/Damascus",
	"af": "Asia/Kabul",
	"pk": "Asia/Karachi",
	"in": "Asia/Kolkata",
	"np": "Asia/Kathmandu",
	"bd": "Asia/Dhaka",
	"lk": "Asia/Colombo",
	"mm": "Asia/Rangoon",
	"th": "Asia/Bangkok",
	"vn": "Asia/Ho_Chi_Minh",
	"kh": "Asia/Phnom_Penh",
	"la": "Asia/Vientiane",
	"my": "Asia/Kuala_Lumpur",
	"sg": "Asia/Singapore",
	"id": "Asia/Jakarta",
	"ph": "Asia/Manila",
	"cn": "Asia/Shanghai",
	"hk": "Asia/Hong_Kong",
	"mo": "Asia/Macau",
	"tw": "Asia/Taipei",
	"jp": "Asia/Tokyo",
	"kr": "Asia/Seoul",
	"kp": "Asia/Pyongyang",
	"mn": "Asia/Ulaanbaatar",
	"eg": "Africa/Cairo",
	"ly": "Africa/Tripoli",
	"tn": "Africa/Tunis",
	"dz": "Africa/Algiers",
	"ma": "Africa/Casablanca",
	"ng": "Africa/Lagos",
	"gh": "Africa/Accra",
	"ke": "Africa/Nairobi",
	"tz": "Africa/Dar_es_Salaam",
	"ug": "Africa/Kampala",
	"et": "Africa/Addis_Ababa",
	"za": "Africa/Johannesburg",
	"zw": "Africa/Harare",
	"ao": "Africa/Luanda",
	"mz": "Africa/Maputo",
	"bw": "Africa/Gaborone",
	"au": "Australia/Sydney",
	"nz": "Pacific/Auckland",
}
