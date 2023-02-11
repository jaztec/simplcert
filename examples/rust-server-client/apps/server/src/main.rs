use eyre::Result;
use tonic::transport::Server;
use proto::greeter::{GreetingResponse, GreetingRequest};
use proto::greeter::greeter_service_server::GreeterServiceServer;

#[derive(Default, Debug)]
pub struct GreetingService;

#[tonic::async_trait]
impl proto::greeter::greeter_service_server::GreeterService for GreetingService {
    async fn greet(
        &self,
        request: tonic::Request<GreetingRequest>,
    ) -> Result<tonic::Response<GreetingResponse>, tonic::Status> {
        Ok(
            tonic::Response::new(GreetingResponse { greeting: format!("Hello {}", request.into_inner().name) })
        )
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    let service = GreeterServiceServer::new(GreetingService{});

    Server::builder()
        .add_service(service)
        .serve(":8000".parse().unwrap())
        .await?;

    Ok(())
}
