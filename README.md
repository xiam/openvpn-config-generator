# ovpn-cfgen

`ovpn-cfgen` is a command line tool that generates certificates and (portable)
configuration files for [OpenVPN](https://openvpn.net/download-open-vpn/) that
_just work_.

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

The following recipe assumes you followed the steps above and that you have a
CA certificate and key (`ca.crt` and `ca.key`), as well as the server
(`server.crt` and `server.key`) and client (`my-laptop.crt` and
`my-laptop.key`) key pairs.

### Generate a server configuration file

Create additional keys `dh.pem` and `key.tlsauth`:

```
openssl dhparam -out dh.pem 2048
# takes a while...

openvpn --genkey --secret key.tlsauth
```

Use the `server-config` command to generate a configuration file for OpenVPN server:

```
ovpn-cfgen server-config
# 2019/05/30 23:11:21 Your new server configuration file was written to: "server.conf"
```

### Generate a client configuration file

```
ovpn-cfgen client-config \
  --remote 127.0.0.1 \
  --cert my-laptop.crt \
  --key my-laptop.key \
  --output my-laptop.ovpn

# 2019/05/30 23:15:10 Your new client configuration file was written to: "my-laptop.ovpn"
```

## Using your new configuration files

Spin up your OpenVPN server:

```
mkdir -p ccd
sudo openvpn --config server.conf
# ...
# Thu May 30 23:44:21 2019 us=100431 my-laptop/127.0.0.1:58334 Incoming Data Channel: Cipher 'AES-256-GCM' initialized with 256 bit key
```

Spin up your OpenVPN client:

```
sudo openvpn --config my-laptop.ovpn
# ...
# Thu May 30 23:44:21 2019 Initialization Sequence Completed
```
