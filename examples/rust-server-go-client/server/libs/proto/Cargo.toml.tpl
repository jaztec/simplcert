[package]
name = "proto"
version = "0.1.0"
edition = "2021"

[dependencies]
bytes = "1.1.0"
prost = "0.11.0"
pbjson = "0.5"
pbjson-types = "0.5"
serde = "1.0"
tonic = { version = "0.8", features = ["gzip"] }


[features]
default = ["proto_full"]
## @@protoc_insertion_point(features)