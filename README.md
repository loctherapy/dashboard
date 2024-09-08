# Dashboard CLI
`dashboard` is a command-line interface (CLI) application written in Go, designed to help you manage and navigate TODOs in your markdown notes. The app groups and sorts undone TODOs by contexts and projects based on frontmatter values in the files.

## Features
- **Grouping and Sorting**: TODOs are grouped by contexts and sorted by project gravity.
- **Navigation**: Use number keys to navigate between contexts and projects.
- **TODO Management**: Mark TODOs as completed directly in your markdown files.

## Installation
To install dashboard, clone the repository and build the Go application:
```
git clone https://github.com/loctherapy/dashboard
cd dashboard
make build
make install
```

## Usage
Run the dashboard CLI in the terminal within your IDE where your repository with notes is opened:
```
./dashboard
```

### Example Note File
Consider a file project_alfa.md:
```
---
context: myself-1
gravity: 8
---

# Project Alfa
- [ ] Do A
- [ ] Do B
- [ ] Do C
```

### Frontmatter Explanation
- `context`: Sets the priority for the context. The lower the number, the higher the priority.
- `gravity`: Indicates the importance of the project. The higher the gravity, the higher the project will appear in the list of TODOs within the context.

### Navigating TODOs
- **Contexts**: Use number keys (2, 3, 4, etc.) to navigate between contexts. The 1 key shows TODOs from all contexts.
- **Projects**: Click on the project name to navigate to the file in your IDE where you can work on the TODOs.

### Marking TODOs as Completed
To mark a TODO as completed, change the status from - [ ] to - [x] in the markdown file. The dashboard doesn't show completed TODOs.

### License
This project is licensed under the MIT License. See the LICENSE file for details.

### Contributing
Contributions are welcome! Please open an issue or submit a pull request.

