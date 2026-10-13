
# TUI Chess

A simple pet project with a terminal user-interface that lets you play chess online or in a local network. Some features are still WIP, but most of the functions operate normally


## Screenshots

![app screenshot](https://dummyimage.com/468x300?text=App+Screenshot+Here)


## Features

- Full legal-move checking with check, mate and stalemate detection
- Castling and promotion
- LAN connection
- Legal move highlighting


## WIP

- Draw rules
- En passant
- Offering draw and resigning
- Promoting to pieces other than queen


## Run Locally

Clone the project

```bash
  git clone https://github.com/uvelyr/chess
```

Go to the project directory

```bash
  cd chess
```

Install dependencies

```bash
  go mod tidy
```

Start the server

```bash
  go run server.go
```

Run the client

```bash
  go run client/client.go
```

Players are identified by a unique token given at the start, so it is not required to launch client from different devices!
## How it works

Project consists of two parts: a server and a client. Running a server (server.go) opens HTTP endpoints, which client (client.go) can interact with via JSON.

| Endpoint | Description |
|---|---|
| `POST /joinGame` | joins a player to the game and returns their color and unique idenfication token (the color is randomised) |
| `POST /inputSource` | accepts the coordinates of the piece that is to be moved and server calculates its legal moves |
| `POST /inputTarget` | accepts the coordinates of the square that the piece is to be moved to. In case the move is illegal, returns `invalid` |
| `GET /getFigureBoard` | returns an [8][8] array of integers, that represent runes. They get translated into symbols on the client side |
| `GET /getLegalBoard` | same as /getFigureBoard, but with legal moves of the chosen piece overlayed on top |
| `GET /queryChanged` | returns an answer only when it's the player's turn (or the game is ) |
| `GET /gameCondition` | returns whether the game is in normal condition, someone is in check, checkmate or it's a stalemate |

identifier token is required for all of these endpoints besides the /joinGame