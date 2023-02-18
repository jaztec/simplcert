// @generated
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GreetingRequest {
    #[prost(string, tag="1")]
    pub name: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GreetingResponse {
    #[prost(string, tag="1")]
    pub greeting: ::prost::alloc::string::String,
}
/// Encoded file descriptor set for the `greeter` package
pub const FILE_DESCRIPTOR_SET: &[u8] = &[
    0x0a, 0xf3, 0x04, 0x0a, 0x0e, 0x67, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72,
    0x6f, 0x74, 0x6f, 0x12, 0x07, 0x67, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x22, 0x25, 0x0a, 0x0f,
    0x47, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
    0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
    0x61, 0x6d, 0x65, 0x22, 0x2e, 0x0a, 0x10, 0x47, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x52,
    0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x67, 0x72, 0x65, 0x65, 0x74,
    0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x67, 0x72, 0x65, 0x65, 0x74,
    0x69, 0x6e, 0x67, 0x32, 0x50, 0x0a, 0x0e, 0x47, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x53, 0x65,
    0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3e, 0x0a, 0x05, 0x47, 0x72, 0x65, 0x65, 0x74, 0x12, 0x18,
    0x2e, 0x67, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x2e, 0x47, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e,
    0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x67, 0x72, 0x65, 0x65, 0x74,
    0x65, 0x72, 0x2e, 0x47, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
    0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x85, 0x01, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x72,
    0x65, 0x65, 0x74, 0x65, 0x72, 0x42, 0x0d, 0x47, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x50,
    0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2b, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73,
    0x2f, 0x72, 0x75, 0x73, 0x74, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x67, 0x6f, 0x2d,
    0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x70, 0x72,
    0x6f, 0x74, 0x6f, 0xa2, 0x02, 0x03, 0x47, 0x58, 0x58, 0xaa, 0x02, 0x07, 0x47, 0x72, 0x65, 0x65,
    0x74, 0x65, 0x72, 0xca, 0x02, 0x07, 0x47, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0xe2, 0x02, 0x13,
    0x47, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64,
    0x61, 0x74, 0x61, 0xea, 0x02, 0x07, 0x47, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x4a, 0x9e, 0x02,
    0x0a, 0x06, 0x12, 0x04, 0x00, 0x00, 0x0f, 0x01, 0x0a, 0x08, 0x0a, 0x01, 0x0c, 0x12, 0x03, 0x00,
    0x00, 0x12, 0x0a, 0x08, 0x0a, 0x01, 0x08, 0x12, 0x03, 0x02, 0x00, 0x42, 0x0a, 0x09, 0x0a, 0x02,
    0x08, 0x0b, 0x12, 0x03, 0x02, 0x00, 0x42, 0x0a, 0x08, 0x0a, 0x01, 0x02, 0x12, 0x03, 0x03, 0x00,
    0x10, 0x0a, 0x0a, 0x0a, 0x02, 0x06, 0x00, 0x12, 0x04, 0x05, 0x00, 0x07, 0x01, 0x0a, 0x0a, 0x0a,
    0x03, 0x06, 0x00, 0x01, 0x12, 0x03, 0x05, 0x08, 0x16, 0x0a, 0x0b, 0x0a, 0x04, 0x06, 0x00, 0x02,
    0x00, 0x12, 0x03, 0x06, 0x02, 0x3a, 0x0a, 0x0c, 0x0a, 0x05, 0x06, 0x00, 0x02, 0x00, 0x01, 0x12,
    0x03, 0x06, 0x06, 0x0b, 0x0a, 0x0c, 0x0a, 0x05, 0x06, 0x00, 0x02, 0x00, 0x02, 0x12, 0x03, 0x06,
    0x0c, 0x1b, 0x0a, 0x0c, 0x0a, 0x05, 0x06, 0x00, 0x02, 0x00, 0x03, 0x12, 0x03, 0x06, 0x26, 0x36,
    0x0a, 0x0a, 0x0a, 0x02, 0x04, 0x00, 0x12, 0x04, 0x09, 0x00, 0x0b, 0x01, 0x0a, 0x0a, 0x0a, 0x03,
    0x04, 0x00, 0x01, 0x12, 0x03, 0x09, 0x08, 0x17, 0x0a, 0x0b, 0x0a, 0x04, 0x04, 0x00, 0x02, 0x00,
    0x12, 0x03, 0x0a, 0x02, 0x12, 0x0a, 0x0c, 0x0a, 0x05, 0x04, 0x00, 0x02, 0x00, 0x05, 0x12, 0x03,
    0x0a, 0x02, 0x08, 0x0a, 0x0c, 0x0a, 0x05, 0x04, 0x00, 0x02, 0x00, 0x01, 0x12, 0x03, 0x0a, 0x09,
    0x0d, 0x0a, 0x0c, 0x0a, 0x05, 0x04, 0x00, 0x02, 0x00, 0x03, 0x12, 0x03, 0x0a, 0x10, 0x11, 0x0a,
    0x0a, 0x0a, 0x02, 0x04, 0x01, 0x12, 0x04, 0x0d, 0x00, 0x0f, 0x01, 0x0a, 0x0a, 0x0a, 0x03, 0x04,
    0x01, 0x01, 0x12, 0x03, 0x0d, 0x08, 0x18, 0x0a, 0x0b, 0x0a, 0x04, 0x04, 0x01, 0x02, 0x00, 0x12,
    0x03, 0x0e, 0x02, 0x16, 0x0a, 0x0c, 0x0a, 0x05, 0x04, 0x01, 0x02, 0x00, 0x05, 0x12, 0x03, 0x0e,
    0x02, 0x08, 0x0a, 0x0c, 0x0a, 0x05, 0x04, 0x01, 0x02, 0x00, 0x01, 0x12, 0x03, 0x0e, 0x09, 0x11,
    0x0a, 0x0c, 0x0a, 0x05, 0x04, 0x01, 0x02, 0x00, 0x03, 0x12, 0x03, 0x0e, 0x14, 0x15, 0x62, 0x06,
    0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
];
include!("greeter.serde.rs");
include!("greeter.tonic.rs");
// @@protoc_insertion_point(module)