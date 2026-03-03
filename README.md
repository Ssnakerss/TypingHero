# Typing Hero

Console-based typing trainer in Go with difficulty levels and typing statistics.

## Features

- Console-based typing trainer
- 10 difficulty levels
- Real-time typing speed calculation (CPM)
- Error counting
- AI-generated texts (via GigaChat API - to be implemented)

## Usage

```bash
go run main.go
```

## Implementation Details

The program currently uses predefined texts for different difficulty levels. In the future, it will be enhanced to use GigaChat API for text generation.

## Planned Enhancements

- Integration with GigaChat API for dynamic text generation
- Improved input handling (support for spaces and special characters)
- Statistics persistence
- Multiple language support
- Typing lessons progression system