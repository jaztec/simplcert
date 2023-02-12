use eyre::Result;
use tonic::transport::{Channel, ClientTlsConfig};
use proto::greeter::{greeter_service_client::GreeterServiceClient, GreetingRequest};

#[tokio::main]
async fn main() -> Result<()> {
    let config = ClientTlsConfig::new()
        .ca_certificate(common::load_root_cert())
        .domain_name("server")
        .identity(common::load_identity("client").unwrap());
    let channel = Channel::from_static("server:8000")
        .tls_config(config)?
        .connect()
        .await?;

    let mut client = GreeterServiceClient::new(channel);

    let req = tonic::Request::new(GreetingRequest {
        name: "World".to_string(),
    });

    let resp = client.greet(req).await?;
    println!("{}", resp.into_inner().greeting);

    Ok(())
}
