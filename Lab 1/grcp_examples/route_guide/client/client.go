/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a simple gRPC client that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It interacts with the route guide service whose definition can be found in routeguide/route_guide.proto.
package main

import (
	"context"
	"flag"
	"io"
	"log"
	"math/rand"
	"time"
	"os"
	
	"ioutil"
	"strconv"
	"encoding/csv"
	"encoding/json"

	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
	pb "google.golang.org/grpc/examples/route_guide/routeguide"
	//"google.golang.org/grpc/testdata"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name used to verify the hostname returned by the TLS handshake")
)

// printFeature gets the feature for the given point.
//func printFeature(client pb.RouteGuideClient, point *pb.Point) {
func printFeature(client pb.RouteGuideClient, pedido int) {
//	log.Printf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)
	log.Printf("Obtieniendo estado de pedido %d", pedido)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	estado, err := client.GetFeature(ctx, pedido)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	log.Println(estado)
}

// printFeatures lists all the features within the given bounding Rectangle.
func printFeatures(client pb.RouteGuideClient, rect *pb.Rectangle) {
	log.Printf("Looking for features within %v", rect)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ListFeatures(ctx, rect)
	if err != nil {
		log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
	}
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}
		log.Println(feature)
	}
}

// runRecordRoute sends a sequence of points to server and expects to get a RouteSummary from server.
func runRecordRoute(client pb.RouteGuideClient, data []byte) {
	
	/*
	// Create a random number of random points
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pointCount := int(r.Int31n(100)) + 2 // Traverse at least two points
	var points []*pb.Point
	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint(r))
	}
	log.Printf("Traversing %d points.", len(points))
	*/

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.RecordRoute(ctx)
	if err != nil {
		log.Fatalf("%v.RecordRoute(_) = _, %v", client, err)
	}

	if err := stream.Send(data); err != nil {
		log.Fatalf("%v.Send(%v) = %v", stream, data, err)
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Route summary: %v", reply)
}


// runRouteChat receives a sequence of route notes, while sending notes for various locations.
func runRouteChat(client pb.RouteGuideClient) {
	notes := []*pb.RouteNote{
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "First message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "Second message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "Third message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "Fourth message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "Fifth message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "Sixth message"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.RouteChat(ctx)
	if err != nil {
		log.Fatalf("%v.RouteChat(_) = _, %v", client, err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Printf("Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
		}
	}()
	for _, note := range notes {
		if err := stream.Send(note); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}

func randomPoint(r *rand.Rand) *pb.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &pb.Point{Latitude: lat, Longitude: long}
}


// loadFeatures loads features from a ~JSON~ CSV* file.
func (c *routeGuideClient) loadFeatures(filePath string) {
	var data []byte
	if filePath != "" {
		var err error
        //leer csv
		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to load default features: %v", err)
		}
        //csv a json
    } else{
        data = exampleData
    }
	if err := json.Unmarshal(data, &c.savedFeatures); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
}



type Retailparcel struct {
	ID       string `json:"id"`
	Producto string `json:"producto"`
	Valor    int `json:"valor"`
	Tienda   string `json:"tienda"`
	Destino  string `json:"destino"`
}

func leer(nombre string) *os.File {

	arch, err := os.Open(nombre)

	if err != nil {
		log.Fatalln("No se pudo abrir el archivo"+nombre+"para lectura", err) //chequear si está bien, o es con coma o qué
	}

	return arch
}


func main() {
	flag.Parse()
	var opts []grpc.DialOption

    /*
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
    */

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewRouteGuideClient(conn)



	csv_file, err := os.Open("retail.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer csv_file.Close()
	r := csv.NewReader(csv_file)
	records, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var rparcel Retailparcel
	var rparcels []Retailparcel
	for _, rec := range records {
		rparcel.ID = rec[0]
		rparcel.Producto = rec[1]
		rparcel.Valor,_ = strconv.Atoi(rec[2])
		rparcel.Tienda = rec[3]
		rparcel.Destino = rec[4]
		rparcels = append(rparcels, rparcel)
	}
	// Convert to JSON
	json_data, err := json.Marshal(rparcels)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//print json data
	//fmt.Println(string(json_data))
	fmt.Println("hola, todo bien ji")


	runRecordRoute(client, json_data)



/*
	// Looking for a valid feature
	printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})

	// Feature missing.
	printFeature(client, &pb.Point{Latitude: 0, Longitude: 0})

	// Looking for features between 40, -75 and 42, -73.
	printFeatures(client, &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	})
*/
	// RecordRoute
	runRecordRoute(client)

	// RouteChat
	runRouteChat(client)
}

// exampleData is a copy of testdata/route_guide_db.json. It's to avoid
// specifying file path with `go run`.
var exampleData = []byte(`[{
    "id": "SA1554KF",
    "producto": "polera",
    "valor":10,
    "tienda":tienda-A,
    "destino":casa-A,
},
{ "id": "TA1558KG",
    "producto": "pantalon",
    "valor":5,
    "tienda":tienda-B,
    "destino":casa-B,
},
{
    "id": "UA1559KH",
    "producto": "vaso",
    "valor":8,
    "tienda":tienda-C,
    "destino":casa-D,
},
{
    "id": "VA1551KI",
    "producto": "camion",
    "valor":2,
    "tienda":tienda-E,
    "destino":casa-F,
}]`)
