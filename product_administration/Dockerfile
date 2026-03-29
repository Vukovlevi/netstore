#1: build React
FROM node:24 AS frontend-build

WORKDIR /app
COPY . .
ARG VITE_STORE_ADMIN_URL
ENV VITE_STORE_ADMIN_URL=$VITE_STORE_ADMIN_URL
RUN npm install
RUN npm run build

#2: PHP + Apache
FROM php:8.2-apache

RUN docker-php-ext-install mysqli pdo pdo_mysql

# copy built React app
COPY --from=frontend-build /app/dist/ /var/www/html/

# copy PHP API
COPY api/ /var/www/html/api/

COPY ./.htaccess /var/www/html/

# Apache config (optional: enable rewrite)
RUN a2enmod rewrite

EXPOSE 80