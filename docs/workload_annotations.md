# Workload Annotations

Only workloads with the annotation `slackwatch.enabled` set to `true` will be watched by slackwatch.

Additional annotations can be set to customize the behavior of slackwatch for a given workload.

## Main annotations

### `slackwatch.enabled`
description: Set to `true` to enable slackwatch for this workload. If this annotation is not present or set to `false`, slackwatch will not watch this workload.


### `slackwatc.include`
description: A comma-seperated list of regex patterns to apply to tags. Only tags which match will be evaluated by slackwatch.

### `slackwatch.exclude`
description: A comma-seperated list of regex patterns to apply to tags. Tags which match will be ignored by slackwatch during evaluation.


## If using automated gitops commits

### `slackwatch.repo`
description: The name of the gitops configuration to use for syncing changes. This should match the `name` field in the gitops configuration.

### `slackwatch.directory`
description: The directory which your application deployment files are located, within your repo. By default it expects the name of the workload to match the directory name. Slackwatch will walk subdirectories below this directory to find deployment files containing the expected tag. This is only used when `slackwatch.repo` is defined.
