-- Users table with Telegram authentication
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	telegram_id INTEGER UNIQUE,
	telegram_username TEXT,
	first_name TEXT,
	last_name TEXT,
	photo_url TEXT,
	auth_date INTEGER,
	hash TEXT,
	name TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Accounts table (checking, savings, cash, credit cards, etc.)
CREATE TABLE IF NOT EXISTS accounts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	type TEXT NOT NULL CHECK(type IN ('checking', 'savings', 'cash', 'credit', 'investment', 'other')),
	currency TEXT NOT NULL DEFAULT 'RUB',
	initial_balance REAL NOT NULL DEFAULT 0,
	current_balance REAL NOT NULL DEFAULT 0,
	color TEXT,
	icon TEXT,
	is_active INTEGER DEFAULT 1,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Categories table (income/expense categories)
CREATE TABLE IF NOT EXISTS categories (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	type TEXT NOT NULL CHECK(type IN ('income', 'expense')),
	color TEXT,
	icon TEXT,
	parent_id INTEGER,
	is_active INTEGER DEFAULT 1,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL
);

-- Transactions table
CREATE TABLE IF NOT EXISTS transactions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	account_id INTEGER NOT NULL,
	category_id INTEGER,
	type TEXT NOT NULL CHECK(type IN ('income', 'expense', 'transfer')),
	amount REAL NOT NULL,
	currency TEXT NOT NULL DEFAULT 'RUB',
	description TEXT,
	transaction_date DATETIME NOT NULL,
	to_account_id INTEGER,
	notes TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
	FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
	FOREIGN KEY (to_account_id) REFERENCES accounts(id) ON DELETE SET NULL
);

-- Budgets table (monthly/yearly budgets per category)
CREATE TABLE IF NOT EXISTS budgets (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	category_id INTEGER NOT NULL,
	amount REAL NOT NULL,
	period TEXT NOT NULL CHECK(period IN ('monthly', 'yearly')),
	start_date DATETIME NOT NULL,
	end_date DATETIME,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_categories_user_id ON categories(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions(category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(transaction_date);
CREATE INDEX IF NOT EXISTS idx_budgets_user_id ON budgets(user_id);