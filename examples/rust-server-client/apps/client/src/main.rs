use eyre::Result;
use proto::greeter::{greeter_service_client::GreeterServiceClient, GreetingRequest};

#[tokio::main]
async fn main() -> Result<()> {
    let mut client = GreeterServiceClient::connect("server:8000").await?;

    let req = tonic::Request::new(GreetingRequest {
        name: "World".into_string(),
    });

    let resp = client.greet(req).await?;
    println!("Hello {}", resp.into_inner().greeting);

    Ok(())
}
