                                ~~~ | LÉEME | ~~~
                                  Laboratorio 2

Integrantes:
José Miguel Quezada | 201773528-7
Ruth Vicuña Vera    | 201673007-9


Supuestos:

  - El usuario ingresará correctamente los datos.
  - El algoritmo a utilizar (centralizado o descentralizado) será definido por input del Cliente Uploader.
  - Si hay un DataNode caído, la nueva propuesta considerará los otros dos nodos activos.
  - Habrá máximo un DataNode caído en cada ejecución.
  - Clente Downloader puede pedir ver la lista de libros disponibles más de una vez seguida.
  - El contador de mensajes para las métricas del informe no considera mensajes con clientes.


Otra información útil:

  - En caso de aparecer un error de permiso denegado, probar ejecutar desde la carpeta Lab2:
    $ chmod -R 777 bin
    
  - En caso de haber problemas de conexión probar ejecutar: 
    $ sudo systemctl stop firewalld
  
  - Para ver capturas del funcionamiento del programa, ingresar a:
    https://drive.google.com/drive/folders/1szN9ojEgKQ_XrMKRxICQFrQqPxE8J1Hu?usp=sharing


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


Instrucciones:

  - En la carpeta Lab2 está todo lo relacionado al Laboratorio 2.
    Allí se encuentran los programas en .go, la carpeta del archivo .proto "book"
    y la carpeta bin, destinada a la ejecución de los programas y almacenamiento 
    del log, los chunks y los libros.

  - Abrir 6 terminales para ejecutar los 4 nodos y 2 clientes y 
    moverse hacia bin según el siguiente comando:

    $ cd SistDistribuidos/Lab2/bin     //desde la carpeta $HOME en las VMs y en las 6 terminales


  - Ejecutar nodos:

VM17:
$ cd NN
$ make 

VM18:
$ cd DNA
$ make

VM19:
$ cd DNB
$ make

VM20:
$ cd DNC
$ make


  - Ejecutar clientes

VM17/18 (Upload):
$ cd CU
$ make

VM19/20 (Download):
$ cd CD
$ make


- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
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