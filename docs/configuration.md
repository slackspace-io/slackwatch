# Configuration

## Configuration File

Slackwatch supports configuration via a TOML file. The default location is `~/.slackwatch/config.toml`. Here's an example:


#### Sections in the configuration file
- [System Configuration](#system-configuration)
- [Notifications Configuration](#notifications-configuration)
- [GitOps Configuration](#gitops-configuration)
- [Complete configuration file](#complete-configuration-file)


---

#### System Configuration
```toml
[system]
schedule = "0 0 9-22/2 * * *"
data_dir = "/app/slackwatch/data"
```
---

#### schedule
default: `0 0 9-22/2 * * *`

description: The `schedule` is a cron expression that defines when slackwatch should run. The default value is every 2 hours between 9am and 10pm.

--- 

#### data_dir
default: `~/.slackwatch/data`

description: The `data_dir` is the directory where slackwatch stores its data. This includes the state of the last run, and any other data that slackwatch needs to persist.

---

#### Notifications Configuration
```toml
[notifications.ntfy]
url = "http://localhost:9090"
topic = "slackwatch"
priorty = 1
reminder = "24h"
token = "dummy"
```
---


#### url
value: url(None)

description: the url of your ntfy server.

---

#### topic
value: string(None)

default: `slackwatch`

description: Topic to publish the notification to.

---

#### priority
value: int(1)

default: `1`

description: Priority of the notification.

---

#### reminder
value: string

default: `24h`

description: How often to resend the notification. 24h means every 24 hours.

---

#### token
value: string(None)

default: `dummy`

description: This is set to only prevent failure when the token is not set. In deployed scenarios you should set an environment variable named `SLACKWATCH_NOTIFICATIONS.NTFY.TOKEN` with the value of the token.

---

#### GitOps Configuration
```toml
[[gitops]]
name = "fleet-slack-house"
#repository_url = "https://github.com/slackspace-io/slackwatch.git"
repository_url = "https://github.com/slackspace-io/fleet-slack-house.git"
branch = "main"
access_token_env_name = "SLACKWATCH_TOKEN"
commit_message = "Updated by slackwatch"
commit_name = "slackwatch"
commit_email = "slackwatch@slackspace.io"

[[gitops]]
name = "noauth"
repository_url = "https://github.com/slackspace-io/slackwatch.git"
branch = "main"
access_token_env_name = "your_github_access_token_for_repoA"
commit_message = "Updated by slackwatch"
commit_name = "slackwatch"
commit_email = "slackwatch@slackspace.io"
```
Section Description: The `gitops` section is an array of configurations. The `name` field is  the key used to identify which gitops configuration to use. This should match the annotation `slackwatch.repo` on the deployment being watched.

---

#### name
value: string

description: The name of the gitops configuration. This should match the annotation `slackwatch.repo` on the deployment being watched.

---

#### repository_url
value: string

description: The URL of the git repository to sync with.

---

#### branch
value: string

description: The branch to sync with.

---

#### access_token_env_name
value: string

description: The name of the environment variable that contains the access token for the repository.


---

#### commit_message
value: string

description: The commit message to use when syncing changes.

---

#### commit_name
value: string

description: The name to use when syncing changes.

---

#### commit_email

value: string

description: The email to use when syncing changes.

---







[[gitops]]



Complete configuration file
```toml
[system]
schedule = "0 0 9-22/2 * * *"
data_dir = "/app/slackwatch/data"

[notifications.ntfy]
url = "http://localhost:9090"
topic = "slackwatch"
priorty = 1
reminder = "24h"
token = "dummy"


[[gitops]]
name = "fleet-slack-house"
#repository_url = "https://github.com/slackspace-io/slackwatch.git"
repository_url = "https://github.com/slackspace-io/fleet-slack-house.git"
branch = "main"
access_token_env_name = "SLACKWATCH_TOKEN"
commit_message = "Updated by slackwatch"
commit_name = "slackwatch"
commit_email = "slackwatch@slackspace.io"

[[gitops]]
name = "noauth"
repository_url = "https://github.com/slackspace-io/slackwatch.git"
branch = "main"
access_token_env_name = "your_github_access_token_for_repoA"
commit_message = "Updated by slackwatch"
commit_name = "slackwatch"
commit_email = "slackwatch@slackspace.io"
```
