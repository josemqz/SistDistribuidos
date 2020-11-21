package main
//CLIENT

import (
	"log"
	"time"
	
	"github.com/josemqz/SistDistribuidos/test/testp"
	"google.golang.org/grpc"
)

func main(){

	conn, err := grpc.Dial(dirNN, grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a Testp")
	defer conn.Close()

	client := logis.NewTestpServiceClient(conn)
	log.Println("Conexión realizada\n")

	ctx, cancel := context.WithTimeout(context.Background(), 1000 * time.Second)
	defer cancel()

	stat1, err := SubirArchivo(client, "./srcBooks/Alicia_a_traves_del_espejo-Carroll_Lewis.pdf", ctx)

}

func SubirArchivo(client book.BookServiceClient, ctx context.Context, f string) (stats Stats, err error) {

	// Get a file handle for the file we 
	// want to upload
	file, err = os.Open(f)
	if (err != nil) {
		log.Fatalf("%s: %s\n", "no se pudo abrir archivo", err)
	}

	// Open a stream-based connection with the 
	// gRPC server
	stream, err := client.RecibirBytes(ctx)
	if (err != nil) {
		log.Fatalf("%s: %s\n", "error llamando función RecibirBytes", err)
	}

	//file size
	fileInfo, _ := file.Stat()
	var fileSize int32 = fileInfo.Size()

	// Start timing the execution
	stats.StartedAt = time.Now()

	// Allocate a buffer with `chunkSize` as the capacity
	// and length (making a 0 array of the size of `chunkSize`)
	buf = make([]byte, fileSize)
	for writing {
		// put as many bytes as `chunkSize` into the
		// buf array.
		n, err = file.Read(buf)

		// ... if `eof` --> `writing=false`...

		stream.Send(&testp.MSG{
				// because we might've read less than
				// `chunkSize` we want to only send up to
				// `n` (amount of bytes read).
				// note: slicing (`:n`) won't copy the 
				// underlying data, so this as fast as taking
				// a "pointer" to the underlying storage.
				arch: buf[:n]})
	}

	// keep track of the end time so that we can take the elapsed
	// time later
	stats.FinishedAt = time.Now()

	// close
	status, err = stream.CloseAndRecv()
}