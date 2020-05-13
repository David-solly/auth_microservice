// Code generated by protoc-gen-go. DO NOT EDIT.
// source: token_service.proto

package v1

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type TokenState int32

const (
	TokenState_WATCH      TokenState = 0
	TokenState_LOGOUT     TokenState = 1
	TokenState_TRACEROUTE TokenState = 2
	TokenState_FREEZE     TokenState = 3
	TokenState_UNFREEZE   TokenState = 4
)

var TokenState_name = map[int32]string{
	0: "WATCH",
	1: "LOGOUT",
	2: "TRACEROUTE",
	3: "FREEZE",
	4: "UNFREEZE",
}

var TokenState_value = map[string]int32{
	"WATCH":      0,
	"LOGOUT":     1,
	"TRACEROUTE": 2,
	"FREEZE":     3,
	"UNFREEZE":   4,
}

func (x TokenState) String() string {
	return proto.EnumName(TokenState_name, int32(x))
}

func (TokenState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_84c67560d86c10bb, []int{0}
}

type TokenStatus int32

const (
	TokenStatus_INVALID    TokenStatus = 0
	TokenStatus_AUTHORIZED TokenStatus = 1
	TokenStatus_RESTRICTED TokenStatus = 2
	TokenStatus_EXPIRED    TokenStatus = 3
)

var TokenStatus_name = map[int32]string{
	0: "INVALID",
	1: "AUTHORIZED",
	2: "RESTRICTED",
	3: "EXPIRED",
}

var TokenStatus_value = map[string]int32{
	"INVALID":    0,
	"AUTHORIZED": 1,
	"RESTRICTED": 2,
	"EXPIRED":    3,
}

func (x TokenStatus) String() string {
	return proto.EnumName(TokenStatus_name, int32(x))
}

func (TokenStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_84c67560d86c10bb, []int{1}
}

