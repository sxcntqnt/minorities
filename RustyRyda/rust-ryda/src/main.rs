use tonic::{transport::Server, Request, Response, Status};
use greet::greeter_server::{Greeter, GreeterServer};
use greet::{HelloRequest, HelloReply};

// The struct that implements the Greeter trait.
#[derive(Default)]
pub struct MyGreeter {}

#[tonic::async_trait]
impl Greeter for MyGreeter {
    async fn say_hello(
        &self,
        request: Request<HelloRequest>,
    ) -> Result<Response<HelloReply>, Status> {
        let reply = HelloReply {
            message: format!("Hello, {}!", request.get_ref().name),
        };
        Ok(Response::new(reply))
    }
}

pub mod greet {
    tonic::include_proto!("greet");
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?;
    let greeter = MyGreeter::default();

    Server::builder()
        .add_service(GreeterServer::new(greeter))
        .serve(addr)
        .await?;

    Ok(())
}
