# AI WebApps Builder (Clone of Bolt.new and Lovable.dev with go backend.)
A web application with Vite frontend and Go backend using Anthropic's API.

## Project Structure
```
├── frontend/
│   ├── node_modules/
│   ├── src/
│   ├── .gitignore
│   ├── eslint.config.js
│   ├── index.html
│   ├── package.json
│   ├── postcss.config.js
│   ├── tailwind.config.js
│   ├── tsconfig.json
│   └── vite.config.ts
└── go-backend/
    ├── constants/
    ├── prompts/
    ├── utils/
    ├── .env
    ├── .env.example
    ├── anthropic.go
    ├── go.mod
    ├── go.sum
    └── main.go
```

## Setup Instructions

### Frontend Setup
1. Navigate to frontend directory:
```
cd frontend
```

2. Install dependencies:
```
npm install
```

3. Run development server:
```
npm run dev
```

### Backend Setup
1. Navigate to go-backend directory:
```
cd go-backend
```

2. Create .env file from .env.example:
```
cp .env.example .env
```

3. Add your Anthropic API key to .env:
```
ANTHROPIC_API_KEY=your_api_key_here
```

4. Run the backend server:
```
go run .
```

## Technologies Used
- Frontend:
  - Vite
  - TypeScript
  - Tailwind CSS
  - ESLint
- Backend:
  - Go
  - Anthropic API

## Development
- Frontend runs on: http://localhost:5173
- Backend runs on: http://localhost:3000

## Environment Variables
Backend:
- ANTHROPIC_API_KEY: Your Anthropic API key (required)

## Docker Support
The project includes Docker configuration for both frontend and backend services.

To run with Docker:
```
docker compose up --build
```
```

Remember to:
1. Never commit your .env file to git (it's already in .gitignore)
2. Keep your API keys secure
3. Update the .env.example file with any new environment variables you add
