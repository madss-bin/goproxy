package pipe

import (
	"bytes"
	"net/url"
	"strings"

	"animex/goproxy/crypto"
)

func BuildProxyUrl(resolved *url.URL, originParam string) []byte {
	encrypted := crypto.EncryptURL(resolved.String())
	var buf bytes.Buffer
	buf.WriteString("/?u=")
	buf.WriteString(encrypted)
	if originParam != "" {
		buf.WriteString("&origin=")
		buf.WriteString(url.QueryEscape(originParam))
	}
	return buf.Bytes()
}

func getUrl(line []byte, base *url.URL) *url.URL {
	u, err := url.Parse(string(line))
	if err == nil && u.IsAbs() {
		return u
	}
	if err != nil {
		return base.ResolveReference(&url.URL{Path: string(line)})
	}
	return base.ResolveReference(u)
}

func ProcessPipeBody(body []byte, baseUrl *url.URL, originParam string) []byte {
	lines := bytes.Split(body, []byte("\n"))
	var result bytes.Buffer
	result.Grow(len(body) + len(lines)*60)

	for i, line := range lines {
		lineTrimmed := bytes.TrimSpace(line)
		if len(lineTrimmed) == 0 {
			if i != len(lines)-1 {
				result.WriteByte('\n')
			}
			continue
		}

		if lineTrimmed[0] == '#' {
			if len(lineTrimmed) > 11 && bytes.HasPrefix(lineTrimmed, []byte("#EXT-X-KEY")) {
				if uriStart := bytes.Index(lineTrimmed, []byte("URI=\"")); uriStart != -1 {
					keyUriStart := uriStart + 5
					if quotePos := bytes.IndexByte(lineTrimmed[keyUriStart:], '"'); quotePos != -1 {
						keyUriEnd := keyUriStart + quotePos
						keyUri := lineTrimmed[keyUriStart:keyUriEnd]
						resolved := getUrl(keyUri, baseUrl)
						proxyUrl := BuildProxyUrl(resolved, originParam)

						result.Write(lineTrimmed[:keyUriStart])
						result.Write(proxyUrl)
						result.Write(lineTrimmed[keyUriEnd:])
					} else {
						result.Write(lineTrimmed)
					}
				} else {
					result.Write(lineTrimmed)
				}
			} else if len(lineTrimmed) > 16 && bytes.HasPrefix(lineTrimmed, []byte("#EXT-X-MAP:URI=\"")) {
				innerUrl := lineTrimmed[16 : len(lineTrimmed)-1]
				resolved := getUrl(innerUrl, baseUrl)
				proxyUrl := BuildProxyUrl(resolved, originParam)

				result.WriteString(`#EXT-X-MAP:URI="`)
				result.Write(proxyUrl)
				result.WriteByte('"')
			} else if len(lineTrimmed) > 20 && (bytes.Contains(lineTrimmed, []byte("URI=")) || bytes.Contains(lineTrimmed, []byte("URL="))) {
				colonPos := bytes.IndexByte(lineTrimmed, ':')
				if colonPos != -1 {
					prefix := lineTrimmed[:colonPos+1]
					attrs := lineTrimmed[colonPos+1:]
					result.Write(prefix)

					inQuotes := false
					var currentAttr bytes.Buffer
					var parsedAttrs [][]byte

					for _, c := range attrs {
						if c == '"' {
							inQuotes = !inQuotes
						}
						if c == ',' && !inQuotes {
							parsed := make([]byte, currentAttr.Len())
							copy(parsed, currentAttr.Bytes())
							parsedAttrs = append(parsedAttrs, parsed)
							currentAttr.Reset()
						} else {
							currentAttr.WriteByte(c)
						}
					}
					parsed := make([]byte, currentAttr.Len())
					copy(parsed, currentAttr.Bytes())
					parsedAttrs = append(parsedAttrs, parsed)

					for j, attr := range parsedAttrs {
						if j > 0 {
							result.WriteByte(',')
						}

						if eqPos := bytes.IndexByte(attr, '='); eqPos != -1 {
							key := bytes.TrimSpace(attr[:eqPos])
							val := bytes.TrimSpace(attr[eqPos+1:])
							val = bytes.Trim(val, `"`)

							if bytes.Equal(key, []byte("URI")) || strings.EqualFold(string(key), "URL") {
								resolved := getUrl(val, baseUrl)
								proxyUrl := BuildProxyUrl(resolved, originParam)
								result.Write(key)
								result.WriteString(`="`)
								result.Write(proxyUrl)
								result.WriteByte('"')
							} else {
								result.Write(attr)
							}
						} else {
							result.Write(attr)
						}
					}
				} else {
					result.Write(lineTrimmed)
				}
			} else {
				result.Write(lineTrimmed)
			}
		} else {
			resolved := getUrl(lineTrimmed, baseUrl)
			proxyUrl := BuildProxyUrl(resolved, originParam)
			result.Write(proxyUrl)
		}

		if i != len(lines)-1 {
			result.WriteByte('\n')
		}
	}
	return result.Bytes()
}
