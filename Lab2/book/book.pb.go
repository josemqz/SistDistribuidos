// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.6.1
// source: book.proto

package book

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Chunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Algoritmo     string `protobuf:"bytes,1,opt,name=algoritmo,proto3" json:"algoritmo,omitempty"`
	NombreLibro   string `protobuf:"bytes,2,opt,name=nombreLibro,proto3" json:"nombreLibro,omitempty"`
	NumChunk      int32  `protobuf:"varint,3,opt,name=numChunk,proto3" json:"numChunk,omitempty"`
	Contenido     []byte `protobuf:"bytes,5,opt,name=contenido,proto3" json:"contenido,omitempty"`
	NombreArchivo string `protobuf:"bytes,6,opt,name=nombreArchivo,proto3" json:"nombreArchivo,omitempty"`
}

func (x *Chunk) Reset() {
	*x = Chunk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_book_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Chunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chunk) ProtoMessage() {}

func (x *Chunk) ProtoReflect() protoreflect.Message {
	mi := &file_book_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Chunk.ProtoReflect.Descriptor instead.
func (*Chunk) Descriptor() ([]byte, []int) {
	return file_book_proto_rawDescGZIP(), []int{0}
}

func (x *Chunk) GetAlgoritmo() string {
	if x != nil {
		return x.Algoritmo
	}
	return ""
}

func (x *Chunk) GetNombreLibro() string {
	if x != nil {
		return x.NombreLibro
	}
	return ""
}

func (x *Chunk) GetNumChunk() int32 {
	if x != nil {
		return x.NumChunk
	}
	return 0
}

func (x *Chunk) GetContenido() []byte {
	if x != nil {
		return x.Contenido
	}
	return nil
}

func (x *Chunk) GetNombreArchivo() string {
	if x != nil {
		return x.NombreArchivo
	}
	return ""
}

//para enviar direcciones de chunks a ClienteDownloader
type ChunksInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NombreLibro string `protobuf:"bytes,1,opt,name=nombreLibro,proto3" json:"nombreLibro,omitempty"`
	CantChunks  int32  `protobuf:"varint,2,opt,name=cantChunks,proto3" json:"cantChunks,omitempty"`
	Info        string `protobuf:"bytes,3,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *ChunksInfo) Reset() {
	*x = ChunksInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_book_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChunksInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChunksInfo) ProtoMessage() {}

func (x *ChunksInfo) ProtoReflect() protoreflect.Message {
	mi := &file_book_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChunksInfo.ProtoReflect.Descriptor instead.
func (*ChunksInfo) Descriptor() ([]byte, []int) {
	return file_book_proto_rawDescGZIP(), []int{1}
}

func (x *ChunksInfo) GetNombreLibro() string {
	if x != nil {
		return x.NombreLibro
	}
	return ""
}

func (x *ChunksInfo) GetCantChunks() int32 {
	if x != nil {
		return x.CantChunks
	}
	return 0
}

func (x *ChunksInfo) GetInfo() string {
	if x != nil {
		return x.Info
	}
	return ""
}

type ExMutua struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tiempo  string `protobuf:"bytes,1,opt,name=tiempo,proto3" json:"tiempo,omitempty"`   //timestamp
	Entidad string `protobuf:"bytes,2,opt,name=entidad,proto3" json:"entidad,omitempty"` //datanode A B o C
}

func (x *ExMutua) Reset() {
	*x = ExMutua{}
	if protoimpl.UnsafeEnabled {
		mi := &file_book_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExMutua) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExMutua) ProtoMessage() {}

func (x *ExMutua) ProtoReflect() protoreflect.Message {
	mi := &file_book_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExMutua.ProtoReflect.Descriptor instead.
func (*ExMutua) Descriptor() ([]byte, []int) {
	return file_book_proto_rawDescGZIP(), []int{2}
}

func (x *ExMutua) GetTiempo() string {
	if x != nil {
		return x.Tiempo
	}
	return ""
}

func (x *ExMutua) GetEntidad() string {
	if x != nil {
		return x.Entidad
	}
	return ""
}

type ListaLibros struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lista string `protobuf:"bytes,1,opt,name=lista,proto3" json:"lista,omitempty"`
}

func (x *ListaLibros) Reset() {
	*x = ListaLibros{}
	if protoimpl.UnsafeEnabled {
		mi := &file_book_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListaLibros) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListaLibros) ProtoMessage() {}

func (x *ListaLibros) ProtoReflect() protoreflect.Message {
	mi := &file_book_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListaLibros.ProtoReflect.Descriptor instead.
func (*ListaLibros) Descriptor() ([]byte, []int) {
	return file_book_proto_rawDescGZIP(), []int{3}
}

func (x *ListaLibros) GetLista() string {
	if x != nil {
		return x.Lista
	}
	return ""
}

type PropuestaLibro struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NombreLibro string `protobuf:"bytes,1,opt,name=nombreLibro,proto3" json:"nombreLibro,omitempty"`
	Propuesta   string `protobuf:"bytes,2,opt,name=propuesta,proto3" json:"propuesta,omitempty"`
	CantChunks  int32  `protobuf:"varint,3,opt,name=cantChunks,proto3" json:"cantChunks,omitempty"`
	//si la propuesta considera usar datanode a, b y/o c
	DatanodeA bool `protobuf:"varint,4,opt,name=datanodeA,proto3" json:"datanodeA,omitempty"`
	DatanodeB bool `protobuf:"varint,5,opt,name=datanodeB,proto3" json:"datanodeB,omitempty"`
	DatanodeC bool `protobuf:"varint,6,opt,name=datanodeC,proto3" json:"datanodeC,omitempty"`
}

func (x *PropuestaLibro) Reset() {
	*x = PropuestaLibro{}
	if protoimpl.UnsafeEnabled {
		mi := &file_book_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PropuestaLibro) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PropuestaLibro) ProtoMessage() {}

func (x *PropuestaLibro) ProtoReflect() protoreflect.Message {
	mi := &file_book_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PropuestaLibro.ProtoReflect.Descriptor instead.
func (*PropuestaLibro) Descriptor() ([]byte, []int) {
	return file_book_proto_rawDescGZIP(), []int{4}
}

func (x *PropuestaLibro) GetNombreLibro() string {
	if x != nil {
		return x.NombreLibro
	}
	return ""
}

func (x *PropuestaLibro) GetPropuesta() string {
	if x != nil {
		return x.Propuesta
	}
	return ""
}

func (x *PropuestaLibro) GetCantChunks() int32 {
	if x != nil {
		return x.CantChunks
	}
	return 0
}

func (x *PropuestaLibro) GetDatanodeA() bool {
	if x != nil {
		return x.DatanodeA
	}
	return false
}

func (x *PropuestaLibro) GetDatanodeB() bool {
	if x != nil {
		return x.DatanodeB
	}
	return false
}

func (x *PropuestaLibro) GetDatanodeC() bool {
	if x != nil {
		return x.DatanodeC
	}
	return false
}

type RespuestaP struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Respuesta bool   `protobuf:"varint,1,opt,name=respuesta,proto3" json:"respuesta,omitempty"`
	DNcaido   string `protobuf:"bytes,2,opt,name=DNcaido,proto3" json:"DNcaido,omitempty"`
}

func (x *RespuestaP) Reset() {
	*x = RespuestaP{}
	if protoimpl.UnsafeEnabled {
		mi := &file_book_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RespuestaP) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RespuestaP) ProtoMessage() {}

func (x *RespuestaP) ProtoReflect() protoreflect.Message {
	mi := &file_book_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RespuestaP.ProtoReflect.Descriptor instead.
func (*RespuestaP) Descriptor() ([]byte, []int) {
	return file_book_proto_rawDescGZIP(), []int{5}
}

func (x *RespuestaP) GetRespuesta() bool {
	if x != nil {
		return x.Respuesta
	}
	return false
}

func (x *RespuestaP) GetDNcaido() string {
	if x != nil {
		return x.DNcaido
	}
	return ""
}

type ACK struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok string `protobuf:"bytes,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *ACK) Reset() {
	*x = ACK{}
	if protoimpl.UnsafeEnabled {
		mi := &file_book_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ACK) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ACK) ProtoMessage() {}

