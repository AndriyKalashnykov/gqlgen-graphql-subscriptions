# https://hub.docker.com/r/nginxinc/nginx-unprivileged/tags?page=1&ordering=last_updated
FROM nginxinc/nginx-unprivileged:1.29@sha256:bd68c793125e39123d1fdf3c0e9b7d4218667b9db049ed46b1b71aa4cb4493ef
USER root
ENV DEBIAN_FRONTEND=noninteractive
ENV APT="apt-get --no-install-recommends --yes"
RUN $APT update && $APT upgrade && $APT install bash curl wget
COPY build/ /usr/share/nginx/html
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf
RUN chown -R nginx:nginx /usr/share/nginx
USER nginx

#FROM node:14.17.6
#WORKDIR /app
#COPY package.json ./
#COPY package-lock.json ./
#COPY ./public ./public
#COPY ./src ./src
#COPY ./package.json ./package.json
#COPY ./package-lock.json ./package-lock.json
#RUN npx browserslist@latest --update-db
#RUN npm ci --production
#CMD ["npm", "run", "start"]
