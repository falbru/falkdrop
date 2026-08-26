FROM node:24-alpine AS builder

RUN corepack enable
RUN corepack prepare pnpm@10.24 --activate
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
ENV CI=true

ARG VITE_API_URL
ARG VITE_KEYCLOAK_URL
ARG VITE_KEYCLOAK_REALM
ARG VITE_KEYCLOAK_CLIENT_ID
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_KEYCLOAK_URL=$VITE_KEYCLOAK_URL
ENV VITE_KEYCLOAK_REALM=$VITE_KEYCLOAK_REALM
ENV VITE_KEYCLOAK_CLIENT_ID=$VITE_KEYCLOAK_CLIENT_ID

WORKDIR /app
COPY package.json .
COPY pnpm-lock.yaml .
RUN echo 'dangerouslyAllowAllBuilds: true' > pnpm-workspace.yaml
RUN pnpm install --frozen-lockfile
COPY . .
RUN pnpm build

FROM nginx:alpine AS runner
WORKDIR /app
COPY --from=builder /app/dist /usr/share/nginx/html
RUN echo 'server {\
    listen 80;\
\
    root /usr/share/nginx/html;\
    index index.html;\
\
    location / {\
        try_files $uri $uri/ /index.html;\
    }\
}' > /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
