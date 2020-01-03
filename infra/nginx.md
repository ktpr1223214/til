---
title: Nginx
---

## basic
* [remote_addrとかx-forwarded-forとかx-real-ipとか](https://christina04.hatenablog.com/entry/2016/10/25/190000)

### docker
* nginx の公式 Docker image では、/var/log/nginx/ 配下の log ファイルは /dev/stdout /dev/stderr にシンボリックリンクが貼られている
  * cf. https://github.com/nginxinc/docker-nginx/blob/a973c221f6cedede4dab3ab36d18240c4d3e3d74/stable/alpine/Dockerfile#L102

## Reference
* [nginx-admins-handbook](https://github.com/trimstray/nginx-admins-handbook)


