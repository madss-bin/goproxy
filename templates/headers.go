package templates

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var defaultHeaders = map[string]string{
	"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:137.0) Gecko/20100101 Firefox/137.0",
	"accept":          "*/*",
	"accept-language": "en-US,en;q=0.5",
	"sec-fetch-dest":  "empty",
	"sec-fetch-mode":  "cors",
	"sec-fetch-site":  "cross-site",
}

type DomainGroup struct {
	suffixes      []string
	origin        string
	referer       string
	customHeaders map[string]string
}

var domainGroups []DomainGroup

func init() {
	domainGroups = []DomainGroup{
		{
			suffixes: []string{".padorupado.ru", ".kwikie.ru", ".owocdn.top"},
			origin:   "https://kwik.cx",
			referer:  "https://kwik.cx/",
			customHeaders: map[string]string{
				"cache-control": "no-cache",
				"pragma":        "no-cache",
			},
		},
		{
			suffixes: []string{".pro", ".live", ".xyz", ".site", ".online", ".wiki"},
			origin:   "https://rapid-cloud.co/",
			referer:  "https://rapid-cloud.co/",
			customHeaders: map[string]string{
				"cache-control": "no-cache",
				"pragma":        "no-cache",
			},
		},
		{
			suffixes: []string{".streamtape.to"},
			origin:   "https://streamtape.to",
			referer:  "https://streamtape.to/",
		},
		{
			suffixes: []string{"vidcache.net"},
			origin:   "https://www.animegg.org",
			referer:  "https://www.animegg.org/",
		},
		{
			suffixes: []string{"krussdomi.com", "revolutionizingtheweb.xyz", "nextgentechnologytrends.xyz", "smartinvestmentstrategies.xyz", "creativedesignstudioxyz.xyz", "breakingdigitalboundaries.xyz", "ultimatetechinnovation.xyz"},
			origin:   "https://krussdomi.com",
			referer:  "https://krussdomi.com/",
		},
		{
			suffixes: []string{".akamaized.net"},
			origin:   "https://players.akamai.com",
			referer:  "https://players.akamai.com/",
		},
		{
			suffixes: []string{"shadowlandschronicles.", "digitalshinecollective.xyz", "thrivequesthub.xyz", "novaedgelabs.xyz"},
			origin:   "https://cloudnestra.com",
			referer:  "https://cloudnestra.com/",
		},
		{
			suffixes: []string{"viddsn.", ".anilike.cyou"},
			origin:   "https://vidwish.live/",
			referer:  "https://vidwish.live/",
		},
		{
			suffixes: []string{"dotstream.", "playcloud1."},
			origin:   "https://megaplay.buzz/",
			referer:  "https://megaplay.buzz/",
		},
		{
			suffixes: []string{".cloudfront.net"},
			origin:   "https://d2zihajmogu5jn.cloudfront.net",
			referer:  "https://d2zihajmogu5jn.cloudfront.net/",
		},
		{
			suffixes: []string{".ttvnw.net"},
			origin:   "https://www.twitch.tv",
			referer:  "https://www.twitch.tv/",
		},
		{
			suffixes: []string{".xx.fbcdn.net"},
			origin:   "https://www.facebook.com",
			referer:  "https://www.facebook.com/",
		},
		{
			suffixes: []string{".anih1.top", ".xyk3.top"},
			origin:   "https://ee.anih1.top",
			referer:  "https://ee.anih1.top/",
		},
		{
			suffixes: []string{".premilkyway.com"},
			origin:   "https://uqloads.xyz",
			referer:  "https://uqloads.xyz/",
		},
		{
			suffixes: []string{".streamcdn.com", ".mediacache.cc"},
			origin:   "https://anime.uniquestream.net",
			referer:  "https://anime.uniquestream.net/",
		},
		{
			suffixes: []string{".raffaellocdn.net", ".feetcdn.com", "clearskydrift45.site"},
			origin:   "https://kerolaunochan.online",
			referer:  "https://kerolaunochan.online/",
		},
		{
			suffixes: []string{"dewbreeze84.online", "cloudydrift38.site", "sunshinerays93.live", "clearbluesky72.wiki", "breezygale56.online", "frostbite27.pro", "frostywinds57.live", "icyhailstorm64.wiki", "icyhailstorm29.online", "windflash93.xyz", "stormdrift27.site", "tempestcloud61.wiki", "sunburst66.pro", "douvid.xyz"},
			origin:   "https://megacloud.blog",
			referer:  "https://megacloud.blog/",
			customHeaders: map[string]string{
				"cache-control": "no-cache",
				"pragma":        "no-cache",
			},
		},
		{
			suffixes: []string{".echovideo.to"},
			origin:   "https://aniwave.se",
			referer:  "https://aniwave.se/",
		},
		{
			suffixes: []string{".vid-cdn.xyz"},
			origin:   "https://anizone.to/",
			referer:  "https://anizone.to/",
		},
		{
			suffixes: []string{".1stkmgv1.com"},
			origin:   "https://animeyy.com",
			referer:  "https://animeyy.com/",
		},
		{
			suffixes: []string{"lightningspark77.pro", "thunderwave48.xyz", "stormwatch95.site", "windyrays29.online", "thunderstrike77.online", "lightningflash39.live", "cloudburst82.xyz", "drizzleshower19.site", "rainstorm92.xyz"},
			origin:   "https://megacloud.club",
			referer:  "https://megacloud.club/",
		},
		{
			suffixes: []string{"cloudburst99.xyz", "frostywinds73.pro", "stormwatch39.live", "sunnybreeze16.live", "mistydawn62.pro", "lightningbolt21.live", "gentlebreeze85.xyz"},
			origin:   "https://videostr.net",
			referer:  "https://videostr.net/",
		},
		{
			suffixes: []string{"vmeas.cloud"},
			origin:   "https://vidmoly.to",
			referer:  "https://vidmoly.to/",
		},
		{
			suffixes: []string{"nextwaveinitiative.xyz"},
			origin:   "https://edgedeliverynetwork.org",
			referer:  "https://edgedeliverynetwork.org/",
		},
		{
			suffixes: []string{"lightningbolts.ru", "lightningbolt.site", "vyebzzqlojvrl.top"},
			origin:   "https://vidsrc.cc",
			referer:  "https://vidsrc.cc/",
		},
		{
			suffixes: []string{"vidlvod.store"},
			origin:   "https://vidlink.pro",
			referer:  "https://vidlink.pro/",
		},
		{
			suffixes: []string{"heatwave90.pro", "humidmist27.wiki", "frozenbreeze65.live", "drizzlerain73.online", "sunrays81.xyz"},
			origin:   "https://kerolaunochan.live",
			referer:  "https://kerolaunochan.live/",
		},
		{
			suffixes: []string{".vkcdn5.com"},
			origin:   "https://vkspeed.com",
			referer:  "https://vkspeed.com/",
		},
		{
			suffixes: []string{"embed.su", "usbigcdn.cc", ".congacdn.cc"},
			origin:   "https://embed.su",
			referer:  "https://embed.su/",
		},
	}
}

