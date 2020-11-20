//SERVER
package main

import (
	"log"
)

func main(){

	
}


// Upload implements the Upload method of the GuploadService
// interface which is responsible for receiving a stream of
// chunks that form a complete file.
func (s *ServerGRPC) RecibirBytes(stream testp.GuploadService_UploadServer) (err error) {
	// while there are messages coming
	for {
		_, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				goto END
			}

			err = errors.Wrapf(err,
				"failed unexpectadely while reading chunks from stream")
			return
		}
	}

	END:
	// once the transmission finished, send the
	// confirmation if nothign went wrong
	err = stream.SendAndClose(&messaging.UploadStatus{
	Message: "Upload received with success",
	Code:    messaging.UploadStatusCode_Ok,
	})
	// ...

	return
}