type TokenAffectRequest struct {
	Token                string     `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	DesiredState         TokenState `protobuf:"varint,2,opt,name=desiredState,proto3,enum=v1.TokenState" json:"desiredState,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *TokenAffectRequest) Reset()         { *m = TokenAffectRequest{} }
func (m *TokenAffectRequest) String() string { return proto.CompactTextString(m) }
func (*TokenAffectRequest) ProtoMessage()    {}
func (*TokenAffectRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_84c67560d86c10bb, []int{0}
}

func (m *TokenAffectRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenAffectRequest.Unmarshal(m, b)
}
func (m *TokenAffectRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenAffectRequest.Marshal(b, m, deterministic)
}
func (m *TokenAffectRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenAffectRequest.Merge(m, src)
}
func (m *TokenAffectRequest) XXX_Size() int {
	return xxx_messageInfo_TokenAffectRequest.Size(m)
}
func (m *TokenAffectRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenAffectRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TokenAffectRequest proto.InternalMessageInfo

func (m *TokenAffectRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *TokenAffectRequest) GetDesiredState() TokenState {
	if m != nil {
		return m.DesiredState
	}
	return TokenState_WATCH
}

type TokenAffectResponse struct {
	// Types that are valid to be assigned to Affected:
	//	*TokenAffectResponse_EffectApplied
	//	*TokenAffectResponse_Error
	Affected             isTokenAffectResponse_Affected `protobuf_oneof:"affected"`
	XXX_NoUnkeyedLiteral struct{}                       `json:"-"`
	XXX_unrecognized     []byte                         `json:"-"`
	XXX_sizecache        int32                          `json:"-"`
}

func (m *TokenAffectResponse) Reset()         { *m = TokenAffectResponse{} }
func (m *TokenAffectResponse) String() string { return proto.CompactTextString(m) }
func (*TokenAffectResponse) ProtoMessage()    {}
func (*TokenAffectResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_84c67560d86c10bb, []int{1}
}

func (m *TokenAffectResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenAffectResponse.Unmarshal(m, b)
}
func (m *TokenAffectResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenAffectResponse.Marshal(b, m, deterministic)
}
func (m *TokenAffectResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenAffectResponse.Merge(m, src)
}
func (m *TokenAffectResponse) XXX_Size() int {
	return xxx_messageInfo_TokenAffectResponse.Size(m)
}
func (m *TokenAffectResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenAffectResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TokenAffectResponse proto.InternalMessageInfo

type isTokenAffectResponse_Affected interface {
	isTokenAffectResponse_Affected()
}

type TokenAffectResponse_EffectApplied struct {
	EffectApplied bool `protobuf:"varint,1,opt,name=effectApplied,proto3,oneof"`
}

type TokenAffectResponse_Error struct {
	Error *ServiceError `protobuf:"bytes,2,opt,name=error,proto3,oneof"`
}

func (*TokenAffectResponse_EffectApplied) isTokenAffectResponse_Affected() {}

func (*TokenAffectResponse_Error) isTokenAffectResponse_Affected() {}

func (m *TokenAffectResponse) GetAffected() isTokenAffectResponse_Affected {
	if m != nil {
		return m.Affected
	}
	return nil
}

func (m *TokenAffectResponse) GetEffectApplied() bool {
	if x, ok := m.GetAffected().(*TokenAffectResponse_EffectApplied); ok {
		return x.EffectApplied
	}
	return false
}

func (m *TokenAffectResponse) GetError() *ServiceError {
	if x, ok := m.GetAffected().(*TokenAffectResponse_Error); ok {
		return x.Error
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*TokenAffectResponse) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*TokenAffectResponse_EffectApplied)(nil),
		(*TokenAffectResponse_Error)(nil),
	}
}

type TokenVerifyRequest struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Service              string   `protobuf:"bytes,2,opt,name=service,proto3" json:"service,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TokenVerifyRequest) Reset()         { *m = TokenVerifyRequest{} }
func (m *TokenVerifyRequest) String() string { return proto.CompactTextString(m) }
func (*TokenVerifyRequest) ProtoMessage()    {}
func (*TokenVerifyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_84c67560d86c10bb, []int{2}
}

func (m *TokenVerifyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenVerifyRequest.Unmarshal(m, b)
}
func (m *TokenVerifyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenVerifyRequest.Marshal(b, m, deterministic)
}
func (m *TokenVerifyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenVerifyRequest.Merge(m, src)
}
func (m *TokenVerifyRequest) XXX_Size() int {
	return xxx_messageInfo_TokenVerifyRequest.Size(m)
}
func (m *TokenVerifyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenVerifyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TokenVerifyRequest proto.InternalMessageInfo

func (m *TokenVerifyRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *TokenVerifyRequest) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

type TokenVerifyResponse struct {
	Status               TokenStatus `protobuf:"varint,1,opt,name=status,proto3,enum=v1.TokenStatus" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *TokenVerifyResponse) Reset()         { *m = TokenVerifyResponse{} }
func (m *TokenVerifyResponse) String() string { return proto.CompactTextString(m) }
func (*TokenVerifyResponse) ProtoMessage()    {}
func (*TokenVerifyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_84c67560d86c10bb, []int{3}
}

