use eyre::Result;
use proto::greeter::greeter_service_server::GreeterServiceServer;
use proto::greeter::{GreetingRequest, GreetingResponse};
use tonic::transport::{Server, ServerTlsConfig, Certificate, Identity};

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
        .client_ca_root(load_root_cert())
        .identity(load_identity("server").unwrap());

    let service = GreeterServiceServer::new(GreetingService {});

    Server::builder()
        .tls_config(config)?
        .add_service(service)
        .serve("0.0.0.0:8000".parse().unwrap())
        .await?;

    Ok(())
}

pub fn load_root_cert() -> Certificate {
    load_cert("/certs/root-ca.crt").unwrap()
}

pub fn load_cert(path: &str) -> Result<Certificate> {
    let pem = std::fs::read(path)?;
    Ok(Certificate::from_pem(pem))
}

pub fn load_identity(name: &str) -> Result<Identity> {
    let crt = load_cert(format!("/certs/{}.crt", name).as_str())?;
    let key = std::fs::read(format!("/certs/{}.key", name))?;

    Ok(Identity::from_pem(crt, key))
}