func (x *ACK) ProtoReflect() protoreflect.Message {
	mi := &file_book_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ACK.ProtoReflect.Descriptor instead.
func (*ACK) Descriptor() ([]byte, []int) {
	return file_book_proto_rawDescGZIP(), []int{6}
}

func (x *ACK) GetOk() string {
	if x != nil {
		return x.Ok
	}
	return ""
}

var File_book_proto protoreflect.FileDescriptor

var file_book_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x62, 0x6f, 0x6f, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x62, 0x6f,
	0x6f, 0x6b, 0x22, 0xa7, 0x01, 0x0a, 0x05, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x1c, 0x0a, 0x09,
	0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x6d, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x6d, 0x6f, 0x12, 0x20, 0x0a, 0x0b, 0x6e, 0x6f,
	0x6d, 0x62, 0x72, 0x65, 0x4c, 0x69, 0x62, 0x72, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x6e, 0x6f, 0x6d, 0x62, 0x72, 0x65, 0x4c, 0x69, 0x62, 0x72, 0x6f, 0x12, 0x1a, 0x0a, 0x08,
	0x6e, 0x75, 0x6d, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x6e, 0x75, 0x6d, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x69, 0x64, 0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x69, 0x64, 0x6f, 0x12, 0x24, 0x0a, 0x0d, 0x6e, 0x6f, 0x6d, 0x62, 0x72, 0x65,
	0x41, 0x72, 0x63, 0x68, 0x69, 0x76, 0x6f, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e,
	0x6f, 0x6d, 0x62, 0x72, 0x65, 0x41, 0x72, 0x63, 0x68, 0x69, 0x76, 0x6f, 0x22, 0x62, 0x0a, 0x0a,
	0x43, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x20, 0x0a, 0x0b, 0x6e, 0x6f,
	0x6d, 0x62, 0x72, 0x65, 0x4c, 0x69, 0x62, 0x72, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x6e, 0x6f, 0x6d, 0x62, 0x72, 0x65, 0x4c, 0x69, 0x62, 0x72, 0x6f, 0x12, 0x1e, 0x0a, 0x0a,
	0x63, 0x61, 0x6e, 0x74, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0a, 0x63, 0x61, 0x6e, 0x74, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x69, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f,
	0x22, 0x3b, 0x0a, 0x07, 0x45, 0x78, 0x4d, 0x75, 0x74, 0x75, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x74,
	0x69, 0x65, 0x6d, 0x70, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x69, 0x65,
	0x6d, 0x70, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x69, 0x64, 0x61, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x69, 0x64, 0x61, 0x64, 0x22, 0x23, 0x0a,
	0x0b, 0x4c, 0x69, 0x73, 0x74, 0x61, 0x4c, 0x69, 0x62, 0x72, 0x6f, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x6c, 0x69, 0x73, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x69, 0x73,
	0x74, 0x61, 0x22, 0xca, 0x01, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x70, 0x75, 0x65, 0x73, 0x74, 0x61,
	0x4c, 0x69, 0x62, 0x72, 0x6f, 0x12, 0x20, 0x0a, 0x0b, 0x6e, 0x6f, 0x6d, 0x62, 0x72, 0x65, 0x4c,
	0x69, 0x62, 0x72, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6e, 0x6f, 0x6d, 0x62,
	0x72, 0x65, 0x4c, 0x69, 0x62, 0x72, 0x6f, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x70, 0x75,
	0x65, 0x73, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x70,
	0x75, 0x65, 0x73, 0x74, 0x61, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x61, 0x6e, 0x74, 0x43, 0x68, 0x75,
	0x6e, 0x6b, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x63, 0x61, 0x6e, 0x74, 0x43,
	0x68, 0x75, 0x6e, 0x6b, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x61, 0x74, 0x61, 0x6e, 0x6f, 0x64,
	0x65, 0x41, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x64, 0x61, 0x74, 0x61, 0x6e, 0x6f,
	0x64, 0x65, 0x41, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x61, 0x74, 0x61, 0x6e, 0x6f, 0x64, 0x65, 0x42,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x64, 0x61, 0x74, 0x61, 0x6e, 0x6f, 0x64, 0x65,
	0x42, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x61, 0x74, 0x61, 0x6e, 0x6f, 0x64, 0x65, 0x43, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x64, 0x61, 0x74, 0x61, 0x6e, 0x6f, 0x64, 0x65, 0x43, 0x22,
	0x44, 0x0a, 0x0a, 0x52, 0x65, 0x73, 0x70, 0x75, 0x65, 0x73, 0x74, 0x61, 0x50, 0x12, 0x1c, 0x0a,
	0x09, 0x72, 0x65, 0x73, 0x70, 0x75, 0x65, 0x73, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x09, 0x72, 0x65, 0x73, 0x70, 0x75, 0x65, 0x73, 0x74, 0x61, 0x12, 0x18, 0x0a, 0x07, 0x44,
	0x4e, 0x63, 0x61, 0x69, 0x64, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x44, 0x4e,
	0x63, 0x61, 0x69, 0x64, 0x6f, 0x22, 0x15, 0x0a, 0x03, 0x41, 0x43, 0x4b, 0x12, 0x0e, 0x0a, 0x02,
	0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x6f, 0x6b, 0x32, 0xb7, 0x03, 0x0a,
	0x0b, 0x42, 0x6f, 0x6f, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2b, 0x0a, 0x0d,
	0x45, 0x6e, 0x76, 0x69, 0x61, 0x72, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x44, 0x4e, 0x12, 0x0b, 0x2e,
	0x62, 0x6f, 0x6f, 0x6b, 0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x1a, 0x0b, 0x2e, 0x62, 0x6f, 0x6f,
	0x6b, 0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x22, 0x00, 0x12, 0x2b, 0x0a, 0x0f, 0x52, 0x65, 0x63,
	0x69, 0x62, 0x69, 0x72, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x44, 0x4e, 0x12, 0x0b, 0x2e, 0x62,
	0x6f, 0x6f, 0x6b, 0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x1a, 0x09, 0x2e, 0x62, 0x6f, 0x6f, 0x6b,
	0x2e, 0x41, 0x43, 0x4b, 0x22, 0x00, 0x12, 0x3c, 0x0a, 0x10, 0x52, 0x65, 0x63, 0x69, 0x62, 0x69,
	0x72, 0x50, 0x72, 0x6f, 0x70, 0x75, 0x65, 0x73, 0x74, 0x61, 0x12, 0x14, 0x2e, 0x62, 0x6f, 0x6f,
	0x6b, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x75, 0x65, 0x73, 0x74, 0x61, 0x4c, 0x69, 0x62, 0x72, 0x6f,
	0x1a, 0x10, 0x2e, 0x62, 0x6f, 0x6f, 0x6b, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x75, 0x65, 0x73, 0x74,
	0x61, 0x50, 0x22, 0x00, 0x12, 0x2b, 0x0a, 0x0d, 0x52, 0x65, 0x63, 0x69, 0x62, 0x69, 0x72, 0x43,
	0x68, 0x75, 0x6e, 0x6b, 0x73, 0x12, 0x0b, 0x2e, 0x62, 0x6f, 0x6f, 0x6b, 0x2e, 0x43, 0x68, 0x75,
	0x6e, 0x6b, 0x1a, 0x09, 0x2e, 0x62, 0x6f, 0x6f, 0x6b, 0x2e, 0x41, 0x43, 0x4b, 0x22, 0x00, 0x28,
	0x01, 0x12, 0x33, 0x0a, 0x0e, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x69, 0x72, 0x4c, 0x6f, 0x67,
	0x44, 0x65, 0x73, 0x12, 0x14, 0x2e, 0x62, 0x6f, 0x6f, 0x6b, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x75,
	0x65, 0x73, 0x74, 0x61, 0x4c, 0x69, 0x62, 0x72, 0x6f, 0x1a, 0x09, 0x2e, 0x62, 0x6f, 0x6f, 0x6b,
	0x2e, 0x41, 0x43, 0x4b, 0x22, 0x00, 0x12, 0x43, 0x0a, 0x13, 0x72, 0x65, 0x63, 0x69, 0x62, 0x69,
	0x72, 0x50, 0x72, 0x6f, 0x70, 0x44, 0x61, 0x74, 0x61, 0x6e, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x2e,
	0x62, 0x6f, 0x6f, 0x6b, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x75, 0x65, 0x73, 0x74, 0x61, 0x4c, 0x69,
	0x62, 0x72, 0x6f, 0x1a, 0x14, 0x2e, 0x62, 0x6f, 0x6f, 0x6b, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x75,
	0x65, 0x73, 0x74, 0x61, 0x4c, 0x69, 0x62, 0x72, 0x6f, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x0c, 0x43,
	0x68, 0x75, 0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x6f, 0x67, 0x12, 0x10, 0x2e, 0x62, 0x6f,
	0x6f, 0x6b, 0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x10, 0x2e,
	0x62, 0x6f, 0x6f, 0x6b, 0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x22,
	0x00, 0x12, 0x33, 0x0a, 0x11, 0x45, 0x6e, 0x76, 0x69, 0x61, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x61,
	0x4c, 0x69, 0x62, 0x72, 0x6f, 0x73, 0x12, 0x09, 0x2e, 0x62, 0x6f, 0x6f, 0x6b, 0x2e, 0x41, 0x43,
	0x4b, 0x1a, 0x11, 0x2e, 0x62, 0x6f, 0x6f, 0x6b, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x61, 0x4c, 0x69,
	0x62, 0x72, 0x6f, 0x73, 0x22, 0x00, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x6f, 0x73, 0x65, 0x6d, 0x71, 0x7a, 0x2f, 0x53, 0x69, 0x73,
	0x74, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x69, 0x64, 0x6f, 0x73, 0x2f, 0x4c, 0x61,
	0x62, 0x32, 0x2f, 0x62, 0x6f, 0x6f, 0x6b, 0x3b, 0x62, 0x6f, 0x6f, 0x6b, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_book_proto_rawDescOnce sync.Once
	file_book_proto_rawDescData = file_book_proto_rawDesc
)

