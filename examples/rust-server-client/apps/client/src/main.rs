use eyre::Result;
use proto::greeter::{greeter_service_client::GreeterServiceClient, GreetingRequest};

#[tokio::main]
async fn main() -> Result<()> {
    let mut client = GreeterServiceClient::connect("server:8000").await?;

    let req = tonic::Request::new(GreetingRequest {
        name: "World".to_string(),
    });

    let resp = client.greet(req).await?;
    println!("{}", resp.into_inner().greeting);

    Ok(())
}
