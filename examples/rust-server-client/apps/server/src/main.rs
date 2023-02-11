#[derive(Default, Debug)]
pub struct GreetingService;

#[tonic_async_trait]
impl proto::greeter::greeter_service_server::GreeterService for GreetingService {
    async fn greet()
}

fn main() {
    println!("Hello, world!");
}
