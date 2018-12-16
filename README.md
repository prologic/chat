# chat

A simple terminal chat with e2e twofish or blowfish encryption over UDP

## Features

-[ ] Diffie Hellman key exchange
- [ ] Blowfish and Twofish e2e encryption
- [ ] UDP transport with p2p support
- [ ] Simple Terminal UI

## Installation

```#!bash
$ go get -u github.com/prologic/chat
```

## Usage

Waiting for a peer:

```#!bash
chat
```

By default binds to UDP port `:1337`.

Connecting to a peer:

```#!bash
chat x.x.x.x:1337
```

## License

MIT
