FROM oven/bun as builder

WORKDIR /app
COPY package*.json ./
COPY bun*.lock ./

RUN bun install --frozen-lockfile

COPY ./ .
RUN bun run build-only

FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
