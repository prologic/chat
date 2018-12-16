# chat

A simple terminal chat with e2e twofish or blowfish encryption over UDP

## Features

- [ ] Diffie Hellman key exchange
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

## Sample

Here is a sample session between two clients (//on localhost//) to demonstrate
the exchange as well as the packets on the wire captured via `ngrep`:

Listening client:

```#!bash
$ ./chat
>>> <prologic> hello
^D
```

Connecting client:

```#!bash
$ ./chat -b :1338 127.0.0.1:1337
>>> hello
>>> ^D
```

Packet dump:

```#!bash
$ ngrep -q -W byline -d lo0 '.*' 'udp port 1337'
interface: lo0 (127.0.0.0/255.0.0.0)
filter: ( udp port 1337 ) and (ip || ip6)
match: .*

U 127.0.0.1:1338 -> 127.0.0.1:1337 #1
{"kind":1,"addr":"","user":"prologic","data":""}

U 127.0.0.1:1338 -> 127.0.0.1:1337 #2
{"kind":2,"addr":"","user":"prologic","data":"xQb8sVksIlv701Xc0qeYZdkFMLNBgjPhkqzY1fQf1B5vSzERJ0Y6JzbEeqHSaDIdBrK3P6M8PNcVpRd5SDDNeIMz5E7lShsx4h3JXp7Z6YsmDVSj+f+DmKF+nm49g20uUICMzoAJ5jNiD+anqtvz847W6fdY2WGiS0gbF37b/4ZD2L47XwmMIL+Gumu/9xSH/2TC7uSbpfDjNdF2+kJtZIToBSvC4ZEVEPceSsg+/gHjbx7WUlU2Gk55eGRHhLxVyxRjZslCSFv6nH+ypr6++VigwtP/emCX6Ow1Cm0p6g/kHbW6kos3srMHTiBvsASn2D3+G3syNZTCkQ0YhDWWaw=="}

U 127.0.0.1:1337 -> 127.0.0.1:1338 #3
{"kind":2,"addr":"","user":"prologic","data":"A0kuKyl5aFh+voSng4SGfE/BLQSHlhrQwv4V7q8SFgg6vTzWRkLMvauY6siIvO+1h5USRaf9mLBofWIw7EpKgipZLsugSLj9hTBaxh0QFu0B4lEGQ3QyYOkJYwlITlmwO6vzV8saXaaromdL26H3FfWcxGIjvEgPTZnUELF0EXHiNHqSAc/X0k7BGeWDXHwaOJktqgVUvj267ai0hsDP2UcWFVkwnNiwHQgzudAQOGsqTILtd0D3ozDJ6mF7F2pwStpSpf5KzKVbcvFJSQr7GuBepMhHjWrzlOhAjP74I6D6gp7kWA/8sxQCKed48hHj/b82+DW1LQWFoIdQeKgN9g=="}

U 127.0.0.1:1338 -> 127.0.0.1:1337 #4
{"kind":2,"addr":"","user":"prologic","data":"xQb8sVksIlv701Xc0qeYZdkFMLNBgjPhkqzY1fQf1B5vSzERJ0Y6JzbEeqHSaDIdBrK3P6M8PNcVpRd5SDDNeIMz5E7lShsx4h3JXp7Z6YsmDVSj+f+DmKF+nm49g20uUICMzoAJ5jNiD+anqtvz847W6fdY2WGiS0gbF37b/4ZD2L47XwmMIL+Gumu/9xSH/2TC7uSbpfDjNdF2+kJtZIToBSvC4ZEVEPceSsg+/gHjbx7WUlU2Gk55eGRHhLxVyxRjZslCSFv6nH+ypr6++VigwtP/emCX6Ow1Cm0p6g/kHbW6kos3srMHTiBvsASn2D3+G3syNZTCkQ0YhDWWaw=="}

U 127.0.0.1:1338 -> 127.0.0.1:1337 #5
{"kind":0,"addr":"","user":"prologic","data":"AAAAAAAAAAAAAAAAAAAAAMstXJjAO53FHrxEXdaMLQc="}
^C
```
## License

MIT
