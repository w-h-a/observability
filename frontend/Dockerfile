FROM node as build
WORKDIR /trace-blame-frontend
COPY package.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx
COPY conf/default.conf /etc/nginx/conf.d/default.conf
RUN rm -rf /usr/share/nginx/html/*
COPY --from=build /trace-blame-frontend/build /usr/share/nginx/html
ENTRYPOINT ["nginx", "-g", "daemon off;"]