# Technology Stack: Indicator

## Languages & Core Frameworks
- **Go (1.22):** The primary programming language, utilizing generics for type-safe numeric processing and channels for high-performance data streaming.

## Architecture
- **Modular Financial Toolkit:** Decoupled packages for indicators (`trend`, `momentum`, etc.), strategies, asset management, and backtesting.
- **MCP Integration:** Native support for the Multi-Client Protocol, enabling AI-driven interactions and analysis.

## Persistence & Repositories
- **SQL Repository:** Support for relational databases (e.g., PostgreSQL, SQLite) for snapshot persistence.
- **File System Repository:** Efficient storage of asset snapshots in local files.
- **In-Memory Repository:** High-speed, transient storage for real-time analysis.

## Market Data Integrations
- **Tiingo Repository:** Built-in connector for fetching historical and real-time market data from the Tiingo API.
- **Sync Framework:** Concurrent synchronization engine for data ingestion from remote sources to local repositories.

## Testing & Quality Assurance
- **Standard Go Testing:** Unit and integration tests using the built-in `testing` package.
- **Data-Driven Tests:** Custom testing framework using CSV test data for rigorous validation of indicator and strategy outputs.
- **CI/CD:** Automated testing and build pipelines via GitHub Actions.
