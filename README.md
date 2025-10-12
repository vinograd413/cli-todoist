# CLI Todoist

**CLI Todoist** is a fast, colorful, and interactive command-line tool for managing your personal tasks.
Add, list, update, and delete your todos right from your terminal, with persistent storage using BoltDB.

---

## Features

- **Add tasks** with custom headers
- **List all tasks** with creation timestamps
- **Update tasks** interactively
- **Delete tasks** (supports multi-select)
- **Keyboard navigation** (arrow keys, Enter, Esc)
- **Colorful UI** for better readability
- **Cross-platform** (works on Linux, macOS, Windows)
- **Local persistent storage** (BoltDB)

---

## Quick Start

1. **Clone the repository:**
   ```sh
   git clone https://github.com/yourusername/cli_todoist.git
   cd cli_todoist
2. Build the app:**
   ```sh
   go build -o cli_todoist ./cmd/cli_todoist
   ```

3. **Run the app:**
   ```sh
   ./cli_todoist
   ```

---

## Usage

- Use **↑/↓** to navigate the menu.
- Press **Enter** to select an option.
- Press **Esc** to exit.
- When adding/updating tasks, type your input and press Enter.
- To delete multiple tasks, enter their numbers separated by commas (e.g., `1,2,3`).

---

## Data Storage

All tasks are stored locally in a BoltDB file (`dbFile.db`).
No data is sent outside your machine.

---

## Contributing

Pull requests, issues, and suggestions are welcome!
Feel free to fork the project and make it your own.

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

## Author

Made with ❤️ by vinograd413
