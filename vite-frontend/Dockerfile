FROM node:18-alpine

WORKDIR /app

# Install dependencies using the lockfile for reproducible builds
COPY package.json package-lock.json ./

RUN npm ci

RUN npm i -g serve

# Copy application source (excluding node_modules via .dockerignore)
COPY . .

RUN npm run build

EXPOSE 3000

CMD [ "serve", "-s", "dist" ]