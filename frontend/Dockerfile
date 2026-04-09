# https://hub.docker.com/r/nginxinc/nginx-unprivileged/tags?page=1&ordering=last_updated
FROM nginxinc/nginx-unprivileged:1.29@sha256:5c8946de2a79b6031975ff6c54594a26bdb2f1b5abab1c03f501575ebb02510f
USER root
ENV DEBIAN_FRONTEND=noninteractive
ENV APT="apt-get --no-install-recommends --yes"
RUN $APT update && $APT upgrade && $APT install bash curl wget
COPY build/ /usr/share/nginx/html
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf
RUN chown -R nginx:nginx /usr/share/nginx
USER nginx
