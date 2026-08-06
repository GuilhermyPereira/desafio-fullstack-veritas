# Estágio de construção (Builder)
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Estágio de produção (Servidor Nginx)
FROM nginx:alpine
# Copia os arquivos estáticos gerados pelo Vite para a pasta pública do Nginx
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]