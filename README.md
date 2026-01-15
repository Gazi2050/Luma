# Luma – A Minimal Programming Language

**Luma** is a high-performance, lightweight, interpreted programming language built from scratch in **Go**. It aims to provide a clean, modern syntax while maintaining the simplicity of a tree-walking interpreter.

---

## 🚀 Installation

Ensure you have [Go](https://go.dev/) (version 1.18+) installed on your system.

### 1. Clone & Setup
```bash
# Clone the repository
git clone <your-repo-url>
cd Luma

# Initialize and tidy dependencies
go mod tidy
```

### 2. Build the Binary
```bash
# Build the luma executable
go build -o luma main.go
```

The `luma` binary is now ready for use in your current directory.

---

## 🕹️ Usage

Luma provides two primary environments for development.

### 1. Interactive REPL
Perfect for testing expressions and small logic snippets.

**Command:**
```bash
./luma
```

**Inside the REPL:**
- **Execute**: Type any code and press `Enter`.
- **Persistent State**: Variables declared in the REPL remain in memory until you exit.
- **Exit**: Use `Ctrl+C` or `Ctrl+D`.

### 2. Script Execution
For running modular `.lu` source files.

**Command:**
```bash
./luma run <path_to_file>.lu
```

**Example:**
```bash
./luma run index.lu
```

---

## ✨ Features

- **Robust Variable System**: Supports both mutable `let` and immutable `const` declarations.
- **Rich Scoping**: Lexical scoping ensures variables are only accessible where they should be.
- **First-Class Functions**: Define, pass, and return functions just like any other value.
- **Modern Data Structures**:
    - **Objects**: Key-value pairs with dot notation (`user.name`).
    - **Arrays**: Dynamic lists with index access (`nums[0]`).
- **Dynamic Typing**: No need for explicit type declarations; Luma handles it at runtime.
- **Standard Library**:
    - `log(...)`: Output any value to the console with clean formatting.
    - `len(...)`: Get the length of strings or arrays instantly.

---

## 📜 Development Guidelines

### ✅ What to Use (Dos)
- **Use `const` by default**: Favor `const` for values that shouldn't change to prevent accidental reassignments.
- **Use Semicolons**: While often optional, ending statements with `;` is best practice for clarity.
- **Leverage Objects**: Group related data into objects to keep your code organized.
- **Use `log()` for Debugging**: It pretty-prints complex structures like nested objects and arrays automatically.

### ❌ What to Avoid (Don'ts)
- **Don't Overuse Global Scope**: Keep variables inside functions to avoid naming collisions.
- **Don't Forget `run` Keyword**: When executing files, always use `./luma run script.lu`. Using `./luma script.lu` will just open the REPL.
- **Avoid Circular References**: While Luma is powerful, deeply circular object references might lead to infinite formatting loops in `log()`.

---

## 📖 Deep Dive
To understand how Luma translates text into logic, check out the implementation details:
- [**The Luma Engine (LANGUAGE_DETAILS.md)**](file:///home/gazi/Projects/my-projects/Luma/LANGUAGE_DETAILS.md)

---

## 📝 Example `index.lu`
```lu
let greeting = "Hello, Luma!";
const PI = 3.14;

fn calculateCircle(r) {
    if (r <= 0) { return 0; }
    return PI * r * r;
}

let data = {
    radius: 10,
    result: calculateCircle(10)
};

log(greeting, "Area:", data.result);
```
