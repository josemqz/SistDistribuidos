
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

    ·Las variables de estado de paquetes son 
        bdg: bodega
        tr: en tránsito
        rec: recibido
        nr: no recibido

    ·Si tienda es de tipo pyme, necesariamente dirá "pyme" en columna tienda en pymes.csv

    ·En la pauta se usan los términos "número de intentos" y de "reintentos", y se interpretó
    que el número total de intentos para retail y pymes es de 3 (además de la restricción del
    valor del producto para pymes)

    ·No se considera el 30% del valor de los paquetes prioritarios al escoger el orden de entrega

    ·En los registros de camiones se realiza "append" de los datos actualizados de cada pedido, 
    por lo que cada ID aparece dos veces.


Instrucciones:

    Una vez abierta la terminal ejecutar:

    $ cd SistDistribuidos/Lab1

    $ export GOROOT=/usr/local/go \
    export GOPATH=$HOME/go \
    export GOBIN=$GOPATH/bin \
    export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN


    --VM------------------------
    $ go build -race -ldflags "-s -w" -o bin/server VM/serverVM.go
	$ bin/server

    $ go build -race -ldflags "-s -w" -o bin/cliente VM/clienteVM.go
	$ bin/cliente

    $ go build -race -ldflags "-s -w" -o bin/camion VM/camionVM.go
	$ bin/camion

    $ go build -race -ldflags "-s -w" -o bin/financiero VM/financieroVM.go
	$ bin/financiero


DATOS VM:

    Máquina Server/Logística
    hostname:dist17
    contraseña:EPCwrFS4

    Máquina Financiero
    hostname:dist18
    contraseña:L8m9s7AS

    Máquina Cliente
    hostname:dist19
    contraseña:VYgDPNJe

    Máquina Camiones
    hostname:dist20
    contraseña:Kcf25KUB


Dirección Logístico
    10.6.40.157


Para errores VM:

    go get google.golang.org/protobuf/cmd/protoc-gen-go
    go get google.golang.org/grpc
    go get github.com/golang/protobuf

    export GO111MODULE=on  # Enable module mode
    $ go get github.com/golang/protobuf/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc
