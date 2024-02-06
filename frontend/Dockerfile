FROM node:20-alpine AS base

WORKDIR /app

# Copy package.json and install deps
COPY package.json .
RUN npm install --production --silent

# Copy everything except what's in dockerignore
COPY . .

EXPOSE 8080

CMD [ "npm", "start" ]