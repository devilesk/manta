// Code generated by protoc-gen-go.
// source: uifontfile_format.proto
// DO NOT EDIT!

package dota

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type CUIFontFilePB struct {
	FontFileName     *string `protobuf:"bytes,1,opt,name=font_file_name,json=fontFileName" json:"font_file_name,omitempty"`
	OpentypeFontData []byte  `protobuf:"bytes,2,opt,name=opentype_font_data,json=opentypeFontData" json:"opentype_font_data,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CUIFontFilePB) Reset()                    { *m = CUIFontFilePB{} }
func (m *CUIFontFilePB) String() string            { return proto.CompactTextString(m) }
func (*CUIFontFilePB) ProtoMessage()               {}
func (*CUIFontFilePB) Descriptor() ([]byte, []int) { return fileDescriptor43, []int{0} }

func (m *CUIFontFilePB) GetFontFileName() string {
	if m != nil && m.FontFileName != nil {
		return *m.FontFileName
	}
	return ""
}

func (m *CUIFontFilePB) GetOpentypeFontData() []byte {
	if m != nil {
		return m.OpentypeFontData
	}
	return nil
}

type CUIFontFilePackagePB struct {
	PackageVersion     *uint32                                        `protobuf:"varint,1,req,name=package_version,json=packageVersion" json:"package_version,omitempty"`
	EncryptedFontFiles []*CUIFontFilePackagePB_CUIEncryptedFontFilePB `protobuf:"bytes,2,rep,name=encrypted_font_files,json=encryptedFontFiles" json:"encrypted_font_files,omitempty"`
	XXX_unrecognized   []byte                                         `json:"-"`
}

func (m *CUIFontFilePackagePB) Reset()                    { *m = CUIFontFilePackagePB{} }
func (m *CUIFontFilePackagePB) String() string            { return proto.CompactTextString(m) }
func (*CUIFontFilePackagePB) ProtoMessage()               {}
func (*CUIFontFilePackagePB) Descriptor() ([]byte, []int) { return fileDescriptor43, []int{1} }

func (m *CUIFontFilePackagePB) GetPackageVersion() uint32 {
	if m != nil && m.PackageVersion != nil {
		return *m.PackageVersion
	}
	return 0
}

func (m *CUIFontFilePackagePB) GetEncryptedFontFiles() []*CUIFontFilePackagePB_CUIEncryptedFontFilePB {
	if m != nil {
		return m.EncryptedFontFiles
	}
	return nil
}

type CUIFontFilePackagePB_CUIEncryptedFontFilePB struct {
	EncryptedContents []byte `protobuf:"bytes,1,opt,name=encrypted_contents,json=encryptedContents" json:"encrypted_contents,omitempty"`
	XXX_unrecognized  []byte `json:"-"`
}

func (m *CUIFontFilePackagePB_CUIEncryptedFontFilePB) Reset() {
	*m = CUIFontFilePackagePB_CUIEncryptedFontFilePB{}
}
func (m *CUIFontFilePackagePB_CUIEncryptedFontFilePB) String() string {
	return proto.CompactTextString(m)
}
func (*CUIFontFilePackagePB_CUIEncryptedFontFilePB) ProtoMessage() {}
func (*CUIFontFilePackagePB_CUIEncryptedFontFilePB) Descriptor() ([]byte, []int) {
	return fileDescriptor43, []int{1, 0}
}

func (m *CUIFontFilePackagePB_CUIEncryptedFontFilePB) GetEncryptedContents() []byte {
	if m != nil {
		return m.EncryptedContents
	}
	return nil
}

func init() {
	proto.RegisterType((*CUIFontFilePB)(nil), "dota.CUIFontFilePB")
	proto.RegisterType((*CUIFontFilePackagePB)(nil), "dota.CUIFontFilePackagePB")
	proto.RegisterType((*CUIFontFilePackagePB_CUIEncryptedFontFilePB)(nil), "dota.CUIFontFilePackagePB.CUIEncryptedFontFilePB")
}

func init() { proto.RegisterFile("uifontfile_format.proto", fileDescriptor43) }

var fileDescriptor43 = []byte{
	// 255 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x90, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0xdd, 0xa8, 0x07, 0xd7, 0xb4, 0xea, 0x52, 0x34, 0x78, 0x0a, 0x45, 0x30, 0x07, 0x0d,
	0xe8, 0x4f, 0x68, 0xb5, 0xea, 0x45, 0x64, 0x41, 0xaf, 0xcb, 0xb0, 0x99, 0x48, 0xb0, 0xd9, 0x59,
	0x92, 0x51, 0xe8, 0xcd, 0x3f, 0xec, 0x7f, 0x90, 0x6d, 0xd3, 0x56, 0x24, 0xc7, 0x7d, 0xef, 0x9b,
	0x79, 0xfb, 0x46, 0x9e, 0x7d, 0x56, 0x25, 0x39, 0x2e, 0xab, 0x39, 0x9a, 0x92, 0x9a, 0x1a, 0x38,
	0xf7, 0x0d, 0x31, 0xa9, 0xbd, 0x82, 0x18, 0xc6, 0x56, 0x0e, 0xa6, 0xaf, 0x4f, 0x33, 0x72, 0x3c,
	0xab, 0xe6, 0xf8, 0x32, 0x51, 0x17, 0x72, 0x18, 0x78, 0xb3, 0x1c, 0x70, 0x50, 0x63, 0x22, 0x52,
	0x91, 0x1d, 0xe8, 0xb8, 0xec, 0x98, 0x67, 0xa8, 0x51, 0x5d, 0x49, 0x45, 0x1e, 0x1d, 0x2f, 0x7c,
	0xd8, 0xea, 0xd8, 0x14, 0xc0, 0x90, 0x44, 0xa9, 0xc8, 0x62, 0x7d, 0xbc, 0x76, 0xc2, 0xd6, 0x3b,
	0x60, 0x18, 0xff, 0x08, 0x39, 0xfa, 0x9b, 0x02, 0xf6, 0x03, 0xde, 0x43, 0xd8, 0xa5, 0x3c, 0xf2,
	0xab, 0x87, 0xf9, 0xc2, 0xa6, 0xad, 0xc8, 0x25, 0x22, 0x8d, 0xb2, 0x81, 0x1e, 0x76, 0xf2, 0xdb,
	0x4a, 0x55, 0x56, 0x8e, 0xd0, 0xd9, 0x66, 0xe1, 0x19, 0x0b, 0xb3, 0xf9, 0x5f, 0x9b, 0x44, 0xe9,
	0x6e, 0x76, 0x78, 0x7b, 0x93, 0x87, 0x2e, 0x79, 0x5f, 0x44, 0x10, 0xef, 0xd7, 0x93, 0xdb, 0x9a,
	0x5a, 0xe1, 0x7f, 0xb1, 0x3d, 0x7f, 0x90, 0xa7, 0xfd, 0xb4, 0xba, 0x96, 0x5b, 0xde, 0x58, 0x72,
	0x8c, 0x8e, 0xdb, 0xe5, 0x61, 0x62, 0x7d, 0xb2, 0x71, 0xa6, 0x9d, 0x31, 0xd9, 0x7f, 0x14, 0xdf,
	0x62, 0xe7, 0x37, 0x00, 0x00, 0xff, 0xff, 0xe0, 0xdc, 0xb8, 0xe1, 0x7b, 0x01, 0x00, 0x00,
}
