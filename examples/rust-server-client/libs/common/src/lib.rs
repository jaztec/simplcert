use eyre::Result;
use tonic::transport::{Certificate, Identity};

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
