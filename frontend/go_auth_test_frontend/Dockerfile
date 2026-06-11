# syntax=docker/dockerfile:1

############################
# 1. Build stage
############################
FROM node:22-alpine AS build

WORKDIR /app

# Install dependencies (cached unless lockfile/package.json change)
COPY package.json package-lock.json ./
RUN npm ci

# Backend URL is baked into the static bundle at build time.
# Override per environment: docker build --build-arg VITE_BACKEND_URL=https://api.example.com
ARG VITE_BACKEND_URL
ENV VITE_BACKEND_URL=${VITE_BACKEND_URL}

# Copy source and build
COPY . .
RUN npm run build

############################
# 2. Runtime stage (nginx)
############################
FROM nginx:1.27-alpine AS runtime

# SPA-friendly nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Static build output
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