func GenerateHeadersForUrl(req *http.Request, targetUrl *url.URL, customOrigin string) {
	for k, v := range defaultHeaders {
		req.Header.Set(k, v)
	}

	if customOrigin != "" {
		req.Header.Set("origin", customOrigin)
		referer := customOrigin
		if !strings.HasSuffix(referer, "/") {
			referer = fmt.Sprintf("%s/", referer)
		}
		req.Header.Set("referer", referer)
	} else {
		hostname := strings.ToLower(targetUrl.Hostname())
		matched := false

		for _, group := range domainGroups {
			for _, suffix := range group.suffixes {
				if strings.HasSuffix(hostname, suffix) {
					req.Header.Set("origin", group.origin)
					req.Header.Set("referer", group.referer)
					if group.customHeaders != nil {
						for k, v := range group.customHeaders {
							req.Header.Set(k, v)
						}
					}
					matched = true
					break
				}
			}
			if matched {
				break
			}
		}

		if !matched {
			scheme := targetUrl.Scheme
			if hostname != "" {
				origin := fmt.Sprintf("%s://%s", scheme, targetUrl.Hostname())
				req.Header.Set("origin", origin)
				referer := fmt.Sprintf("%s/", origin)
				req.Header.Set("referer", referer)
			}
		}
	}
}
