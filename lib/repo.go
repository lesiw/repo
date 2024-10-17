package repo

import (
	"fmt"
	"net/http"
	neturl "net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"lesiw.io/defers"
)

func Clone(prefix, url string, force bool) (string, error) {
	parts, err := urlToPath(url)
	if err != nil {
		return "", err
	}

	dirpath := filepath.Join(parts[:len(parts)-1]...)
	dirpath = filepath.Join(prefix, dirpath)
	fullpath := filepath.Join(parts...)
	fullpath = filepath.Join(prefix, fullpath)
	if info, err := os.Stat(fullpath); err == nil {
		if !info.IsDir() {
			return "", fmt.Errorf("'%s' exists and is not a directory", err)
		}
		if force {
			if err := os.RemoveAll(fullpath); err != nil {
				return "", fmt.Errorf("could not remove directory '%s': %s",
					fullpath, err)
			}
		} else {
			return fullpath, nil
		}
	}

	newdirs, err := MkdirAll(dirpath, 0755)
	if err != nil {
		return "", fmt.Errorf("could not create repo directory: %s", err)
	}

	defers.Add(func() { rmDirs(newdirs) }) // Stronger guarantee for cli use.
	defer rmDirs(newdirs)                  // Best effort for library use.

	if err := cloneRepo(url, fullpath); err != nil {
		return "", fmt.Errorf("could not clone repository: %s", err)
	}
	return fullpath, nil
}

func splitUrl(url string) (prefix, path, suffix string) {
	var ok bool
	if prefix, path, ok = strings.Cut(url, "@"); ok {
		prefix = prefix + "@"
		path = strings.Replace(path, ":", "/", 1)
	} else {
		path = prefix
		prefix = ""
	}
	parsed, err := neturl.Parse(path)
	if err != nil {
		return
	}
	if prefix == "" && parsed.Scheme != "" {
		prefix = parsed.Scheme + "://"
	}
	if path, ok = strings.CutSuffix(parsed.Host+parsed.Path, ".git"); ok {
		suffix += ".git"
	}
	if parsed.RawQuery != "" {
		suffix += "?" + parsed.RawQuery
	}
	return
}

func mergeUrl(prefix, path, suffix string) string {
	if strings.Contains(prefix, "@") && !strings.Contains(prefix, "://") {
		path = strings.Replace(path, "/", ":", 1)
	} else if prefix == "" {
		prefix = "https://"
	}
	return prefix + path + suffix
}

func urlToPath(url string) (path []string, err error) {
	_, rawpath, _ := splitUrl(url)
	for _, p := range strings.Split(rawpath, "/") {
		if p != "" {
			path = append(path, p)
		}
	}
	if len(path) < 1 {
		err = fmt.Errorf("failed to derive path from url: %s", url)
	}
	return
}

func cloneRepo(loc, path string) error {
	loc = followRedirects(loc)

	cmd := exec.Command("git", "clone", loc, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stderr // Stdout must only contain repodir.
	cmd.Stderr = os.Stderr

	defers.Add(func() { _ = cmd.Wait() })

	if os.Getenv("REPOVERBOSE") == "1" {
		fmt.Fprintf(os.Stderr, "%q\n", cmd.Args)
	}
	return cmd.Run()
}

func followRedirects(url string) (ret string) {
	prefix, path, suffix := splitUrl(url)
	defer func() { ret = mergeUrl(prefix, path, suffix) }()
	resp, err := new(http.Client).Get("https://" + path)
	if err != nil {
		return
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Don't update the path if we ended up on the same domain.
		// It's likely to be a login page for a private git server.
		// Most git hosts will automatically handle moved repositories when
		// they are cloned, so ignoring valid redirects here should be safe.
		host, _, _ := strings.Cut(path, "/")
		if host != resp.Request.URL.Host {
			path = resp.Request.URL.Host + resp.Request.URL.Path
		}
	}
	return
}

func rmDirs(dirs []string) {
	for i := len(dirs) - 1; i >= 0; i-- {
		_ = os.Remove(dirs[i])
	}
}
