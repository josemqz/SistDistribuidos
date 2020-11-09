
                                ~~~ | LÉEME | ~~~
                                  Laboratorio 1

Integrantes:
José Miguel Quezada | 201773528-7
Ruth Vicuña Vera    | 201673007-9


Supuestos e información útil:

    · Se les asignaron los siguientes roles a las máquinas virtuales:
        - MV 17: Servidor/Logística
        - MV 18: Financiero
        - MV 19: Cliente
        - MV 20: Camiones

    · Las variables de estado de paquetes son 
        bdg: bodega
        tr: en tránsito
        rec: recibido
        nr: no recibido

    · Si tienda es de tipo pyme, necesariamente dirá "pyme" en columna tienda en pymes.csv.

    · En la pauta se usan los términos "número de intentos" y de "reintentos", y se interpretó
    que el número total de intentos para retail y pymes es de 3 (además de la restricción del
    valor del producto para pymes).

    · No se considera el 30% del valor de los paquetes prioritarios al escoger el orden de entrega.

    · En los registros de camiones se realiza "append" de los datos actualizados de cada pedido, 
    por lo que cada ID aparece dos veces.

    · Si se realizan los pedidos de pyme y luego los de retail, estos últimos quedarán con id
    de seguimiento igual a la cantidad de pedidos pyme, por un motivo que no pudimos esclarecer.

    · Al realizar ambos tipos de pedidos aparecen ocasionalmente paquetes vacíos, sin motivo 
    conocido. Sin embargo, es algo que no afecta al funcionamiento requerido ni a lo especificado. 

    · Para imprimir el balance final de Financiero, al terminar los envíos, se debe presionar 
    ctrl + C 
    

Instrucciones:

DATOS VM:

    Máquina Server/Logística
        hostname:dist17
        contraseña: EPCwrFS4

    Máquina Financiero
        hostname:dist18
        contraseña: L8m9s7AS

    Máquina Cliente
        hostname:dist19
        contraseña: VYgDPNJe

    Máquina Camiones
        hostname:dist20
        contraseña: Kcf25KUB


Una vez abiertas las terminales en las máquinas virtuales, ejecutar:

    $ cd SistDistribuidos/Lab1

dist17
    $ make logistica

dist18
    $ make financiero

dist19
    $ make cliente

dist20
    $ make camion 


Dirección Logístico:
    10.6.40.157