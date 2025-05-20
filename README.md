# grip2-cli

Command line client for interfacing with GRIP2 Panel. A rewrite of the original 
[grip-cli](https://bitbucket.resideo.com/projects/GRIP/repos/grip-cli/browse) program from Python to Go. 


## Quirks
This program is entirely written in Go's standard library, significantly reducing dependency overhead and resulting binary size 
(compressed and stripped can be down to 500MB). Theoretically, software maintinance should be easier as all code present is what
will run (assuming this implementation is satisfactory for the usecase and is able to implement the necessary commands with relative ease).

## Building
```
./build.sh
```
Resulting binary `cli` will be created.



Can also provide a couple arguments to `build.sh` like `stripped` (for smaller resulting executable, will try to call `upx` if available
on the system) and `deploy` (will compile for `arm64`)


eg.
```
./build.sh stripped deploy
```

## TODO (as of #24cf72)
1. [ ] Implement basic mqtt client
    [ ] Test by implementing `armdisarm`
2. [ ] Implement some `getinfo` or equivilant command
3. [ ] Implement `trouble`

## Current Issues
None yet...
