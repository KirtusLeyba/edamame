![edamame logo](https://github.com/KirtusLeyba/edamame/blob/main/edamame_logo.png)

# edamame
Network visualization using golang and raylib.

### Screenshots
Example using an Erdos-Renyi random graph:
![random graph example](https://github.com/KirtusLeyba/edamame/blob/main/edamame_screenshot.png)

### building from source
```bash
git clone https://github.com/KirtusLeyba/edamame
cd edamame
# Get the raylib-go bindings if you don't have them
go get github.com/gen2brain/raylib-go/raygui
# compile
make
```

### Installing the latest version with go install
```bash
go install github.com/KirtusLeyba/edamame@latest
```

### Interactive Mode
Recommended for networks with < 2000 nodes.
Use the GUI in interactive mode:
```bash
edamame
```

### Headless Mode
Use headless mode with large networks (> 2000 nodes).
```bash
edamame -headless \
-nodeFilePath path-to-node-csv \
-edgeFilePath path-to-edge-csv \
-outputFilePath path-to-save-img \
-maxWorkers number-of-go-routines-to-use \
-maxIters number-of-iters-for-layout-algorithm \
-repulsion repulsive-force-in-layout-algorithm
```

### Licence
edamame is open-source software licensed according to the MIT license.
See [license](https://github.com/KirtusLeyba/edamame/blob/main/LICENSE)
