# Typing Hero

Console and web-based typing trainer in Go with difficulty levels and typing statistics.

## Features

- **Console-based typing trainer** with colorful ANSI output
- **Web-based interface** accessible at http://localhost:8080
- 10 difficulty levels with progressively complex texts
- Real-time typing speed calculation (WPM)
- Error counting and accuracy percentage
- Color-coded visual feedback for typing performance
- User statistics tracking (best, worst, average speeds)
- Persistent storage of typing sessions in SQLite database
- Context-aware graceful shutdown handling

## Usage

### Console Version

```bash
go run cmd/main.go
```

### Web Version

The web server automatically starts alongside the console version at http://localhost:8080

## Architecture

Typing Hero is a dual-mode application with both console and web interfaces that run concurrently:

- **Main Application** (`cmd/main.go`): Initializes database, handles signals, and runs both console and web versions concurrently
- **Console Module** (`internal/console/console.go`): Interactive terminal interface with colorful output and statistics
- **Web Module** (`internal/web/web.go`): HTTP server providing web interface
- **Database Layer** (`internal/database/database.go`): SQLite persistence for typing sessions
- **Models** (`internal/models/`): Data structures and text pools for different difficulty levels

## Implementation Details

The program uses predefined texts for different difficulty levels, stored in `internal/models/text.go`. Each difficulty level (1-10) has its own pool of texts that increase in length and complexity.

User typing sessions are stored in an SQLite database (`data/typing_sessions.sqlite`) with the following schema:

- **users** table: Stores user information (id, name, nickname, timestamps)
- **typing_sessions** table: Stores session results (user_id, difficulty, wpm, error_rate, time_taken, timestamps)

The application uses Go's context package for graceful shutdown and concurrent execution of console and web interfaces.

## Planned Enhancements

- User authentication and personalized statistics
- Web interface for viewing statistics and progress
- Dynamic text generation via API integration
- Multiple language support
- Typing lessons progression system
- Improved web UI with real-time typing feedback
- Mobile-responsive design for web version
- Export/import of typing statistics
- Achievement system and gamification elements