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
	DesiredState         TokenState `protobuf:"varint,2,opt,name=desired_state,json=desiredState,proto3,enum=v1.TokenState" json:"desired_state,omitempty"`
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
	EffectApplied        bool          `protobuf:"varint,1,opt,name=effect_applied,json=effectApplied,proto3" json:"effect_applied,omitempty"`
	Error                *ServiceError `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
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

func (m *TokenAffectResponse) GetEffectApplied() bool {
	if m != nil {
		return m.EffectApplied
	}
	return false
}

func (m *TokenAffectResponse) GetError() *ServiceError {
	if m != nil {
		return m.Error
	}
	return nil
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
	Status               TokenStatus       `protobuf:"varint,1,opt,name=status,proto3,enum=v1.TokenStatus" json:"status,omitempty"`
	UserId               string            `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Claims               map[string]string `protobuf:"bytes,3,rep,name=claims,proto3" json:"claims,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Error                *ServiceError     `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
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

func (m *TokenVerifyResponse) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *TokenVerifyResponse) GetClaims() map[string]string {
	if m != nil {
		return m.Claims
	}
	return nil
}

func (m *TokenVerifyResponse) GetError() *ServiceError {
	if m != nil {
		return m.Error
	}
	return nil
}

type TokenRenewRequest struct {
	RefreshToken         string   `protobuf:"bytes,1,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
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
	Tokens               *TokenPair    `protobuf:"bytes,1,opt,name=tokens,proto3" json:"tokens,omitempty"`
	Error                *ServiceError `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
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

func (m *TokenResponse) GetTokens() *TokenPair {
	if m != nil {
		return m.Tokens
	}
	return nil
}

func (m *TokenResponse) GetError() *ServiceError {
	if m != nil {
		return m.Error
	}
	return nil
}

type TokenPair struct {
	AuthToken            string   `protobuf:"bytes,1,opt,name=auth_token,json=authToken,proto3" json:"auth_token,omitempty"`
	RefreshToken         string   `protobuf:"bytes,2,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
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
	proto.RegisterMapType((map[string]string)(nil), "v1.TokenVerifyResponse.ClaimsEntry")
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
	// 590 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0xdf, 0x6f, 0xd2, 0x5e,
	0x14, 0x5f, 0xcb, 0xe8, 0xc6, 0x29, 0xf0, 0xed, 0xee, 0xbe, 0x71, 0x64, 0xd1, 0x84, 0x74, 0x51,
	0xc9, 0x1e, 0x30, 0xeb, 0x7c, 0x40, 0x4d, 0x4c, 0x1a, 0xb8, 0xb2, 0x26, 0x73, 0x2c, 0x97, 0x32,
	0xcd, 0x1e, 0x6c, 0x2a, 0xbd, 0xc4, 0x66, 0x48, 0xf1, 0xb6, 0xc5, 0xf0, 0xe4, 0x9f, 0xed, 0xab,
	0xb9, 0x3f, 0xe8, 0xc0, 0xcd, 0xc5, 0xc4, 0xb7, 0x9e, 0x5f, 0x9f, 0x73, 0x3e, 0x9f, 0x7b, 0x4e,
	0x61, 0x3f, 0x4b, 0x6e, 0xe8, 0x2c, 0x48, 0x29, 0x5b, 0xc4, 0x63, 0xda, 0x9e, 0xb3, 0x24, 0x4b,
	0x90, 0xbe, 0x38, 0xb1, 0x03, 0x40, 0x3e, 0x0f, 0xb9, 0x93, 0x09, 0x1d, 0x67, 0x84, 0x7e, 0xcb,
	0x69, 0x9a, 0xa1, 0xff, 0xa1, 0x2c, 0x0a, 0x1a, 0x5a, 0x53, 0x6b, 0x55, 0x88, 0x34, 0xd0, 0x29,
	0xd4, 0x22, 0x9a, 0xc6, 0x8c, 0x46, 0x41, 0x9a, 0x85, 0x19, 0x6d, 0xe8, 0x4d, 0xad, 0x55, 0x77,
	0xea, 0xed, 0xc5, 0x49, 0x5b, 0x80, 0x0c, 0xb9, 0x97, 0x54, 0x55, 0x92, 0xb0, 0xec, 0x08, 0xf6,
	0x37, 0x1a, 0xa4, 0xf3, 0x64, 0x96, 0x52, 0xf4, 0x14, 0xea, 0x54, 0x78, 0x82, 0x70, 0x3e, 0x9f,
	0xc6, 0x34, 0x12, 0xad, 0x76, 0x49, 0x4d, 0x7a, 0x5d, 0xe9, 0x44, 0xcf, 0xa0, 0x4c, 0x19, 0x4b,
	0x98, 0x68, 0x65, 0x3a, 0x16, 0x6f, 0x35, 0x94, 0x24, 0x30, 0xf7, 0x13, 0x19, 0xb6, 0x7b, 0x8a,
	0xc6, 0x15, 0x65, 0xf1, 0x64, 0xf9, 0x30, 0x8d, 0x06, 0xec, 0x28, 0x1d, 0x04, 0x6a, 0x85, 0xac,
	0x4c, 0xfb, 0xa7, 0xa6, 0x86, 0x5d, 0xc1, 0xa8, 0x61, 0x9f, 0x83, 0xc1, 0x09, 0xe7, 0xa9, 0x00,
	0xaa, 0x3b, 0xff, 0x6d, 0x30, 0xce, 0x53, 0xa2, 0xc2, 0xe8, 0x00, 0x76, 0xf2, 0x94, 0xb2, 0x20,
	0x8e, 0x14, 0xb4, 0xc1, 0x4d, 0x2f, 0x42, 0x6f, 0xc0, 0x18, 0x4f, 0xc3, 0xf8, 0x6b, 0xda, 0x28,
	0x35, 0x4b, 0x2d, 0xd3, 0x39, 0x2a, 0x10, 0x36, 0x5b, 0xb5, 0xbb, 0x22, 0x0b, 0xcf, 0x32, 0xb6,
	0x24, 0xaa, 0xe4, 0x56, 0x84, 0xed, 0x07, 0x45, 0x38, 0x7c, 0x05, 0xe6, 0x5a, 0x39, 0xb2, 0xa0,
	0x74, 0x43, 0x97, 0x8a, 0x3b, 0xff, 0xe4, 0x7a, 0x2c, 0xc2, 0x69, 0xbe, 0xe2, 0x2d, 0x8d, 0xd7,
	0x7a, 0x47, 0xb3, 0x3b, 0xb0, 0x27, 0xa6, 0x21, 0x74, 0x46, 0xbf, 0xaf, 0xe4, 0x3b, 0x82, 0x1a,
	0xa3, 0x13, 0x46, 0xd3, 0x2f, 0xc1, 0xba, 0x8c, 0x55, 0xe5, 0x14, 0x05, 0xf6, 0x0f, 0xa8, 0xaa,
	0x4a, 0x59, 0xf4, 0xb2, 0x60, 0xaa, 0x09, 0xa6, 0x8f, 0x0b, 0xa6, 0x2a, 0xe3, 0x3e, 0x8a, 0xff,
	0x32, 0xfa, 0x27, 0xa8, 0x29, 0xf8, 0x62, 0xb5, 0x0c, 0x31, 0xae, 0x7c, 0x2d, 0xd3, 0xa9, 0x15,
	0x13, 0x5c, 0x86, 0x31, 0x23, 0x2a, 0xf8, 0xd7, 0xab, 0x35, 0x80, 0x4a, 0x51, 0x8c, 0x9e, 0x00,
	0x84, 0x79, 0xb6, 0xa9, 0x47, 0x85, 0x7b, 0x44, 0xca, 0x5d, 0xc5, 0xf4, 0x7b, 0x14, 0xeb, 0x40,
	0x75, 0xbd, 0x0f, 0xa7, 0x26, 0x07, 0x51, 0x5b, 0x2a, 0x0c, 0x84, 0x60, 0x7b, 0x9c, 0x44, 0x92,
	0x6f, 0x99, 0x88, 0xef, 0xe3, 0xf7, 0x00, 0xb7, 0x77, 0x86, 0x2a, 0x50, 0xfe, 0xe0, 0xfa, 0xdd,
	0x33, 0x6b, 0x0b, 0x01, 0x18, 0xe7, 0x83, 0xfe, 0x60, 0xe4, 0x5b, 0x1a, 0xaa, 0x03, 0xf8, 0xc4,
	0xed, 0x62, 0x32, 0x18, 0xf9, 0xd8, 0xd2, 0x79, 0xec, 0x1d, 0xc1, 0xf8, 0x1a, 0x5b, 0x25, 0x54,
	0x85, 0xdd, 0xd1, 0x85, 0xb2, 0xb6, 0x8f, 0xfb, 0x60, 0xae, 0x2d, 0x31, 0x32, 0x61, 0xc7, 0xbb,
	0xb8, 0x72, 0xcf, 0xbd, 0x9e, 0xb5, 0xc5, 0x51, 0xdc, 0x91, 0x7f, 0x36, 0x20, 0xde, 0x35, 0xee,
	0x49, 0x54, 0x82, 0x87, 0x3e, 0xf1, 0xba, 0x3e, 0xee, 0x59, 0x3a, 0x4f, 0xc6, 0x1f, 0x2f, 0x3d,
	0x82, 0x7b, 0x56, 0xc9, 0x59, 0xed, 0x80, 0xa2, 0x85, 0x5e, 0xc0, 0x6e, 0x9f, 0xce, 0x28, 0xe3,
	0x53, 0x5a, 0xbf, 0xbf, 0xff, 0xe1, 0xde, 0x9a, 0x47, 0x3d, 0xd9, 0x5b, 0x30, 0xe5, 0x1d, 0x48,
	0x19, 0x1f, 0xdd, 0xb9, 0x0e, 0x59, 0x79, 0xf0, 0x87, 0xab, 0xf9, 0x6c, 0x88, 0x1f, 0xda, 0xe9,
	0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2e, 0x84, 0xd8, 0x4e, 0xe7, 0x04, 0x00, 0x00,
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
	VerifyToken(ctx context.Context, in *TokenVerifyRequest, opts ...grpc.CallOption) (*TokenVerifyResponse, error)
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

func (c *tokenServiceClient) VerifyToken(ctx context.Context, in *TokenVerifyRequest, opts ...grpc.CallOption) (*TokenVerifyResponse, error) {
	out := new(TokenVerifyResponse)
	err := c.cc.Invoke(ctx, "/v1.TokenService/VerifyToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokenServiceServer is the server API for TokenService service.
type TokenServiceServer interface {
	Generate(context.Context, *TokenRequest) (*TokenResponse, error)
	VerifyToken(context.Context, *TokenVerifyRequest) (*TokenVerifyResponse, error)
}

// UnimplementedTokenServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTokenServiceServer struct {
}

func (*UnimplementedTokenServiceServer) Generate(ctx context.Context, req *TokenRequest) (*TokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Generate not implemented")
}
func (*UnimplementedTokenServiceServer) VerifyToken(ctx context.Context, req *TokenVerifyRequest) (*TokenVerifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyToken not implemented")
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

func _TokenService_VerifyToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenVerifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).VerifyToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.TokenService/VerifyToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).VerifyToken(ctx, req.(*TokenVerifyRequest))
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
		{
			MethodName: "VerifyToken",
			Handler:    _TokenService_VerifyToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "token_service.proto",
}
