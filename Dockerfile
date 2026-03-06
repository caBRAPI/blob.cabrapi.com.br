FROM node:20-alpine

WORKDIR /app

# Install dependencies first to maximize Docker layer cache usage.
COPY package.json package-lock.json* ./
RUN npm install --no-audit --no-fund

# Copy source after deps; .dockerignore keeps context lean.
COPY . .

# Prepare writable runtime paths for SQLite, objects and logs.
RUN mkdir -p /app/data/blob-storage /app/data/api/logs

ENV NODE_ENV=production

EXPOSE 3000

# Runtime keeps TS entrypoint via tsx and ensures local DB schema is prepared.
CMD ["sh", "-c", "npm run db:prepare && npm start"]