func (m *TokenVerifyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenVerifyResponse.Unmarshal(m, b)
}
func (m *TokenVerifyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenVerifyResponse.Marshal(b, m, deterministic)
}
func (m *TokenVerifyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenVerifyResponse.Merge(m, src)
}
func (m *TokenVerifyResponse) XXX_Size() int {
	return xxx_messageInfo_TokenVerifyResponse.Size(m)
}
func (m *TokenVerifyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenVerifyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TokenVerifyResponse proto.InternalMessageInfo

func (m *TokenVerifyResponse) GetStatus() TokenStatus {
	if m != nil {
		return m.Status
	}
	return TokenStatus_INVALID
}

type TokenRenewRequest struct {
	RefreshToken         string   `protobuf:"bytes,1,opt,name=refreshToken,proto3" json:"refreshToken,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TokenRenewRequest) Reset()         { *m = TokenRenewRequest{} }
func (m *TokenRenewRequest) String() string { return proto.CompactTextString(m) }
func (*TokenRenewRequest) ProtoMessage()    {}
func (*TokenRenewRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_84c67560d86c10bb, []int{4}
}

func (m *TokenRenewRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenRenewRequest.Unmarshal(m, b)
}
func (m *TokenRenewRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenRenewRequest.Marshal(b, m, deterministic)
}
func (m *TokenRenewRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenRenewRequest.Merge(m, src)
}
func (m *TokenRenewRequest) XXX_Size() int {
	return xxx_messageInfo_TokenRenewRequest.Size(m)
}
func (m *TokenRenewRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenRenewRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TokenRenewRequest proto.InternalMessageInfo

func (m *TokenRenewRequest) GetRefreshToken() string {
	if m != nil {
		return m.RefreshToken
	}
	return ""
}

type TokenRequest struct {
	Claims               map[string]string `protobuf:"bytes,1,rep,name=claims,proto3" json:"claims,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *TokenRequest) Reset()         { *m = TokenRequest{} }
func (m *TokenRequest) String() string { return proto.CompactTextString(m) }
func (*TokenRequest) ProtoMessage()    {}
func (*TokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_84c67560d86c10bb, []int{5}
}

func (m *TokenRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenRequest.Unmarshal(m, b)
}
func (m *TokenRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenRequest.Marshal(b, m, deterministic)
}
func (m *TokenRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenRequest.Merge(m, src)
}
func (m *TokenRequest) XXX_Size() int {
	return xxx_messageInfo_TokenRequest.Size(m)
}
func (m *TokenRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TokenRequest proto.InternalMessageInfo

func (m *TokenRequest) GetClaims() map[string]string {
	if m != nil {
		return m.Claims
	}
	return nil
}

type TokenResponse struct {
	// Types that are valid to be assigned to Response:
	//	*TokenResponse_Tokens
	//	*TokenResponse_Error
	Response             isTokenResponse_Response `protobuf_oneof:"response"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *TokenResponse) Reset()         { *m = TokenResponse{} }
func (m *TokenResponse) String() string { return proto.CompactTextString(m) }
func (*TokenResponse) ProtoMessage()    {}
func (*TokenResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_84c67560d86c10bb, []int{6}
}

func (m *TokenResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenResponse.Unmarshal(m, b)
}
func (m *TokenResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenResponse.Marshal(b, m, deterministic)
}
func (m *TokenResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenResponse.Merge(m, src)
}
func (m *TokenResponse) XXX_Size() int {
	return xxx_messageInfo_TokenResponse.Size(m)
}
func (m *TokenResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TokenResponse proto.InternalMessageInfo

type isTokenResponse_Response interface {
	isTokenResponse_Response()
}

type TokenResponse_Tokens struct {
	Tokens *TokenPair `protobuf:"bytes,1,opt,name=tokens,proto3,oneof"`
}

type TokenResponse_Error struct {
	Error *ServiceError `protobuf:"bytes,2,opt,name=error,proto3,oneof"`
}

func (*TokenResponse_Tokens) isTokenResponse_Response() {}

func (*TokenResponse_Error) isTokenResponse_Response() {}

func (m *TokenResponse) GetResponse() isTokenResponse_Response {
	if m != nil {
		return m.Response
	}
	return nil
}

func (m *TokenResponse) GetTokens() *TokenPair {
	if x, ok := m.GetResponse().(*TokenResponse_Tokens); ok {
		return x.Tokens
	}
	return nil
}

func (m *TokenResponse) GetError() *ServiceError {
	if x, ok := m.GetResponse().(*TokenResponse_Error); ok {
		return x.Error
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*TokenResponse) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*TokenResponse_Tokens)(nil),
		(*TokenResponse_Error)(nil),
	}
}

type TokenPair struct {
	AuthToken            string   `protobuf:"bytes,1,opt,name=authToken,proto3" json:"authToken,omitempty"`
	RefreshToken         string   `protobuf:"bytes,2,opt,name=refreshToken,proto3" json:"refreshToken,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TokenPair) Reset()         { *m = TokenPair{} }
func (m *TokenPair) String() string { return proto.CompactTextString(m) }
func (*TokenPair) ProtoMessage()    {}
func (*TokenPair) Descriptor() ([]byte, []int) {
	return fileDescriptor_84c67560d86c10bb, []int{7}
}

func (m *TokenPair) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenPair.Unmarshal(m, b)
}
func (m *TokenPair) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenPair.Marshal(b, m, deterministic)
}
func (m *TokenPair) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenPair.Merge(m, src)
}
func (m *TokenPair) XXX_Size() int {
	return xxx_messageInfo_TokenPair.Size(m)
}
func (m *TokenPair) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenPair.DiscardUnknown(m)
}

