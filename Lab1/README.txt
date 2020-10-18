
                                ~~~ | LÉEME | ~~~


Integrantes:
José Miguel Quezada | 201773528-7
Ruth Vicuña         | 201673007-9


Supuestos e información útil:

    ·Se les asignaron los siguientes roles a las máquinas virtuales:
        - MV 17: Servidor/Logística
        - MV 18: Finanzas
        - MV 19: Cliente
        - MV 20: Camiones

    ·Si tienda es de tipo pyme, necesariamente dirá "pyme" en columna tienda en pymes.csv

    ·En la pauta se usan los términos "número de intentos" y de "reintentos", y se interpretó
    que el número total de intentos para retail y pymes es de 3 (además de la restricción del
    valor del producto para pymes)

    ·No se considera el 30% del valor de los paquetes prioritarios al escoger el orden de entrega

    ·Las variables de estado de paquetes son 
        bdg: bodega
        tr: en tránsito
        rec: recibido
        nr: no recibido


Instrucciones:

    Ejecutar en terminal:

    *este es mi caso personal
    export GOROOT=/media/joseesmuyoriginal/opt/go \
    
    export GOROOT=/usr/local/go
    export GOPATH=$HOME/go
    export GOBIN=$GOPATH/bin
    export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN


    protoc -I logis --go_out=./logis --go_opt=paths=source_relative \
    --go-grpc_out=./logis --go-grpc_opt=paths=source_relative logis/logis.proto

    > go build -race -ldflags "-s -w" -o bin/server GRPC_server/server.go
	> bin/server

    > go build -race -ldflags "-s -w" -o bin/cliente GRPC_cliente/cliente.go
	> bin/cliente

    > go build -race -ldflags "-s -w" -o bin/cliente Camiones/camiontb.go
	> bin/camiontb

    > go build -race -ldflags "-s -w" -o bin/financiero financiero.go
	> bin/financiero