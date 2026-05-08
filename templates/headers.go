package templates

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
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
	patterns      []*regexp.Regexp
	origin        string
	referer       string
	customHeaders map[string]string
}

var domainGroups []DomainGroup

func init() {
	domainGroups = []DomainGroup{
		{
			patterns: compileRegexes(`(?i)\.padorupado\.ru$`, `(?i)\.kwikie\.ru$`, `(?i)\.owocdn\.top$`),
			origin:   "https://kwik.cx",
			referer:  "https://kwik.cx/",
			customHeaders: map[string]string{
				"cache-control": "no-cache",
				"pragma":        "no-cache",
			},
		},
		{
			patterns: compileRegexes(`(?i)[a-z]+\d+.(pro|live|xyz|site|online|wiki)$`),
			origin:   "https://rapid-cloud.co/",
			referer:  "https://rapid-cloud.co/",
			customHeaders: map[string]string{
				"cache-control": "no-cache",
				"pragma":        "no-cache",
			},
		},
		{
			patterns: compileRegexes(`(?i)\.streamtape\.to$`),
			origin:   "https://streamtape.to",
			referer:  "https://streamtape.to/",
		},
		{
			patterns: compileRegexes(`(?i)vidcache\.net$`),
			origin:   "https://www.animegg.org",
			referer:  "https://www.animegg.org/",
		},
		{
			patterns: compileRegexes(`(?i)krussdomi\.com$`, `(?i)revolutionizingtheweb\.xyz$`, `(?i)nextgentechnologytrends\.xyz$`, `(?i)smartinvestmentstrategies\.xyz$`, `(?i)creativedesignstudioxyz\.xyz$`, `(?i)breakingdigitalboundaries\.xyz$`, `(?i)ultimatetechinnovation\.xyz$`),
			origin:   "https://krussdomi.com",
			referer:  "https://krussdomi.com/",
		},
		{
			patterns: compileRegexes(`(?i)\.akamaized\.net$`),
			origin:   "https://players.akamai.com",
			referer:  "https://players.akamai.com/",
		},
		{
			patterns: compileRegexes(`(?i)(?:^|\.)shadowlandschronicles\.`, `(?i)digitalshinecollective\.xyz$`, `(?i)thrivequesthub\.xyz$`, `(?i)novaedgelabs\.xyz$`),
			origin:   "https://cloudnestra.com",
			referer:  "https://cloudnestra.com/",
		},
		{
			patterns: compileRegexes(`(?i)(?:^|\.)viddsn\.`, `(?i)\.anilike\.cyou$`),
			origin:   "https://vidwish.live/",
			referer:  "https://vidwish.live/",
		},
		{
			patterns: compileRegexes(`(?i)(?:^|\.)dotstream\.`, `(?i)(?:^|\.)playcloud1\.`),
			origin:   "https://megaplay.buzz/",
			referer:  "https://megaplay.buzz/",
		},
		{
			patterns: compileRegexes(`(?i)\.cloudfront\.net$`),
			origin:   "https://d2zihajmogu5jn.cloudfront.net",
			referer:  "https://d2zihajmogu5jn.cloudfront.net/",
		},
		{
			patterns: compileRegexes(`(?i)\.ttvnw\.net$`),
			origin:   "https://www.twitch.tv",
			referer:  "https://www.twitch.tv/",
		},
		{
			patterns: compileRegexes(`(?i)\.xx\.fbcdn\.net$`),
			origin:   "https://www.facebook.com",
			referer:  "https://www.facebook.com/",
		},
		{
			patterns: compileRegexes(`(?i)\.anih1\.top$`, `(?i)\.xyk3\.top$`),
			origin:   "https://ee.anih1.top",
			referer:  "https://ee.anih1.top/",
		},
		{
			patterns: compileRegexes(`(?i)\.premilkyway\.com$`),
			origin:   "https://uqloads.xyz",
			referer:  "https://uqloads.xyz/",
		},
		{
			patterns: compileRegexes(`(?i)\.streamcdn\.com$`, `(?i)\.mediacache\.cc$`),
			origin:   "https://anime.uniquestream.net",
			referer:  "https://anime.uniquestream.net/",
		},
		{
			patterns: compileRegexes(`(?i)\.raffaellocdn\.net$`, `(?i)\.feetcdn\.com$`, `(?i)clearskydrift45\.site$`),
			origin:   "https://kerolaunochan.online",
			referer:  "https://kerolaunochan.online/",
		},
		{
			patterns: compileRegexes(`(?i)dewbreeze84\.online$`, `(?i)cloudydrift38\.site$`, `(?i)sunshinerays93\.live$`, `(?i)clearbluesky72\.wiki$`, `(?i)breezygale56\.online$`, `(?i)frostbite27\.pro$`, `(?i)frostywinds57\.live$`, `(?i)icyhailstorm64\.wiki$`, `(?i)icyhailstorm29\.online$`, `(?i)windflash93\.xyz$`, `(?i)stormdrift27\.site$`, `(?i)tempestcloud61\.wiki$`, `(?i)sunburst66\.pro$`, `(?i)douvid\.xyz$`),
			origin:   "https://megacloud.blog",
			referer:  "https://megacloud.blog/",
			customHeaders: map[string]string{
				"cache-control": "no-cache",
				"pragma":        "no-cache",
			},
		},
		{
			patterns: compileRegexes(`(?i)\.echovideo\.to$`),
			origin:   "https://aniwave.se",
			referer:  "https://aniwave.se/",
		},
		{
			patterns: compileRegexes(`(?i)\.vid-cdn\.xyz$`),
			origin:   "https://anizone.to/",
			referer:  "https://anizone.to/",
		},
		{
			patterns: compileRegexes(`(?i)\.1stkmgv1\.com$`),
			origin:   "https://animeyy.com",
			referer:  "https://animeyy.com/",
		},
		{
			patterns: compileRegexes(`(?i)lightningspark77\.pro$`, `(?i)thunderwave48\.xyz$`, `(?i)stormwatch95\.site$`, `(?i)windyrays29\.online$`, `(?i)thunderstrike77\.online$`, `(?i)lightningflash39\.live$`, `(?i)cloudburst82\.xyz$`, `(?i)drizzleshower19\.site$`, `(?i)rainstorm92\.xyz$`),
			origin:   "https://megacloud.club",
			referer:  "https://megacloud.club/",
		},
		{
			patterns: compileRegexes(`(?i)cloudburst99\.xyz$`, `(?i)frostywinds73\.pro$`, `(?i)stormwatch39\.live$`, `(?i)sunnybreeze16\.live$`, `(?i)mistydawn62\.pro$`, `(?i)lightningbolt21\.live$`, `(?i)gentlebreeze85\.xyz$`),
			origin:   "https://videostr.net",
			referer:  "https://videostr.net/",
		},
		{
			patterns: compileRegexes(`(?i)vmeas\.cloud$`),
			origin:   "https://vidmoly.to",
			referer:  "https://vidmoly.to/",
		},
		{
			patterns: compileRegexes(`(?i)nextwaveinitiative\.xyz$`),
			origin:   "https://edgedeliverynetwork.org",
			referer:  "https://edgedeliverynetwork.org/",
		},
		{
			patterns: compileRegexes(`(?i)lightningbolts\.ru$`, `(?i)lightningbolt\.site$`, `(?i)vyebzzqlojvrl\.top$`),
			origin:   "https://vidsrc.cc",
			referer:  "https://vidsrc.cc/",
		},
		{
			patterns: compileRegexes(`(?i)vidlvod\.store$`),
			origin:   "https://vidlink.pro",
			referer:  "https://vidlink.pro/",
		},
		{
			patterns: compileRegexes(`(?i)heatwave90\.pro$`, `(?i)humidmist27\.wiki$`, `(?i)frozenbreeze65\.live$`, `(?i)drizzlerain73\.online$`, `(?i)sunrays81\.xyz$`),
			origin:   "https://kerolaunochan.live",
			referer:  "https://kerolaunochan.live/",
		},
		{
			patterns: compileRegexes(`(?i)\.vkcdn5\.com$`),
			origin:   "https://vkspeed.com",
			referer:  "https://vkspeed.com/",
		},
		{
			patterns: compileRegexes(`(?i)embed\.su$`, `(?i)usbigcdn\.cc$`, `(?i)\.congacdn\.cc$`),
			origin:   "https://embed.su",
			referer:  "https://embed.su/",
		},
	}
}

func compileRegexes(patterns ...string) []*regexp.Regexp {
	regexes := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		regexes[i] = regexp.MustCompile(pattern)
	}
	return regexes
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
		hostname := targetUrl.Hostname()
		matched := false

		for _, group := range domainGroups {
			for _, re := range group.patterns {
				if re.MatchString(hostname) {
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
				origin := fmt.Sprintf("%s://%s", scheme, hostname)
				req.Header.Set("origin", origin)
				referer := fmt.Sprintf("%s/", origin)
				req.Header.Set("referer", referer)
			}
		}
	}
}