var xxx_messageInfo_TokenPair proto.InternalMessageInfo

func (m *TokenPair) GetAuthToken() string {
	if m != nil {
		return m.AuthToken
	}
	return ""
}

func (m *TokenPair) GetRefreshToken() string {
	if m != nil {
		return m.RefreshToken
	}
	return ""
}

type ServiceError struct {
	Error                string   `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Code                 int32    `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServiceError) Reset()         { *m = ServiceError{} }
func (m *ServiceError) String() string { return proto.CompactTextString(m) }
func (*ServiceError) ProtoMessage()    {}
func (*ServiceError) Descriptor() ([]byte, []int) {
	return fileDescriptor_84c67560d86c10bb, []int{8}
}

func (m *ServiceError) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServiceError.Unmarshal(m, b)
}
func (m *ServiceError) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServiceError.Marshal(b, m, deterministic)
}
func (m *ServiceError) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServiceError.Merge(m, src)
}
func (m *ServiceError) XXX_Size() int {
	return xxx_messageInfo_ServiceError.Size(m)
}
func (m *ServiceError) XXX_DiscardUnknown() {
	xxx_messageInfo_ServiceError.DiscardUnknown(m)
}

var xxx_messageInfo_ServiceError proto.InternalMessageInfo

func (m *ServiceError) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *ServiceError) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func init() {
	proto.RegisterEnum("v1.TokenState", TokenState_name, TokenState_value)
	proto.RegisterEnum("v1.TokenStatus", TokenStatus_name, TokenStatus_value)
	proto.RegisterType((*TokenAffectRequest)(nil), "v1.TokenAffectRequest")
	proto.RegisterType((*TokenAffectResponse)(nil), "v1.TokenAffectResponse")
	proto.RegisterType((*TokenVerifyRequest)(nil), "v1.TokenVerifyRequest")
	proto.RegisterType((*TokenVerifyResponse)(nil), "v1.TokenVerifyResponse")
	proto.RegisterType((*TokenRenewRequest)(nil), "v1.TokenRenewRequest")
	proto.RegisterType((*TokenRequest)(nil), "v1.TokenRequest")
	proto.RegisterMapType((map[string]string)(nil), "v1.TokenRequest.ClaimsEntry")
	proto.RegisterType((*TokenResponse)(nil), "v1.TokenResponse")
	proto.RegisterType((*TokenPair)(nil), "v1.TokenPair")
	proto.RegisterType((*ServiceError)(nil), "v1.ServiceError")
}

func init() {
	proto.RegisterFile("token_service.proto", fileDescriptor_84c67560d86c10bb)
}

