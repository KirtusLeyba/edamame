# edamame
Network visualization using golang and raylib.

### building from source
```bash
git clone https://github.com/KirtusLeyba/edamame
cd edamame
make
```

### Installing in the go workspace
```bash
go install https://github.com/KirtusLeyba/edamame
```

### Interactive Mode
Recommended for networks with < 2000 nodes.
Use the GUI in interactive mode:
```bash
./edamame
```

<!--	Headless bool
	NodeFilePath, EdgeFilePath, OutputFilePath string
	MaxWorkers, MaxIters int
	Repulsion float64-->

### Headless Mode
Use headless mode with large networks (> 2000 nodes).
```bash
./edamame -headless -nodeFilePath path-to-node-csv -edgeFilePath path-to-edge-csv -outputFilePath path-to-save-img -maxWorkers number-of-go-routines-to-use -maxIters number-of-iters-for-layout-algorithm -repulsion repulsive-force-in-layout-algorithm
```

### Licence
edamame is open-source software licensed according to the MIT license.
See [license](https://github.com/KirtusLeyba/edamame/LICENSE.md)
