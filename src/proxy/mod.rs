/*
 *     Copyright 2024 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

use crate::shutdown;
use crate::Result as ClientResult;
use bytes::Bytes;
use http_body_util::{combinators::BoxBody, BodyExt, Empty, Full};
use hyper::server::conn::http1;
use hyper::service::service_fn;
use hyper::upgrade::Upgraded;
use hyper::{Method, Request, Response};
use hyper_util::rt::tokio::TokioIo;
use std::net::SocketAddr;
use tokio::net::TcpListener;
use tokio::net::TcpStream;
use tokio::pin;
use tokio::sync::mpsc;
use tracing::{info, instrument, Span};

// Proxy is the proxy server.
#[derive(Debug)]
pub struct Proxy {
    // addr is the address of the proxy server.
    addr: SocketAddr,

    // shutdown is used to shutdown the proxy server.
    shutdown: shutdown::Shutdown,

    // _shutdown_complete is used to notify the proxy server is shutdown.
    _shutdown_complete: mpsc::UnboundedSender<()>,
}

// Proxy implements the proxy server.
impl Proxy {
    // new creates a new Proxy.
    pub fn new(
        addr: SocketAddr,
        shutdown: shutdown::Shutdown,
        shutdown_complete_tx: mpsc::UnboundedSender<()>,
    ) -> Self {
        Self {
            addr,
            shutdown,
            _shutdown_complete: shutdown_complete_tx,
        }
    }

    // run starts the proxy server.
    #[instrument(skip_all)]
    pub async fn run(&self) -> ClientResult<()> {
        let listener = TcpListener::bind(self.addr).await?;

        // Start the proxy server and wait for it to finish.
        info!("proxy server listening on {}", self.addr);

        loop {
            // Clone the shutdown channel.
            let mut shutdown = self.shutdown.clone();

            let (tcp, remote_address) = listener.accept().await?;

            let io = TokioIo::new(tcp);
            info!("accepted connection from {}", remote_address);

            tokio::task::spawn(async move {
                // Pin the connection object so we can use tokio::select! below.
                let conn = http1::Builder::new()
                    .preserve_header_case(true)
                    .title_case_headers(true)
                    .serve_connection(io, service_fn(handler));
                pin!(conn);

                // Iterate the timeouts.  Use tokio::select! to wait on the
                // result of polling the connection itself,
                // and also on tokio::time::sleep for the current timeout duration.
                tokio::select! {
                    res = conn.as_mut() => {
                        // Polling the connection returned a result.
                        // In this case print either the successful or error result for the connection
                        // and break out of the loop.
                        match res {
                            Ok(()) => println!("after polling conn, no error"),
                            Err(e) =>  println!("error serving connection: {:?}", e),
                        };
                    }
                    _ = shutdown.recv() => {
                        // tokio::time::sleep returned a result.
                        // Call graceful_shutdown on the connection and continue the loop.
                        conn.as_mut().graceful_shutdown();
                    }
                }
            });
        }
    }
}

// handle starts to handle the request.
#[instrument(skip_all, fields(uri, method))]
pub async fn handler(
    request: Request<hyper::body::Incoming>,
) -> Result<Response<BoxBody<Bytes, hyper::Error>>, hyper::Error> {
    info!("handle request: {:?}", request);

    // Span record the uri and method.
    Span::current().record("uri", request.uri().to_string().as_str());
    Span::current().record("method", request.method().as_str());

    // Handle CONNECT request.
    if Method::CONNECT == request.method() {
        return handle_https(request).await;
    }

    return handle_http(request).await;
}

// handle_http handles the http request.
#[instrument(skip_all)]
pub async fn handle_http(
    request: Request<hyper::body::Incoming>,
) -> Result<Response<BoxBody<Bytes, hyper::Error>>, hyper::Error> {
    info!("handle http request: {:?}", request);
    let mut resp = Response::new(full("CONNECT must be to a socket address"));
    *resp.status_mut() = http::StatusCode::BAD_REQUEST;

    Ok(resp)
}

// handle_https handles the https request.
#[instrument(skip_all)]
pub async fn handle_https(
    request: Request<hyper::body::Incoming>,
) -> Result<Response<BoxBody<Bytes, hyper::Error>>, hyper::Error> {
    if let Some(addr) = host_addr(request.uri()) {
        tokio::task::spawn(async move {
            match hyper::upgrade::on(request).await {
                Ok(upgraded) => {
                    if let Err(e) = tunnel(upgraded, addr).await {
                        eprintln!("server io error: {}", e);
                    };
                }
                Err(e) => eprintln!("upgrade error: {}", e),
            }
        });

        Ok(Response::new(empty()))
    } else {
        eprintln!("CONNECT host is not socket addr: {:?}", request.uri());
        let mut resp = Response::new(full("CONNECT must be to a socket address"));
        *resp.status_mut() = http::StatusCode::BAD_REQUEST;

        Ok(resp)
    }
}

// empty returns an empty body.
fn empty() -> BoxBody<Bytes, hyper::Error> {
    Empty::<Bytes>::new()
        .map_err(|never| match never {})
        .boxed()
}

// full returns a body with the given chunk.
fn full<T: Into<Bytes>>(chunk: T) -> BoxBody<Bytes, hyper::Error> {
    Full::new(chunk.into())
        .map_err(|never| match never {})
        .boxed()
}

// host_addr returns the host address of the uri.
#[instrument(skip_all)]
fn host_addr(uri: &hyper::Uri) -> Option<String> {
    uri.authority().map(|auth| auth.to_string())
}

// tunnel proxies the data between the client and the remote server.
async fn tunnel(upgraded: Upgraded, addr: String) -> std::io::Result<()> {
    // Connect to remote server
    let mut server = TcpStream::connect(addr).await?;
    let mut upgraded = TokioIo::new(upgraded);

    // Proxying data
    let (from_client, from_server) =
        tokio::io::copy_bidirectional(&mut upgraded, &mut server).await?;

    // Print message when done
    println!(
        "client wrote {} bytes and received {} bytes",
        from_client, from_server
    );

    Ok(())
}
