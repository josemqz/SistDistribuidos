
                                ~~~ | LÉEME | ~~~

Integrantes:
José Miguel Quezada | 201773528-7
Ruth Vicuña         |

Supuestos e información útil:
    ·Si tienda es de tipo pyme, necesariamente dice "pyme" en columna tienda en pymes.csv
    ·Las variables de estado de paquetes son 
        bdg: bodega, tr: en tránsito , rec: recibido , nr: no recibido

Instrucciones:

    export GOROOT=/media/joseesmuyoriginal/opt/go \  *esto es mi caso específico
    export GOPATH=$HOME/go \
    export GOBIN=$GOPATH/bin \
    export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN

    protoc -I logis --go_out=./logis --go_opt=paths=source_relative \
    --go-grpc_out=./logis --go-grpc_opt=paths=source_relative \
    logis/logis.proto

    go build -race -ldflags "-s -w" -o bin/server server/main.go *en las MVs deberían estar los mains directamente
	bin/server