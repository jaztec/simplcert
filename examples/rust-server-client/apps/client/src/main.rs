use eyre::Result;
use proto::greeter::{greeter_service_client::GreeterServiceClient, GreetingRequest};
use std::time::Duration;
use tokio::{task, time};
use tonic::transport::{Channel, ClientTlsConfig};

#[tokio::main]
async fn main() -> Result<()> {
    let config = ClientTlsConfig::new()
        .ca_certificate(common::load_root_cert())
        .domain_name("server")
        .identity(common::load_identity("client").unwrap());
    let channel = Channel::from_static("https://server:8000")
        .tls_config(config)?
        .connect()
        .await?;

    let mut client = GreeterServiceClient::new(channel);

    greet(&mut client).await?;
    let forever = task::spawn(async move {
        let mut interval = time::interval(Duration::from_secs(10));

        loop {
            interval.tick().await;
            greet(&mut client).await.unwrap();
        }
    });

    forever.await?;

    Ok(())
}

async fn greet(client: &mut GreeterServiceClient<Channel>) -> Result<()> {
    let req = tonic::Request::new(GreetingRequest {
        name: "World".to_string(),
    });

    let resp = client.greet(req).await?;
    println!("{}", resp.into_inner().greeting);

    Ok(())
}
