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
budget_tui/
├── cmd/
│   ├── budget_tui.go           # Main application entry point (simplified version)
│   └── budget_tui_simple.go    # Clean main file for new structure
├── internal/
│   ├── budget/
│   │   ├── budget.go           # Budget data model and operations
│   │   └── storage.go          # Data persistence (moved from root)
│   ├── import/
│   │   ├── csv_parser.go       # CSV import and format detection
│   │   └── import_history.go   # Import session tracking
│   ├── tui/
│   │   ├── banner.go           # ASCII art banner
│   │   ├── dashboard.go        # Dashboard rendering
│   │   ├── styles.go           # Original styling (has conflicts)
│   │   └── components.go       # Clean TUI components
│   ├── storage/
│   │   └── storage.go          # Storage utilities
│   └── analytics/
│       └── analytics.go        # Budget analytics
├── pkg/
│   └── categorizer/
│       └── categorizer.go      # Smart categorization engine
├── configs/                    # Configuration files
├── docs/                       # Documentation
├── testdata/
│   └── sample_bank_statement.csv  # Demo data
├── go.mod                      # Go module file
├── go.sum                      # Go dependencies
└── README.md                   # This file
```

### Running the Application

**Current Working Version (Flat Structure):**
```bash
# Run the original working version from root
go run main.go
```

**New Restructured Version (In Progress):**
```bash
# Run the restructured version
cd cmd && go run budget_tui_simple.go
```

*Note: The directory restructure is in progress. The main working application is still `main.go` in the root directory.*

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