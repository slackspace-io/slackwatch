use crate::config::Settings;
use git2::{
    Commit, Cred, ErrorCode, IndexAddOption, PushOptions, RemoteCallbacks, Repository, Signature,
};
use std::error::Error;
use std::fs::OpenOptions;
use std::io::Write;
use std::path::Path;

fn delete_local_repo() -> Result<(), std::io::Error> {
    let local_path = Path::new("/tmp/repos/");
    std::fs::remove_dir_all(local_path)?;
    Ok(())
}
/// Clone or open a repository based on the provided URL and local path.
fn clone_or_open_repo(
    repo_url: &str,
    repo_path: &Path,
    access_token: &str,
) -> Result<Repository, git2::Error> {
    match Repository::open(repo_path) {
        Ok(repo) => Ok(repo),
        Err(e) if e.code() == ErrorCode::NotFound => {
            let mut cb = RemoteCallbacks::new();
            log::info!("Setting credentials");
            cb.credentials(move |_url, _username, _allowed_types| {
                Cred::userpass_plaintext("x-access-token", access_token)
            });
            log::info!("Setting credentials Done");

            let mut fo = git2::FetchOptions::new();
            fo.remote_callbacks(cb);

            let mut builder = git2::build::RepoBuilder::new();
            builder.fetch_options(fo);
            log::info!("Building repo");
            builder.clone(repo_url, repo_path)
        }
        Err(e) => Err(e),
    }
}

/// Perform file modifications needed for your workload.
/// This is a placeholder function; implement your file editing logic here.
fn edit_files(local_path: &Path) -> Result<(), Box<dyn Error>> {
    // Example: Append a line to a file in the repository.
    let file_path = local_path.join("file_to_edit.txt");
    let mut file = OpenOptions::new().append(true).open(file_path)?;
    writeln!(file, "New line added by gitops operation.")?;
    Ok(())
}

/// Stage all changes in the repository.
fn stage_changes(repo: &Repository) -> Result<(), git2::Error> {
    let mut index = repo.index()?;
    index.add_all(["*"].iter(), IndexAddOption::DEFAULT, None)?;
    index.write()?;
    Ok(())
}

/// Commit the staged changes to the repository.
fn commit_changes<'a>(repo: &'a Repository, message: &str) -> Result<Commit<'a>, git2::Error> {
    let sig = Signature::now("Your Name", "your_email@example.com")?;
    let oid = repo.index()?.write_tree()?;
    let tree = repo.find_tree(oid)?;
    let parent_commit = find_last_commit(repo)?;
    let commit = repo.commit(Some("HEAD"), &sig, &sig, message, &tree, &[&parent_commit])?;
    Ok(repo.find_commit(commit)?)
}

/// Find the last commit in the repository.
fn find_last_commit(repo: &Repository) -> Result<Commit<'_>, git2::Error> {
    let obj = repo.head()?.resolve()?.peel(git2::ObjectType::Commit)?;
    obj.into_commit()
        .map_err(|_| git2::Error::from_str("Couldn't find commit"))
}

/// Push the changes to the remote repository.
fn push_changes(repo: &Repository, access_token: &str) -> Result<(), git2::Error> {
    let mut cb = RemoteCallbacks::new();
    cb.credentials(|_url, username_from_url, _allowed_types| {
        Cred::userpass_plaintext(username_from_url.unwrap(), access_token)
    });

    let mut opts = PushOptions::new();
    opts.remote_callbacks(cb);

    let mut remote = repo.find_remote("origin")?;
    remote.push(&["refs/heads/main:refs/heads/main"], Some(&mut opts))?;
    Ok(())
}

/// Main function to run git operations.
pub fn run_git_operations(settings: Settings) -> Result<(), Box<dyn Error>> {
    for gitops_config in settings.gitops {
        let repo_url = gitops_config.repository_url;
        let branch = gitops_config.branch;
        let name = gitops_config.name;
        let access_token_env_name = gitops_config.access_token_env_name;
        //try to get environment variable otherwise default
        let access_token = std::env::var(access_token_env_name).unwrap_or_default();
        log::info!("Access token: {}", access_token);
        let local_path = Path::new("/tmp/repos/").join(name);
        log::info!("Running git operations for repository: {}", repo_url);
        log::info!("Local path: {:?}", local_path);
        delete_local_repo()?;
        //let repo = clone_or_open_repo(&repo_url, &local_path, &access_token)?;

        //        edit_files(&local_path)?;
        //        stage_changes(&repo)?;
        //        commit_changes(&repo, "Automated commit by gitops.rs")?;
        //        push_changes(&repo, &access_token)?;
    }

    Ok(())
}
//) -> Result<(), Box<dyn Error>> {
//    let repo = clone_or_open_repo(repo_url, local_path)?;
//
//    //edit_files(local_path)?;
//    //stage_changes(&repo)?;
//    //commit_changes(&repo, "Automated commit by gitops.rs")?;
//    //push_changes(&repo, access_token)?;
//
//    Ok(())
//}
