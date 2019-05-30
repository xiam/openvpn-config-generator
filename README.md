# ovpn-cfgen

`ovpn-cfgen` is a configuration tool for openvpn that generates configuration
files and certificates the easy way.

## Get `ovpn-cfgen`

```
go install github.com/xiam/openvpn-config-generator/cmd/ovpn-cfgen
```

## Create a self-signed root CA certificate

```
ovpn-cfgen build-ca
# 2019/05/29 21:53:28 Your new CA certificate was successfully generated.
# 2019/05/29 21:53:28 certificate: "/home/rev/ca.crt"
# 2019/05/29 21:53:28 private key: "/home/rev/ca.key"

openssl x509 -in ca.crt -noout -text
```

## Create a server certificate

ovpn-cfgen build-key-server
# 2019/05/29 21:54:22 Your new server certificate was successfully generated.
# 2019/05/29 21:54:22 certificate: "/home/rev/server.crt"
# 2019/05/29 21:54:22 private key: "/home/rev/server.key"

openssl x509 -in server.crt -noout -text
```

## Create a client certificate

```
ovpn-cfgen build-key --name my-laptop                                                                  ~
# 2019/05/29 21:54:59 Your new client certificate was successfully generated.
# 2019/05/29 21:54:59 certificate: "/home/rev/my-laptop.crt"
# 2019/05/29 21:54:59 private key: "/home/rev/my-laptop.key"

openssl x509 -in my-laptop.crt -noout -text
```
