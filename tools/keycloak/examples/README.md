# authentication-learning

## setup
``` bash
$ openssl req -x509 -newkey rsa:2048 -keyout myservice.key -out myservice.cert -days 365 -nodes -subj "/CN=samltest.ktpr1223214.com"

```