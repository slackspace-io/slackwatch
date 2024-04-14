use oci_distribution::client::{Client, ClientConfig};
use oci_distribution::secrets::RegistryAuth;
use oci_distribution::Reference;

//pub async fn test_call() -> Result<Vec<String>, Box<dyn std::error::Error>> {
//    let reference = Reference::try_from("binwiederhier/ntfy")?;
//    let auth = RegistryAuth::Anonymous;
//    let config = ClientConfig::default();
//    let mut client = Client::new(config);
//
//    let tags = client.list_tags(&reference, &auth, None, None).await?;
//    //sort tags
//    let mut tags = tags.tags;
//    tags.sort();
//    let first_tag = tags.first().unwrap();
//    let last_tag = tags.last().unwrap();
//    println!("Available tags for {}: {:?}", reference, tags);
//    println!("First tag: {}", first_tag);
//    println!("Last tag: {}", last_tag);
//
//    Ok(tags)
//}

pub async fn get_tags_for_image(image: &str) -> Result<Vec<String>, Box<dyn std::error::Error>> {
    let reference = Reference::try_from(image)?;
    let auth = RegistryAuth::Anonymous;
    let config = ClientConfig::default();
    let mut client = Client::new(config);
    let max_tags = Some(1500);
    log::info!("Fetching tags for image: {:?}", reference.tag());
    let tags = client
        .list_tags(&reference, &auth, max_tags, reference.tag())
        .await?;
    log::info!("Available tags for {}: {:?}", reference, tags.tags);
    //length of tags
    log::info!("Number of tags: {}", tags.tags.len());

    Ok(tags.tags)
}
