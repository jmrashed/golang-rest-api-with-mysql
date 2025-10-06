# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.2.0] - 2025-10-06

### Added
- Database seeder for initial data population
- Automatic seeding of roles, permissions, and admin user
- Swagger UI integration for interactive API documentation
- Environment variable loading with godotenv
- Enhanced database connection handling for empty passwords
- Static file serving for Swagger documentation
- Default admin user (username: admin, password: admin123)

### Fixed
- Database connection issues with empty password configuration
- Duplicate User struct declaration in model package
- Unused import statements causing compilation errors
- Swagger UI file serving and routing conflicts
- Environment variable handling for empty values

### Changed
- Improved database configuration with automatic database creation
- Enhanced error handling in database connection
- Updated Swagger documentation with correct server URLs
- Consolidated model definitions in single file

### Technical
- Added seeder package for database initialization
- Integrated Swagger UI with http-swagger package
- Enhanced route configuration with documentation endpoints
- Improved environment variable parsing logic

## [1.1.0] - 2025-10-06

### Added
- JWT-based authentication with access and refresh tokens
- Role-based access control (RBAC) with admin, moderator, and user roles
- Permission-based access control for fine-grained authorization
- Secure password hashing with bcrypt
- User registration and login endpoints
- User profile management (get, update, change password)
- Multi-device logout functionality
- Todo CRUD operations with user ownership
- Pagination, filtering, and sorting for data retrieval
- Rate limiting middleware (60 requests/minute)
- In-memory caching for improved performance
- Request/response logging middleware
- CORS support for cross-origin requests
- Input validation for all API requests
- Comprehensive error handling with structured responses
- Database schema with users, roles, permissions, and todos
- Health check endpoint with service status
- Docker support with multi-stage builds
- Docker Compose for development environment
- CI/CD pipeline with GitHub Actions
- OpenAPI/Swagger documentation
- Comprehensive unit and integration tests
- End-to-end tests for critical user flows
- Makefile for development tasks
- Environment-based configuration

### Security
- Secure token storage with hashed refresh tokens
- Password validation requirements
- Token expiration (15 min access, 7 days refresh)
- Protection against common vulnerabilities
- Rate limiting to prevent abuse

### Performance
- Database connection pooling
- Optimized database queries with indexing
- Caching middleware for GET requests
- Efficient pagination implementation

### Documentation
- Complete API documentation
- OpenAPI 3.0 specification
- Postman collection examples
- cURL examples for all endpoints
- Development and deployment guides

## [1.0.0] - 2024-01-01

### Added
- JWT-based authentication with access and refresh tokens
- Role-based access control (RBAC) with admin, moderator, and user roles
- Permission-based access control for fine-grained authorization
- Secure password hashing with bcrypt
- User registration and login endpoints
- User profile management (get, update, change password)
- Multi-device logout functionality
- Todo CRUD operations with user ownership
- Pagination, filtering, and sorting for data retrieval
- Rate limiting middleware (60 requests/minute)
- In-memory caching for improved performance
- Request/response logging middleware
- CORS support for cross-origin requests
- Input validation for all API requests
- Comprehensive error handling with structured responses
- Database schema with users, roles, permissions, and todos
- Health check endpoint with service status
- Docker support with multi-stage builds
- Docker Compose for development environment
- CI/CD pipeline with GitHub Actions
- OpenAPI/Swagger documentation
- Comprehensive unit and integration tests
- End-to-end tests for critical user flows
- Makefile for development tasks
- Environment-based configuration

### Security
- Secure token storage with hashed refresh tokens
- Password validation requirements
- Token expiration (15 min access, 7 days refresh)
- Protection against common vulnerabilities
- Rate limiting to prevent abuse

### Performance
- Database connection pooling
- Optimized database queries with indexing
- Caching middleware for GET requests
- Efficient pagination implementation

### Documentation
- Complete API documentation
- OpenAPI 3.0 specification
- Postman collection examples
- cURL examples for all endpoints
- Development and deployment guides