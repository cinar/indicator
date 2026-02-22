# Product Guidelines: Indicator

## Core Principles
- **Reliability:** All financial calculations must be accurate and well-tested.
- **Performance:** Indicators and strategies must be efficient, especially when processing large datasets.
- **Maintainability:** Code should be modular, documented, and easy to extend.
- **Transparency:** All algorithms and logic should be clear and well-documented.

## Prose & Documentation
- **Clear & Concise:** Use simple language to explain complex financial concepts.
- **Technical Accuracy:** Ensure all technical terms and financial formulas are used correctly.
- **Consistent Terminology:** Use consistent naming conventions across the library and documentation.
- **Comprehensive Examples:** Provide clear, working examples for all indicators and strategies.

## User Experience (API Design)
- **Idiomatic Go:** Follow standard Go idioms and best practices for API design.
- **Fluent APIs:** Prefer clear and intuitive interfaces for configuring indicators and strategies.
- **Error Handling:** Provide descriptive error messages to help users diagnose issues.
- **Configurability:** All indicators and strategies should be fully configurable with sensible defaults (if applicable).

## Testing & Quality
- **High Coverage:** Target at least 90% code coverage for all indicators and strategies.
- **Data-Driven Tests:** Use CSV-based test data for validating calculations.
- **Regression Testing:** Ensure new changes don't break existing functionality.
