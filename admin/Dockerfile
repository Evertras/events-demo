FROM node:12 as builder

WORKDIR /usr/src/app

COPY package.json .

RUN npm install

COPY . .

RUN npm run build

FROM nginx:1.16-alpine

COPY --from=builder /usr/src/app/dist/ /usr/share/nginx/html
COPY deploy/nginx/nginx.conf /etc/nginx/conf.d/default.conf

