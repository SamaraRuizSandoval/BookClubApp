# ---- build stage ----
FROM node:20-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build   # puts optimized static files in /app/build

# ---- run stage ----
FROM nginx:1.27-alpine
COPY --from=build /app/build /usr/share/nginx/html

EXPOSE 80
