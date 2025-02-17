# Frontend Dockerfile
FROM node:18-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy all other files
COPY . .

# Add this to ensure Vite works in Docker
ENV VITE_HOST=0.0.0.0

EXPOSE 5173

# For development
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
