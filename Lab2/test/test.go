//CLIENT
package main

import (
	"log"
	"time"
)

func main(){

	conn, err := grpc.Dial(dirNN, grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a NameNode")
	defer conn.Close()

	camion := logis.NewLogisServiceClient(conn)
	log.Println("Conexión realizada\n")

}

func (c *ClientGRPC) UploadFile(ctx context.Context, f string) (stats Stats, err error) {

	// Get a file handle for the file we 
	// want to upload
	file, err = os.Open(f)

	// Open a stream-based connection with the 
	// gRPC server
	stream, err := c.client.Upload(ctx)

	// Start timing the execution
	stats.StartedAt = time.Now()

	// Allocate a buffer with `chunkSize` as the capacity
	// and length (making a 0 array of the size of `chunkSize`)
	buf = make([]byte, c.chunkSize)
	for writing {
		// put as many bytes as `chunkSize` into the
		// buf array.
		n, err = file.Read(buf)

		// ... if `eof` --> `writing=false`...

		stream.Send(&messaging.Chunk{
				// because we might've read less than
				// `chunkSize` we want to only send up to
				// `n` (amount of bytes read).
				// note: slicing (`:n`) won't copy the 
				// underlying data, so this as fast as taking
				// a "pointer" to the underlying storage.
				Content: buf[:n],
		})
	}

	// keep track of the end time so that we can take the elapsed
	// time later
	stats.FinishedAt = time.Now()

	// close
	status, err = stream.CloseAndRecv()
}