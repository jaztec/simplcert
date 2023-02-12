use eyre::Result;
use proto::greeter::greeter_service_server::GreeterServiceServer;
use proto::greeter::{GreetingRequest, GreetingResponse};
use tonic::transport::{Server, ServerTlsConfig};

#[derive(Default, Debug)]
pub struct GreetingService;

#[tonic::async_trait]
impl proto::greeter::greeter_service_server::GreeterService for GreetingService {
    async fn greet(
        &self,
        request: tonic::Request<GreetingRequest>,
    ) -> Result<tonic::Response<GreetingResponse>, tonic::Status> {
        Ok(tonic::Response::new(GreetingResponse {
            greeting: format!("Hello {}", request.into_inner().name),
        }))
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    let config = ServerTlsConfig::new()
        .client_ca_root(common::load_root_cert())
        .identity(common::load_identity("server").unwrap());

    let service = GreeterServiceServer::new(GreetingService {});

    Server::builder()
        .tls_config(config)?
        .add_service(service)
        .serve("0.0.0.0:8000".parse().unwrap())
        .await?;

    Ok(())
}
