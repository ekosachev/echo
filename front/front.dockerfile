FROM node:20-alpine AS builder

WORKDIR /app

COPY package.json package_lock.json* yarn.lock* ./

RUN npm install

COPY . .
RUN npm run build

FROM node:20-alpine AS runner

COPY package.json package-lock.json* yarn.lock* ./

RUN npm install --omit=dev

COPY --from=builder /app/.next ./.next

COPY --from=builder /app/public ./public


EXPOSE 3000
CMD ["npm", "start"]