var fileDescriptor_84c67560d86c10bb = []byte{
	// 532 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0xcf, 0x6f, 0x12, 0x51,
	0x10, 0x86, 0xa5, 0x50, 0x98, 0x05, 0xdc, 0xbe, 0x7a, 0x20, 0x4d, 0x0f, 0xcd, 0x1e, 0x94, 0xf4,
	0x80, 0xe9, 0x6a, 0x62, 0xf5, 0xa0, 0x59, 0xe1, 0x09, 0x24, 0x6d, 0x69, 0x1e, 0x4b, 0x35, 0x3d,
	0x68, 0x56, 0x18, 0xe2, 0xa6, 0xb8, 0x8b, 0x6f, 0x1f, 0x18, 0x4e, 0xfe, 0xeb, 0xe6, 0xfd, 0x00,
	0x76, 0xd5, 0x98, 0x78, 0xdb, 0xf9, 0x66, 0xe6, 0xfb, 0xe6, 0x9b, 0x37, 0x0b, 0xc7, 0x22, 0x79,
	0xc0, 0xf8, 0x73, 0x8a, 0x7c, 0x1d, 0x4d, 0xb1, 0xb3, 0xe4, 0x89, 0x48, 0x88, 0xb5, 0xbe, 0x70,
	0x3f, 0x01, 0x09, 0x64, 0xca, 0x9f, 0xcf, 0x71, 0x2a, 0x18, 0x7e, 0x5f, 0x61, 0x2a, 0xc8, 0x63,
	0x28, 0xab, 0x86, 0x56, 0xf1, 0xac, 0xd8, 0xae, 0x31, 0x1d, 0x10, 0x0f, 0xea, 0x33, 0x4c, 0x23,
	0x8e, 0xb3, 0xb1, 0x08, 0x05, 0xb6, 0xac, 0xb3, 0x62, 0xbb, 0xe9, 0x35, 0x3b, 0xeb, 0x8b, 0x8e,
	0xe2, 0x50, 0x28, 0xcb, 0xd5, 0xb8, 0x29, 0x1c, 0xe7, 0xf8, 0xd3, 0x65, 0x12, 0xa7, 0x48, 0x9e,
	0x40, 0x03, 0x15, 0xe2, 0x2f, 0x97, 0x8b, 0x08, 0x67, 0x4a, 0xa8, 0x3a, 0x28, 0xb0, 0x3c, 0x4c,
	0xda, 0x50, 0x46, 0xce, 0x13, 0xae, 0xb4, 0x6c, 0xcf, 0x91, 0x5a, 0x63, 0x6d, 0x82, 0x4a, 0x7c,
	0x50, 0x60, 0xba, 0xe0, 0x1d, 0x40, 0x35, 0x54, 0xad, 0x38, 0x73, 0x7b, 0xc6, 0xd4, 0x1d, 0xf2,
	0x68, 0xbe, 0xf9, 0xb7, 0xa9, 0x16, 0x1c, 0x9a, 0xad, 0x28, 0x8d, 0x1a, 0xdb, 0x86, 0xee, 0x1b,
	0x33, 0xfa, 0x96, 0xc5, 0x8c, 0xfe, 0x14, 0x2a, 0xa9, 0x08, 0xc5, 0x2a, 0x55, 0x3c, 0x4d, 0xef,
	0x51, 0xce, 0xff, 0x2a, 0x65, 0x26, 0xed, 0xbe, 0x84, 0x23, 0x05, 0x33, 0x8c, 0xf1, 0xc7, 0x76,
	0x08, 0x17, 0xea, 0x1c, 0xe7, 0x1c, 0xd3, 0xaf, 0x41, 0x66, 0x96, 0x1c, 0xe6, 0xfe, 0x84, 0xba,
	0x69, 0xd4, 0x3d, 0x2f, 0xa0, 0x32, 0x5d, 0x84, 0xd1, 0x37, 0xa9, 0x58, 0x6a, 0xdb, 0xde, 0xe9,
	0x4e, 0xd1, 0x54, 0x74, 0xba, 0x2a, 0x4d, 0x63, 0xc1, 0x37, 0xcc, 0xd4, 0x9e, 0xbc, 0x02, 0x3b,
	0x03, 0x13, 0x07, 0x4a, 0x0f, 0xb8, 0x31, 0x7a, 0xf2, 0x53, 0xee, 0x63, 0x1d, 0x2e, 0x56, 0x5b,
	0xdf, 0x3a, 0x78, 0x6d, 0x5d, 0x16, 0xdd, 0x18, 0x1a, 0x86, 0x7e, 0xef, 0x59, 0x6d, 0x4b, 0x7b,
	0xb6, 0xbd, 0xc6, 0x6e, 0x82, 0xdb, 0x30, 0x92, 0x8f, 0x60, 0xd2, 0xff, 0xf7, 0x5e, 0xdc, 0xd0,
	0xbb, 0xd7, 0x50, 0xdb, 0x91, 0x91, 0x53, 0xa8, 0x85, 0x2b, 0x91, 0x5b, 0xcf, 0x1e, 0xf8, 0x63,
	0x7f, 0xd6, 0x5f, 0xf6, 0x77, 0x09, 0xf5, 0xac, 0xa6, 0x34, 0xaa, 0x87, 0x32, 0x0f, 0xaf, 0x02,
	0x42, 0xe0, 0x60, 0x9a, 0xcc, 0xb4, 0xfb, 0x32, 0x53, 0xdf, 0xe7, 0xd7, 0x00, 0xfb, 0x4b, 0x26,
	0x35, 0x28, 0x7f, 0xf0, 0x83, 0xee, 0xc0, 0x29, 0x10, 0x80, 0xca, 0xd5, 0xa8, 0x3f, 0x9a, 0x04,
	0x4e, 0x91, 0x34, 0x01, 0x02, 0xe6, 0x77, 0x29, 0x1b, 0x4d, 0x02, 0xea, 0x58, 0x32, 0xf7, 0x9e,
	0x51, 0x7a, 0x4f, 0x9d, 0x12, 0xa9, 0x43, 0x75, 0x72, 0x63, 0xa2, 0x83, 0xf3, 0x3e, 0xd8, 0x99,
	0xc3, 0x20, 0x36, 0x1c, 0x0e, 0x6f, 0xee, 0xfc, 0xab, 0x61, 0xcf, 0x29, 0x48, 0x16, 0x7f, 0x12,
	0x0c, 0x46, 0x6c, 0x78, 0x4f, 0x7b, 0x9a, 0x95, 0xd1, 0x71, 0xc0, 0x86, 0xdd, 0x80, 0xf6, 0x1c,
	0x4b, 0x16, 0xd3, 0x8f, 0xb7, 0x43, 0x46, 0x7b, 0x4e, 0xc9, 0x7b, 0x6b, 0x2e, 0xc2, 0xd8, 0x22,
	0xcf, 0xa0, 0xda, 0xc7, 0x18, 0xb9, 0x9c, 0xd2, 0xf9, 0xfd, 0x1a, 0x4e, 0x8e, 0x32, 0x88, 0xde,
	0xf0, 0x97, 0x8a, 0xfa, 0xe3, 0x9f, 0xff, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x9f, 0x05, 0x70, 0x53,
	0x08, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// TokenServiceClient is the client API for TokenService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TokenServiceClient interface {
	Generate(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*TokenResponse, error)
}

type tokenServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTokenServiceClient(cc grpc.ClientConnInterface) TokenServiceClient {
	return &tokenServiceClient{cc}
}

func (c *tokenServiceClient) Generate(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*TokenResponse, error) {
	out := new(TokenResponse)
	err := c.cc.Invoke(ctx, "/v1.TokenService/Generate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokenServiceServer is the server API for TokenService service.
type TokenServiceServer interface {
	Generate(context.Context, *TokenRequest) (*TokenResponse, error)
}

// UnimplementedTokenServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTokenServiceServer struct {
}

func (*UnimplementedTokenServiceServer) Generate(ctx context.Context, req *TokenRequest) (*TokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Generate not implemented")
}

func RegisterTokenServiceServer(s *grpc.Server, srv TokenServiceServer) {
	s.RegisterService(&_TokenService_serviceDesc, srv)
}

func _TokenService_Generate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).Generate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.TokenService/Generate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).Generate(ctx, req.(*TokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TokenService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.TokenService",
	HandlerType: (*TokenServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Generate",
			Handler:    _TokenService_Generate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "token_service.proto",
}
