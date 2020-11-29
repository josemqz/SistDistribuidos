                                ~~~ | LÉEME | ~~~
                                  Laboratorio 2

Integrantes:
José Miguel Quezada | 201773528-7
Ruth Vicuña Vera    | 201673007-9

Pendientes:

- Datanodes:  Recibe y almacena correctamente los chunks del datanode alpha
- C Down: Solicita (a los Data Nodes) y guarda todos los chunks que componen a un libro mediante Protocol Buffers
- El Name Node se ejecuta en una máquina sin Data Nodes (???)
- Exclusión mutua
- Conexion debe ser bilateral?
- Verificar y cambiar puertos en cada función de conexión


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

  CU - NN   50516
  CU - DNA  50517
  CU - DNB  50518
  CU - DNC  50519


Supuestos:

  - Habrá un máximo de 100 DataNodes en espera para la exclusión mutua centralizada.
  - El algoritmo a utilizar (centralizado o distribuido) será definido por input del Cliente Uploader.
  - Si hay un datanode caído, la nueva propuesta considerará los otros dos nodos activos.
  - (****)''''''''Habrá máximo un datanode caído en cada ejecución.


Otra información útil:

  - Cliente Uploader se ejecutará en las máquinas virtuales 17 y 18 
  - Cliente Downloader se ejecutará en las máquinas virtuales 19 y 20
  
  - Clente Downloader puede pedir ver la lista de libros disponibles más de una vez seguida


Instrucciones:

    VM17: NameNode + Cliente Uploader
    VM18: DataNodeA + Cliente Uploader
    VM19: DataNodeB + Cliente Downloader
    VM20: DataNodeC + Cliente Downloader

    $ protoc -I testp testp/testp.proto --go_out=./testp --go-grpc_out=./testp --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative

    $ make...




