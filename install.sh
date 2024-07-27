echo "Cloning repository..."
git clone https://github.com/Ashy5000/cryptocurrency
cd cryptocurrency
echo "Building node software... ('consensus client')"
go build -o builds/node/node
echo "Building smart contract software... ('execution client')"
cd contracts
cargo build --release
cd ..