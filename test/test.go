package main
//CLIENT

import (
	"log"
	"time"
	"context"
	"io"
	//"stats"
	"os"
	
	"github.com/josemqz/SistDistribuidos/test/testp"
	"google.golang.org/grpc"
	//"golang.org/x/perf/internal/stats"
)

func failOnError(err error, msg string) {
	if (err != nil) {
		log.Fatalf("%s: %s\n", msg, err)
	}
}

func main(){

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a Testp")
	defer conn.Close()

	client := testp.NewTestpServiceClient(conn)
	log.Println("Conexión realizada\n")

	ctx, cancel := context.WithTimeout(context.Background(), 1000 * time.Second)
	defer cancel()

	err = SubirArchivo(client, ctx, "./srcBooks/Alicia_a_traves_del_espejo-Carroll_Lewis.pdf")

	failOnError(err,"ble")

}

func SubirArchivo(client testp.TestpServiceClient, ctx context.Context, f string) (err error) {

	// Get a file handle for the file we 
	// want to upload
	file, err := os.Open(f)
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
	var fileSize int32 = int32(fileInfo.Size())

	// Start timing the execution
	//stats.StartedAt = time.Now()

	// Allocate a buffer with `chunkSize` as the capacity
	// and length (making a 0 array of the size of `chunkSize`)
	buf := make([]byte, fileSize)
	for {
		// put as many bytes as `chunkSize` into the
		// buf array.
		n, err := file.Read(buf)

		if (err != nil){
			if err == io.EOF{
				break
			}
			failOnError(err, "aaaaah")
			return err
		}

		// ... if `eof` --> `writing=false`...

		stream.Send(&testp.MSG{
				// because we might've read less than
				// `chunkSize` we want to only send up to
				// `n` (amount of bytes read).
				// note: slicing (`:n`) won't copy the 
				// underlying data, so this as fast as taking
				// a "pointer" to the underlying storage.
				Arch: buf[:n]})
	}

	// keep track of the end time so that we can take the elapsed
	// time later
	//stats.FinishedAt = time.Now()

	// close
	_, err = stream.CloseAndRecv()

	return nil
}