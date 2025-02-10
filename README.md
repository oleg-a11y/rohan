# Notification Service

This project is an interview notification service that integrates with Notion and Telegram. The service retrieves information about interviews from a Notion database and sends notifications to Telegram.

## Installation

To install the dependencies, run the following command:

```bash
make install
```

## Running the Project

To run the project in development mode, use the following command:

```bash
make run
```

## Building the Project

To build the binary file, run the command:

```bash
make build
```

The built file will be located in the `bin` directory and will be named `app.exe`.

## Makefile Targets

- **run**: Runs the project directly in development mode.
- **build**: Builds the binary and outputs it as `app.exe` in the `bin` directory.
- **install**: Installs the dependencies by running `go mod tidy`.

## Usage

The service automatically performs the following tasks:

1. **Every day at 10:00 AM**: Retrieves all interviews for the current day and sends them to Telegram.
2. **Every minute**: Checks for interviews that will start in 10 minutes and sends notifications to Telegram.

## Configuration

Before running the service, ensure that you have correctly set the following environment variables in the configuration file:

- `NotionAPIKey`: The API key for accessing Notion.
- `NotionDatabaseID`: The ID of the Notion database from which interviews will be retrieved.
- `TelegramBotToken`: The token for your Telegram bot.
- `TelegramChatID`: The ID of the Telegram chat where notifications will be sent.
- `TelegramThreadID`: The ID of the Telegram thread where notifications will be sent.
