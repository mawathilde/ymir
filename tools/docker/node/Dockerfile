FROM node:20-alpine

RUN apk add --no-cache curl && \
    corepack enable && \
    corepack prepare pnpm@latest --activate

WORKDIR /usr/src/app

CMD ["pnpm", "install"]
