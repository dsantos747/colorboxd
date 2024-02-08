FROM node:20-alpine AS build

WORKDIR /app

COPY package.json .
COPY tsconfig.json .
RUN npm install

COPY . .

RUN npm run build

# Production stage
FROM nginx:1.24 AS production

# Replace default conf with mine
RUN rm -rf /etc/nginx/conf.d/*
COPY nginx.conf /etc/nginx/conf.d/

# Copy only the build items
COPY --from=build /app/build /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]