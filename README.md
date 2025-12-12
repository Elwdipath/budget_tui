# Budget TUI - Terminal Budget Application

A powerful terminal-based budget management tool with bank statement import and smart categorization.

## Features

- 🏠 **Dashboard** - Real-time financial overview with ASCII art banner
- 📁 **Bank Import** - CSV bank statement import with automatic format detection
- 🏷️ **Smart Categorization** - 40+ rules for automatic transaction categorization
- 📊 **Spending Analytics** - Category breakdown and transaction history
- 💾 **Data Persistence** - JSON-based storage with import history tracking

## Quick Start

```bash
# Build and run
go run *.go

# Or build binary
go build .
./budget_tui
```

## Usage

### Dashboard Navigation
- **j/k** - Navigate categories and transactions
- **i** - Add income
- **e** - Add expense  
- **t** - View all transactions
- **b** - Import bank statement
- **h** - Toggle help
- **q** - Quit

### Bank Statement Import
1. Press `[b]` from dashboard
2. Press `Tab` to select sample file
3. Press `Enter` to detect format and preview
4. Press `c` to confirm import

### Supported Bank Formats
- Chase
- Bank of America  
- Wells Fargo
- Generic CSV format

## Project Structure

```
├── main.go              # Main application and UI
├── budget.go            # Budget data model and operations
├── storage.go           # Data persistence
├── csv_parser.go        # CSV import and format detection
├── categorizer.go       # Smart categorization engine
├── import_history.go    # Import session tracking
├── dashboard.go         # Dashboard rendering
├── banner.go            # ASCII art banner
├── styles.go            # Lipgloss styling
├── analytics.go         # Budget analytics
└── sample_bank_statement.csv  # Demo data
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling

## Development

```bash
# Install dependencies
go mod tidy

# Run with development output
go run *.go

# Build for production
go build -o budget_tui
```

## License

MIT License