func file_book_proto_rawDescGZIP() []byte {
	file_book_proto_rawDescOnce.Do(func() {
		file_book_proto_rawDescData = protoimpl.X.CompressGZIP(file_book_proto_rawDescData)
	})
	return file_book_proto_rawDescData
}

var file_book_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_book_proto_goTypes = []interface{}{
	(*Chunk)(nil),          // 0: book.Chunk
	(*ChunksInfo)(nil),     // 1: book.ChunksInfo
	(*ExMutua)(nil),        // 2: book.ExMutua
	(*ListaLibros)(nil),    // 3: book.ListaLibros
	(*PropuestaLibro)(nil), // 4: book.PropuestaLibro
	(*RespuestaP)(nil),     // 5: book.RespuestaP
	(*ACK)(nil),            // 6: book.ACK
}
var file_book_proto_depIdxs = []int32{
	0, // 0: book.BookService.EnviarChunkDN:input_type -> book.Chunk
	0, // 1: book.BookService.RecibirChunksDN:input_type -> book.Chunk
	4, // 2: book.BookService.RecibirPropuesta:input_type -> book.PropuestaLibro
	0, // 3: book.BookService.RecibirChunks:input_type -> book.Chunk
	4, // 4: book.BookService.escribirLogDes:input_type -> book.PropuestaLibro
	4, // 5: book.BookService.recibirPropDatanode:input_type -> book.PropuestaLibro
	1, // 6: book.BookService.ChunkInfoLog:input_type -> book.ChunksInfo
	6, // 7: book.BookService.EnviarListaLibros:input_type -> book.ACK
	0, // 8: book.BookService.EnviarChunkDN:output_type -> book.Chunk
	6, // 9: book.BookService.RecibirChunksDN:output_type -> book.ACK
	5, // 10: book.BookService.RecibirPropuesta:output_type -> book.RespuestaP
	6, // 11: book.BookService.RecibirChunks:output_type -> book.ACK
	6, // 12: book.BookService.escribirLogDes:output_type -> book.ACK
	4, // 13: book.BookService.recibirPropDatanode:output_type -> book.PropuestaLibro
	1, // 14: book.BookService.ChunkInfoLog:output_type -> book.ChunksInfo
	3, // 15: book.BookService.EnviarListaLibros:output_type -> book.ListaLibros
	8, // [8:16] is the sub-list for method output_type
	0, // [0:8] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_book_proto_init() }
func file_book_proto_init() {
	if File_book_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_book_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Chunk); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_book_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChunksInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_book_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExMutua); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_book_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListaLibros); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_book_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PropuestaLibro); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_book_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RespuestaP); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_book_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ACK); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_book_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_book_proto_goTypes,
		DependencyIndexes: file_book_proto_depIdxs,
		MessageInfos:      file_book_proto_msgTypes,
	}.Build()
	File_book_proto = out.File
	file_book_proto_rawDesc = nil
	file_book_proto_goTypes = nil
	file_book_proto_depIdxs = nil
}