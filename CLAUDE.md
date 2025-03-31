# Hostel Booking System Development Guidelines

## Build Commands
- Backend build: `go build -o main ./cmd/api`
- Backend run: `docker-compose up --build -d backend`
- Frontend dev: `cd frontend/hostel-frontend && npm start`
- Frontend build: `cd frontend/hostel-frontend && npm run build`

## Test Commands
- Backend tests: `go test ./...`
- Single test: `go test ./path/to/package -run TestName`
- Frontend tests: `cd frontend/hostel-frontend && npm test`

## Lint/Format Commands
- Check Go imports: `./backend/check_imports.sh`
- Frontend lint: Uses ESLint with react-app config

## Code Style Guidelines
- **Go Imports**: Standard lib first, project imports second, third-party last
- **Go Naming**: CamelCase, interfaces with "Interface" suffix
- **Go Functions**: New* constructors, descriptive names, context as first param
- **Go Errors**: If err != nil pattern, utils.ErrorResponse for HTTP errors
- **React Components**: Functional components with hooks, consistent prop types
- **Frontend Structure**: Component/page separation, context for global state