---
icon: lucide/shield
---

# Zscaler scripts

The `zs` command manages the Zscaler corporate proxy when one is installed, with three verbs:
`enable`, `disable`, and `certs`. All no-op cleanly when the relevant Zscaler launchd plists are
absent — safe on personal machines. It is now a [`dot`](../tooling/dot.md) applet (same behaviour;
the shell snippets below are illustrative).

## `zs enable` { #zs-enable }

```sh
zs enable
```

Loads the Zscaler service and tunnel launch daemons if they're not already running:

```sh title=".local/share/scripts/zs enable (core)"
ZSCALER_SERVICE="/Library/LaunchDaemons/com.zscaler.service.plist"
ZSCALER_TUNNEL="/Library/LaunchDaemons/com.zscaler.tunnel.plist"

if [ -f "${ZSCALER_SERVICE}" ] && ! sudo launchctl list | grep 'com.zscaler.service' > /dev/null; then
    sudo launchctl load "${ZSCALER_SERVICE}"
    started=1
fi

if [ -f "${ZSCALER_TUNNEL}" ] && ! sudo launchctl list | grep 'com.zscaler.tunnel' > /dev/null; then
    sudo launchctl load "${ZSCALER_TUNNEL}"
    started=1
fi

# Wait 10s for the tunnel to come up
[ -n "${started:-}" ] && sleep 10
```

## `zs disable` { #zs-disable }

```sh
zs disable
```

Unloads both daemons (tunnel first, service second).

## `zs certs` { #zs-certs }

```sh
zs certs -- <command> [args...]
```

Runs a command with the Zscaler root CA injected into the trust path. Useful for tools like
Node, Python, and curl that don't read the system keychain by default.

Flow:

1. If `~/.local/share/certificates/zscaler.pem` exists, use it.
2. Otherwise, extract the Zscaler root CA from the system keychain via
   `security find-certificate -c "Zscaler Root CA" -p` and cache it in
   `~/.local/share/certificates/`.
3. Export `ZSCALER_CERTIFICATE` and `NODE_EXTRA_CA_CERTS` pointing at the cert.
4. `exec` the command.

```sh
zs certs -- npm install
zs certs -- curl https://artifactory.example/api/something
```

If no certificate is found anywhere, the command still runs — with a warning, but without the
extra CA. You'll get TLS errors if Zscaler is intercepting; that's the signal to import the
cert into your keychain.
