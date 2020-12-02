                                ~~~ | LÉEME | ~~~
                                  Laboratorio 2

Integrantes:
José Miguel Quezada | 201773528-7
Ruth Vicuña Vera    | 201673007-9

Pendientes:

- El Name Node se ejecuta en una máquina sin Data Nodes (???)
- Exclusión mutua
- Crear copias DataNode (B y C)
- Mutex en exclusión mutua
- Hacer makefiles
- Probaaaaaaaaaar
- Informe


Supuestos:

  - Puede haber hasta un máximo de 100 DataNodes en espera para la exclusión mutua centralizada.
  - El algoritmo a utilizar (centralizado o descentralizado) será definido por input del Cliente Uploader.
  - Si hay un DataNode caído, la nueva propuesta considerará los otros dos nodos activos.
  - Habrá máximo un DataNode caído en cada ejecución.
  - Clente Downloader puede pedir ver la lista de libros disponibles más de una vez seguida
  - El contador de mensajes para las métricas del informe no considera mensajes con clientes.

Otra información útil:


Instrucciones:

DATOS VM:

    VM17: NameNode + Cliente Uploader
    VM18: DataNodeA + Cliente Uploader
    VM19: DataNodeB + Cliente Downloader
    VM20: DataNodeC + Cliente Downloader

    Máquina Namenode + Cliente Uploader
        hostname:dist17
        contraseña: EPCwrFS4
        IP: 10.6.40.157

    Máquina DataNodeA + Cliente Uploader
        hostname:dist18
        contraseña: L8m9s7AS
        IP: 10.6.40.158

    Máquina DataNodeB + Cliente Downloader
        hostname:dist19
        contraseña: VYgDPNJe
        IP: 10.6.40.159

    Máquina DataNodeC + Cliente Downloader
        hostname:dist20
        contraseña: Kcf25KUB
        IP: 10.6.40.160

Puertos:

  DNA - DNB 50500
  DNB - DNA 50501

  DNA - DNC 50502
  DNC - DNA 50503

  DNB - DNC 50504
  DNC - DNB 50505

  DNA - NN  50506
  DNB - NN  50507
  DNC - NN  50508

  NN - DNA  50509
  NN - DNB  50510
  NN - DNC  50511

  CD - NN   50512
  CD - DNA  50513
  CD - DNB  50514
  CD - DNC  50515

  CU - DNA  50517
  CU - DNB  50518
  CU - DNC  50519


    $ export GO111MODULE=on
    $ export GOROOT=/media/joseesmuyoriginal/opt/go
    $ export GOPATH=$HOME/go
    $ export GOBIN=$GOPATH/bin
    $ export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN:$GOROOT/bin

    $ protoc -I book book/book.proto --go_out=./book --go-grpc_out=./book --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative

    $ make...
        go build -race -ldflags "-s -w" -o bin/DNA DNA.go


