# GitHub Proxy Utility

Work w/ [purple4pur/docker_compose/github_proxy](https://github.com/purple4pur/docker_compose/tree/master/github_proxy).

- Replace 302 Redirect location of `raw.githubusercontent.com`.
- Replace 302 Redirect location of `objects.githubusercontent.com`.
- Replace `github.com` links.
- Replace `github.githubassets.com` links.
- Replace `avatars.githubusercontent.com` links.

## Limitation / Known issue

- Slower speed compared to direct access.
- Not support login.

## Example usage

```sh
GH_PROXY_HOSTNAME=hostname.com go run main.go 8080->http://localhost:8080 8081->http://localhost:8081
```
