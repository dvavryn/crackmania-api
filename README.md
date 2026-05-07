*This is part of the Project ft_transcendence of the 42 curriculum*<br/>
*API author: dvavryn(Dominic Vavryn)*<br/>
*Project members: dvavryn, hanjkim, dplotzl, oohnivch*<br/>
<hr/>

# GOTTA-GO-FAST-API

## Description
The gotta-go-fast-api is and AI that calculates the controll inputs of  the AI opponent for the web-application gotta-go-fast aka. Crackmania.

## Usage
```bash
./gotta-go-fast-api.bin
```
Server needs to be started before the game can be started!
Start API -> Open Game -> look what is happening

If nothing is happening, press <b>F12</b> in the browser, open Console, look what is happening.
Also look at the logs.
Recompile in debug mode to see better debug messages!


### Makefile
||description|
|---|---|
|make|compiles main.go to .bin|
|make all|similar to make|
|make build|recompiles main.go to .bin|
|make run|compiles if .bin missing and executes current state of .bin|
|make debug|compiles .bin in debug mode|
|make clean|deletes .bin|

## Tree
.
├── Dockerfile
├── src
│   ├── ai
│   │   ├── ai.go
│   │   └── structs.go
│   ├── config
│   │   ├── config.go
│   │   └── structs.go
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   └── Makefile
└── config.json
