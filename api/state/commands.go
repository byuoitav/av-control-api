package state

/*
var (
	errCommandNotFound    = errors.New("command not found")
	errCommandEnvNotFound = errors.New("command environment not found")
	errNotEnoughParams    = errors.New("not enough parameters provided")
)

func getCommand(d api.Device, cmd string, env string) (string, *int, error) {
	c, ok := d.Type.Commands[cmd]
	if !ok {
		return "", nil, errCommandNotFound
	}

	cmdURL, ok := c.URLs[env]
	if !ok {
		return "", nil, errCommandEnvNotFound
	}

	u, err := url.Parse(cmdURL)
	if err != nil {
		return "", nil, fmt.Errorf("unable to parse url: %w", err)
	}

	for reg, proxy := range d.Proxy {
		if reg.MatchString(cmd) {
			// use this proxy
			var host strings.Builder

			oldhost := strings.Split(u.Host, ":")
			newhost := strings.Split(proxy, ":")

			switch len(newhost) {
			case 1: // no port on the proxy url
				host.WriteString(newhost[0])

				// add on the old port if there was one
				if len(oldhost) > 1 {
					host.WriteString(":")
					host.WriteString(oldhost[1])
				}
			case 2: // port present on proxy url
				host.WriteString(newhost[0])
				host.WriteString(":")
				host.WriteString(newhost[1])
			default:
				return "", nil, fmt.Errorf("invalid proxy value %q", proxy)
			}

			u.Host = host.String()
			break
		}
	}

	// un-urlencode the `{{}}` around templated values
	s, err := url.PathUnescape(u.String())
	if err != nil {
		return "", nil, fmt.Errorf("unable to unescape path: %w", err)
	}

	return s, c.Order, nil
}

func fillURL(url string, params map[string]string) (string, error) {
	for param, val := range params {
		param = "{{" + param + "}}"
		url = strings.Replace(url, param, val, -1)
	}

	if strings.Contains(url, "{{") {
		return url, errNotEnoughParams
	}

	return url, nil
}
*/
