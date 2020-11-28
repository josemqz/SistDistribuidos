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
- 


Supuestos:

- Hará un máximo de 100 datanodes en espera para la exclusión mutua centralizada.
- El algoritmo a utilizar (centralizado o distribuido) será definido por input del usuario uploader.
- Si hay un datanode caído, la nueva propuesta considerará los otros dos nodos activos.
- (****)''''''''Habrá máximo un datanode caído en cada ejecución.

Otra información útil:


Instrucciones:

    VM17: NameNode
    VM18: DataNodeA
    VM19: DataNodeB
    VM20: DataNodeC



