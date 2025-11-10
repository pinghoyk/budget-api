# Changelog

## [1.0.0] - 2025-01-10

### Added
- **Telegram Authentication**
  - Login/Register via Telegram Login Widget
  - HMAC-SHA256 verification of Telegram auth data
  - JWT token generation with 7-day expiration
  - Authentication middleware for protected routes

- **Database Schema**
  - Users table with Telegram fields
  - Accounts table for different account types
  - Categories table for income/expense categories
  - Transactions table with automatic balance updates
  - Budgets table for category-based budget tracking
  - Indexes for improved query performance

- **Account Management**
  - Create, read, update, delete accounts
  - Support for multiple account types (checking, savings, cash, credit, investment)
  - Multi-currency support (default: RUB)
  - Track initial and current balance
  - Custom colors and icons

- **Category Management**
  - Create, read, update, delete categories
  - Support for income and expense categories
  - Hierarchical categories with parent_id support
  - Custom colors and icons
  - Filter by category type

- **Transaction Management**
  - Create, read, update, delete transactions
  - Support for income, expense, and transfer types
  - Automatic balance updates on create/delete
  - Advanced filtering (by account, category, type, date range, amount)
  - Pagination support
  - Transaction history with timestamps

- **Budget Management**
  - Create, read, update, delete budgets
  - Monthly and yearly budget periods
  - Budget status tracking (spent, remaining, percentage)
  - Budget exceeded detection

- **Statistics & Reports**
  - Category summary with total amounts and transaction counts
  - Monthly balance reports (income vs expense)
  - Date range filtering for reports

- **API Documentation**
  - Comprehensive README with setup instructions
  - Detailed API_DOCS.md with all endpoints and examples
  - .env.example template for configuration

### Security
- JWT-based authentication
- Telegram HMAC verification
- User data isolation
- Environment-based secrets management
- No security vulnerabilities detected (CodeQL verified)

### Technical Details
- Go 1.24.4
- Chi router for HTTP handling
- SQLite database with WAL mode
- JWT tokens with golang-jwt/jwt/v5
- Modular architecture with separate handlers and storage layers
