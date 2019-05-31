# ovpn-cfgen

`ovpn-cfgen` is a configuration tool for openvpn that generates certificates
and configuration files for OpenVPN that _just work_.

## Installing `ovpn-cfgen`

```
go install github.com/xiam/openvpn-config-generator/cmd/ovpn-cfgen
```

## Using `ovpn-cfgen` to generate certificates

### Create a self-signed root CA certificate

```
ovpn-cfgen build-ca
# 2019/05/29 21:53:28 Your new CA certificate was successfully generated.
# 2019/05/29 21:53:28 certificate: "/home/rev/ca.crt"
# 2019/05/29 21:53:28 private key: "/home/rev/ca.key"

openssl x509 -in ca.crt -noout -text
```

### Create a server certificate

```
ovpn-cfgen build-key-server
# 2019/05/29 21:54:22 Your new server certificate was successfully generated.
# 2019/05/29 21:54:22 certificate: "/home/rev/server.crt"
# 2019/05/29 21:54:22 private key: "/home/rev/server.key"

openssl x509 -in server.crt -noout -text
```

### Create a client certificate

```
ovpn-cfgen build-key --name my-laptop                                                                  ~
# 2019/05/29 21:54:59 Your new client certificate was successfully generated.
# 2019/05/29 21:54:59 certificate: "/home/rev/my-laptop.crt"
# 2019/05/29 21:54:59 private key: "/home/rev/my-laptop.key"

openssl x509 -in my-laptop.crt -noout -text
```

## Using `ovpn-cfgen` to generate config files for OpenVPN

The following recipe assumes you followed the steps above to build your CA
certificate and key (`ca.crt` and `ca.key`), as well as the server
(`server.crt` and `server.key`) and client (`my-laptop.crt` and
`my-laptop.key`) key pairs.

### Generate a server configuration file

Create `dh.pem` and `key.tlsauth`:

```
openssl dhparam -out dh.pem 2048
```

```
openvpn --genkey --secret key.tlsauth
```

Use the `server-config` command to generate a configuration file for OpenVPN server:

```
ovpn-cfgen server-config
# 2019/05/30 23:11:21 Your new configuration file for OpenVPN server was written to: "server.conf"
```

### Generate a client configuration file

```
ovpn-cfgen client-config \
  --remote 192.168.1.9 \
  --cert my-laptop.crt \
  --key my-laptop.key \
  --output my-laptop.ovpn
2019/05/30 23:15:10 Your new configuration file for OpenVPN client was written to: "my-laptop.ovpn"
```
