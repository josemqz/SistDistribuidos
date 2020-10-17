
                                ~~~ | LÉEME | ~~~

Integrantes:
José Miguel Quezada | 201773528-7
Ruth Vicuña         | 201673007-9

Supuestos e información útil:
    ·Si tienda es de tipo pyme, necesariamente dice "pyme" en columna tienda en pymes.csv
    ·En la pauta se usan los términos "número de intentos" y de "reintentos", y se interpretó
    que el número total de intentos para retail y pymes es de 3 (además de la restricción del
    valor del producto para pymes)
    ·No se considera el 30% del valor de los paquetes prioritarios al escoger el orden de entrega
    ·Las variables de estado de paquetes son 
        bdg: bodega, tr: en tránsito , rec: recibido , nr: no recibido

Instrucciones:

    Ejecutar en terminal:

    export GOROOT=/media/joseesmuyoriginal/opt/go \  *esto es mi caso específico
    export GOPATH=$HOME/go \
    export GOBIN=$GOPATH/bin \
    export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN

    protoc -I logis --go_out=./logis --go_opt=paths=source_relative \
    --go-grpc_out=./logis --go-grpc_opt=paths=source_relative logis/logis.proto

    go build -race -ldflags "-s -w" -o bin/server server/main.go
	bin/server

    go build -race -ldflags "-s -w" -o bin/cliente cliente/main.go
	bin/cliente