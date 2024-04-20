package analysis

import (
	"dbwf-ls/lsp"
	"fmt"
	"log"
	"regexp"
)

func wordAtCursor(line string, position lsp.Position, logger *log.Logger) (string, error) {
	re, err := regexp.Compile("\\W")
	if err != nil {
		logger.Printf("Regexp Compile %s", err)
		return "", err
	}

	// Because the flocking cursor is 1 step ahead of the line while typing
	// So this can fail, quietly, damn.
	char := position.Character
	if char == len(line) {
		char--
	}

	if loc := re.FindStringIndex(line[char : char+1]); loc != nil {
		return "", nil
	}

	start, end := 0, 0
	if locs := re.FindAllStringIndex(line, -1); locs != nil {
		for _, loc := range locs {
			if loc[0] > position.Character {
				end = loc[0]
			}
			if loc[1] <= position.Character {
				start = loc[1]
			}
			if end != 0 {
				break
			}
		}
	}
	if end == 0 {
		end = len(line)
	}

	return line[start:end], nil
}

func leadingSpaces(line string, logger *log.Logger) (string, error) {
	re, err := regexp.Compile("^\\s*")
	if err != nil {
		logger.Panicf("Regexp Compile %s", err)
		return "", err
	}

	return re.FindString(line), nil
}

type Keyword struct {
	hover       lsp.MarkupContent
	diag        Diag
	completions []Completion
}

type Completion struct {
	insertText, detail string
	kind               int
	documentation      lsp.MarkupContent
}

type Diag struct {
	help     string
	severity int
}

// A whole lot of keywords
var Keywords = map[string]Keyword{
	"name": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
				"`name` string <= 4096 characters",
				"Default `\"Untitled\"`",
				"Example `\"A multitask job\"`",
				"An optional name for the job. The maximum length is 4096 bytes in UTF-8 encoding.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#name)",
			),
		},
		diag: Diag{
			help:     "`name` is missing. Hint: start by typing `name`",
			severity: 2,
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n",
					"name: \"Untitled workflow\"",
					"description: \"Workflow description\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "name and description declaration",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n",
						"Snippet for name and description declaration",
						"name: \"Untitled workflow\"",
						"description: \"Workflow description\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n", "name: \"Untitled workflow\""),
				kind:       lsp.CompletionItemKind["Snippet"],
				detail:     "name declaration",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n",
						"Snippet for name declaration",
						"name: \"Untitled workflow\"",
					),
				},
			},
			{
				insertText: "name",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	"description": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n",
				"`description` string <= 1024 characters",
				"Example `\"This job contain multiple tasks that are required to produce the weekly shark sightings report.\"`",
				"optional description for the job. The maximum length is 1024 characters in UTF-8 encoding.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#description)",
			),
		},
		diag: Diag{
			help:     "`description` is missing. Hint: start by typing `description`",
			severity: 4,
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n", "description: \"Meaningful description\""),
				kind:       lsp.CompletionItemKind["Snippet"],
				detail:     "description declaration",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n",
						"Snippet for description declaration",
						"description: \"Meaningful description\"",
					),
				},
			},
			{
				insertText: "description",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	"email_notifications": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n",
				"`email_notifications` object",
				"Default `{}`",
				"An optional set of email addresses that is notified when runs of this job begin or complete as well as when this job is deleted.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#email_notifications)",
			),
		},
		diag: Diag{
			help:     "`email_notifications` is missing. Hint: start by typing `email`",
			severity: 2,
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n",
					"email_notifications:",
					"  on_failure:\n    - \"some.one@some.org\"",
					"  on_duration_warning_threshold_exceeded:\n    - \"some.one@some.org\"",
					"  on_success:\n    - \"some.one@some.org\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "email noti declaration (recommended)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n",
						"Snippet for email noti declaration",
						"email_notifications:",
						"  on_failure:\n    - \"some.one@some.org\"",
						"  on_duration_warning_threshold_exceeded:\n    - \"some.one@some.org\"",
						"  on_success:\n    - \"some.one@some.org\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n",
					"email_notifications:",
					"  on_start:\n    - \"some.one@some.org\"",
					"  on_failure:\n    - \"some.one@some.org\"",
					"  on_duration_warning_threshold_exceeded:\n    - \"some.one@some.org\"",
					"  on_success:\n    - \"some.one@some.org\"",
					"  no_alert_for_skipped_runs: false",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "email noti declaration (all)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n%s\n%s\n",
						"Snippet for email noti declaration",
						"email_notifications:",
						"  on_start:\n    - \"some.one@some.org\"",
						"  on_failure:\n    - \"some.one@some.org\"",
						"  on_duration_warning_threshold_exceeded:\n    - \"some.one@some.org\"",
						"  on_success:\n    - \"some.one@some.org\"",
						"  no_alert_for_skipped_runs: false",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n",
					"email_notifications:",
					"  on_failure:\n    - \"some.one@some.org\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "email noti declaration (minimal)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n",
						"Snippet for email noti declaration",
						"email_notifications:",
						"  on_failure:\n    - \"some.one@some.org\"",
					),
				},
			},
			{
				insertText: "email_notification",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	// "webhook_notifications": {},
	"notification_settings": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
				"`notification_settings` object",
				"Default `{}`",
				"Example:",
				"```yaml",
				"notification_settings:",
				"  no_alert_for_skipped_runs: false",
				"  no_alert_for_canceled_runs: false",
				"```",
				"Optional notification settings that are used when sending notifications to each of the `email_notifications` and `webhook_notifications` for this job.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#notification_settings)",
			),
		},
		diag: Diag{
			help:     "`notification_settings` is missing. Hint: start by typing `noti`",
			severity: 4,
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n",
					"notification_settings:",
					"  no_alert_for_skipped_runs: false",
					"  no_alert_for_canceled_runs: false",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "notification settings",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n",
						"Snippet for notification settings",
						"notification_settings:",
						"  no_alert_for_skipped_runs: false",
						"  no_alert_for_canceled_runs: false",
					),
				},
			},
			{
				insertText: "notification_settings",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	"timeout_seconds": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
				"`timeout_seconds` int32",
				"Default `0`",
				"Example `86400`",
				"An optional timeout applied to each run of this job. A value of 0 means no timeout.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#timeout_seconds)",
			),
		},
		diag: Diag{
			help:     "`timeout_seconds` is missing. Hint: start by typing `timeout`",
			severity: 2,
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n", "timeout_seconds: 0"),
				kind:       lsp.CompletionItemKind["Snippet"],
				detail:     "timeout in seconds",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n",
						"Snippet for setting workflow timeout",
						"timeout_seconds: 0",
					),
				},
			},
			{
				insertText: "timeout_seconds",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	"health": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
				"`health` object",
				"An optional set of health rules that can be defined for this job.",
				"Example",
				"```yaml",
				"health:",
				"  rules:",
				"    - metrics: \"RUN_DURATION_SECONDS\"",
				"      op: \"GREATER_THAN\"",
				"      value: 10800",
				"```",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#health)",
			),
		},
		diag: Diag{
			help:     "`health` is missing. Hint: start by typing `health`",
			severity: 2,
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
					"health:",
					"  rules:",
					"    - metric: \"RUN_DURATION_SECONDS\"",
					"      op: \"GREATER_THAN\"",
					"      value: 10800 # 3 hours",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "workflow health",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n%s\n",
						"Snippet for setting workflow health alert",
						"health:",
						"  rules:",
						"    - metric: \"RUN_DURATION_SECONDS\"",
						"      op: \"GREATER_THAN\"",
						"      value: 10800 # 3 hours",
					),
				},
			},
			{
				insertText: "health",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	"schedule": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
				"`schedule` object",
				"An optional periodic schedule for this job. The default behavior is that the job only runs when triggered by clicking “Run Now” in the Jobs UI or sending an API request to `runNow`.",
				"Example",
				"```yaml",
				"schedule:",
				"  quartz_cronz_expression: \"0 0 0 * * ?\"",
				"  timezone_id: \"UTC\"",
				"  pause_status: `\"PAUSED\"`",
				"```",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#schedule)",
			),
		},
		diag: Diag{
			help:     "`schedule` is missing. Hint: start by typing `schedule`",
			severity: 2,
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n",
					"schedule:",
					"  quartz_cron_expression: \"0 0 0 * * ?\" # Everyday at 0am",
					"  timezone_id: \"UTC\"",
					"  pause_status: \"UNPAUSED\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "schedule (unpaused)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n",
						"Snippet to schedule workflow (unpaused)",
						"schedule:",
						"  quartz_cron_expression: \"0 0 0 * * ?\" # Everyday at 0am",
						"  timezone_id: \"UTC\"",
						"  pause_status: \"UNPAUSED\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n",
					"schedule:",
					"  quartz_cron_expression: \"0 0 0 * * ?\" # Everyday at 0am",
					"  timezone_id: \"UTC\"",
					"  pause_status: \"PAUSED\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "schedule (paused)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n",
						"Snippet to schedule workflow (paused)",
						"schedule:",
						"  quartz_cron_expression: \"0 0 0 * * ?\" # Everyday at 0am",
						"  timezone_id: \"UTC\"",
						"  pause_status: \"PAUSED\"",
					),
				},
			},
			{
				insertText: "schedule",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	// "trigger": {},
	// "continuous": {},
	"max_concurrent_runs": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
				"`max_concurrent_runs` int32",
				"Default `1`",
				"Example `10`",
				"An optional maximum allowed number of concurrent runs of the job. Set this value if you want to be able to execute multiple runs of the same job concurrently. This is useful for example if you trigger your job on a frequent schedule and want to allow consecutive runs to overlap with each other, or if you want to trigger multiple runs which differ by their input parameters. This setting affects only new runs. For example, suppose the job’s concurrency is 4 and there are 4 concurrent active runs. Then setting the concurrency to 3 won’t kill any of the active runs. However, from then on, new runs are skipped unless there are fewer than 3 active runs. This value cannot exceed 1000. Setting this value to `0` causes all new runs to be skipped.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#max_concurrent_runs)",
			),
		},
		diag: Diag{
			help:     "`max_concurrent_runs` is missing. Hint: start by typing `max`",
			severity: 4,
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n", "max_concurrent_runs: 1"),
				kind:       lsp.CompletionItemKind["Snippet"],
				detail:     "max concurrent runs",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n",
						"Snippet for setting concurrency",
						"max_concurrent_runs: 1",
					),
				},
			},
			{
				insertText: "max_concurrent_runs",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	"task_key": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n",
				"`task_key` string [ 1 .. 100 ] characters ^[\\w\\-\\_]+$",
				"Example: `\"Task_Key\"`",
				"A unique name for the task. This field is used to refer to this task from other tasks. This field is required and must be unique within its parent job. On Update or Reset, this field is used to reference the tasks to be updated or reset.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#tasks-task_key)",
			),
		},
	},
	"tasks": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n",
				"`tasks` Array of object <= 100 items",
				"A list of task specifications to be executed by this job.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#tasks)",
			),
		},
		diag: Diag{
			help:     "`tasks` is missing. Hint: start by typing `task`",
			severity: 1,
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
					"tasks:",
					"  - task_key: \"task_name\"",
					"    description: \"task description\"",
					"    existing_cluster_id: \"some_cluster_id\"",
					"    <task type declaration>",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "tasks (existing cluster)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n%s\n",
						"Snippet for declaring tasks",
						"tasks:",
						"  - task_key: \"task_name\"",
						"    description: \"task description\"",
						"    existing_cluster_id: \"some_cluster_id\"",
						"    <task type declaration>",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
					"tasks:",
					"  - task_key: \"task_name\"",
					"    description: \"task description\"",
					"    job_cluster_key: \"some_cluster_key\"",
					"    <task type declaration>",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "tasks (cluster key)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n%s\n",
						"Snippet for declaring tasks",
						"tasks:",
						"  - task_key: \"task_name\"",
						"    description: \"task description\"",
						"    job_cluster_key: \"some_cluster_key\"",
						"    <task type declaration>",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n", "- task_key: \"task_name\""),
				kind:       lsp.CompletionItemKind["Snippet"],
				detail:     "task key (name only)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n",
						"Snippet for declaring tasks",
						"- task_key: \"task_name\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
					"- task_key: \"task_name\"",
					"  description: \"task description\"",
					"  existing_cluster_id: \"some_cluster_id\"",
					"  <task type declaration>",
					"  depends_on:",
					"    - task_key: \"some_task_key\"",
					"  run_if: \"ALL_SUCCESS\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "single task (existing cluster)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
						"Snippet for declaring single task",
						"- task_key: \"task_name\"",
						"  description: \"task description\"",
						"  existing_cluster_id: \"some_cluster_id\"",
						"  <task type declaration>",
						"  depends_on:",
						"    - task_key: \"some_task_key\"",
						"  run_if: \"ALL_SUCCESS\"",
					)},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
					"- task_key: \"task_name\"",
					"  description: \"task description\"",
					"  job_cluster_key: \"some_cluster_key\"",
					"  <task type declaration>",
					"  depends_on:",
					"    - task_key: \"some_task_key\"",
					"  run_if: \"ALL_SUCCESS\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "single task (cluster key)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
						"Snippet for declaring single task",
						"- task_key: \"task_name\"",
						"  description: \"task description\"",
						"  job_cluster_key: \"some_cluster_key\"",
						"  <task type declaration>",
						"  depends_on:",
						"    - task_key: \"some_task_key\"",
						"  run_if: \"ALL_SUCCESS\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n",
					"spark_python_task:",
					"  python_file: \"file:/path/to/file\"",
					"  parameters:",
					"    - \"param-key=param-value\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "spark python task",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n",
						"Snippet for declaring single task",
						"spark_python_task:",
						"  python_file: \"file:/path/to/file\"",
						"  parameters:",
						"    - \"param-key=param-value\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
					"python_wheel_task:",
					"  package_name: \"some_package\"",
					"  entry_point: \"some_entry_point\"",
					"  parameters:",
					"    - \"param-key=param-value\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "python wheel task",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n%s\n",
						"Snippet for declaring single task",
						"python_wheel_task:",
						"  package_name: \"some_package\"",
						"  entry_point: \"some_entry_point\"",
						"  parameters:",
						"    - \"param-key=param-value\"",
					),
				},
			},
		},
	},
	"depends_on": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
				"`depends_on` Array of object",
				"An optional array of objects specifying the dependency graph of the task. All tasks specified in this field must complete before executing this task. The task will run only if the `run_if` condition is true. The key is `task_key`, and the value is the name assigned to the dependent task.",
				"Example:",
				"```yaml",
				"depends_on:",
				"  - task_key: \"some_task\"",
				"```",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#tasks-depends_on)",
			),
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n",
					"depends_on:",
					"  - task_key: \"task_name\"",
					"run_if: \"ALL_SUCCESS\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "depends on",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n",
						"Snippet for declare dependencies",
						"depends_on:",
						"  - task_key: \"task_name\"",
						"run_if: \"ALL_SUCCESS\"",
					),
				},
			},
			{
				insertText: "depends_on",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	"run_if": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
				"`run_if` string",
				"Enum: `ALL_SUCCESS` | `ALL_DONE` | `NONE_FAILED` | `AT_LEAST_ONE_SUCCESS` | `ALL_FAILED` | `AT_LEAST_ONE_FAILED`",
				"Default `\"ALL_SUCCESS\"`",
				"An optional value specifying the condition determining whether the task is run once its dependencies have been completed.",
				"  `ALL_SUCCESS`: All dependencies have executed and succeeded",
				"  `AT_LEAST_ONE_SUCCESS`: At least one dependency has succeeded",
				"  `NONE_FAILED`: None of the dependencies have failed and at least one was executed",
				"  `ALL_DONE`: All dependencies have been completed",
				"  `AT_LEAST_ONE_FAILED`: At least one dependency failed",
				"  `ALL_FAILED`: All dependencies have failed",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#tasks-run_if)",
			),
		},
	},
	"spark_python_task": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n",
				"`spark_python_task` object",
				"If spark_python_task, indicates that this task must run a Python file.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#tasks-spark_python_task)",
			),
		},
	},
	"python_wheel_task": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n",
				"`python_wheel_task` object",
				"If python_wheel_task, indicates that this job must execute a PythonWheel.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#tasks-python_wheel_task)",
			),
		},
	},
	"job_clusters": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n",
				"`job_clusters` Array of object <= 100 items",
				"A list of job cluster specifications that can be shared and reused by tasks of this job. Libraries cannot be declared in a shared job cluster. You must declare dependent libraries in task settings.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#job_clusters)",
			),
		},
	},
	"job_cluster_key": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n",
				"`job_cluster_key` [ 1 .. 100 ] characters ^[\\w\\-\\_]+$object <= 100 items",
				"Example `\"auto_scaling_cluster\"`",
				"A unique name for the job cluster. This field is required and must be unique within the job. `JobTaskSettings` may refer to this field to determine which cluster to launch for the task execution.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#job_clusters-job_cluster_key)",
			),
		},
	},
	"new_cluster": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n",
				"`new_cluster` object",
				"If new_cluster, a description of a cluster that is created for each task.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#job_clusters-new_cluster)",
			),
		},
	},
	"existing_cluster_id": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n",
				"`existing_cluster_id` string",
				"Example `\"0923-164208-meows279\"`",
				"If existing_cluster_id, the ID of an existing cluster that is used for all runs. When running jobs or tasks on an existing cluster, you may need to manually restart the cluster if it stops responding. We suggest running jobs and tasks on new clusters for greater reliability",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#tasks-existing_cluster_id)",
			),
		},
	},
	"cluster": {
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n",
					"job_clusters:",
					"  - job_cluster_key: \"job_cluster\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "job clusters chunk (name only)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n",
						"Snippet for declaring job clusters chunk",
						"job_clusters:",
						"  - job_cluster_key: \"job_cluster\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
					"job_clusters:",
					"  - job_cluster_key: \"job_cluster\"",
					"    new_cluster:",
					"      autoscale:",
					"        min_workers: 5",
					"        max_workers: 15",
					"      spark_conf:",
					"        spark.sql.shuffle.partitions: \"auto\"",
					"      runtime_engine: \"PHOTON\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "job clusters chunk (new, suggested settings)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
						"Snippet for declaring new job clusters chunk",
						"job_clusters:",
						"  - job_cluster_key: \"job_cluster\"",
						"    new_cluster:",
						"      autoscale:",
						"        min_workers: 5",
						"        max_workers: 15",
						"      spark_conf:",
						"        spark.sql.shuffle.partitions: \"auto\"",
						"      runtime_engine: \"PHOTON\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n", "job_cluster_key: \"job_cluster\""),
				kind:       lsp.CompletionItemKind["Snippet"],
				detail:     "job cluster key (name only)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n",
						"Snippet for declare job cluster key",
						"job_cluster_key: \"job_cluster\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
					"job_cluster_key: \"job_cluster\"",
					"new_cluster:",
					"  autoscale:",
					"    min_workers: 5",
					"    max_workers: 15",
					"  spark_conf:",
					"    spark.sql.shuffle.partitions: \"auto\"",
					"  runtime_engine: \"PHOTON\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "job cluster (new, suggested settings)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
						"Snippet for declaring new job cluster",
						"job_cluster_key: \"job_cluster\"",
						"new_cluster:",
						"  autoscale:",
						"    min_workers: 5",
						"    max_workers: 15",
						"  spark_conf:",
						"    spark.sql.shuffle.partitions: \"auto\"",
						"  runtime_engine: \"PHOTON\"",
					),
				},
			},
			// "existing_cluster_id": {},
		},
	},
	"docker_image": {
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n",
					"docker_image:",
					"  url: \"ecs-url\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "docker image url",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n",
						"Snippet for declare docker image url",
						"docker_image:",
						"  url: \"ecs-url\"",
					),
				},
			},
			{
				insertText: "docker_image",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	// "git_source": {},
	"tags": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
				"`tags` object",
				"Default `{}`",
				"Example",
				"```yaml",
				"tags:",
				"  cost-center: \"engineering\"",
				"  team: \"jobs\"",
				"```",
				"A map of tags associated with the job. These are forwarded to the cluster as cluster tags for jobs clusters, and are subject to the same limitations as cluster tags. A maximum of 25 tags can be added to the job.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#tags)",
			),
		},
		diag: Diag{
			help:     "`tags` is missing. Hint: start by typing `tags`",
			severity: 2,
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n",
					"tags:",
					"  tag-key: \"tag-value\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "workflow tags",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n",
						"Snippet for declare workflow tags",
						"tags:",
						"  tag-key: \"tag-value\"",
					),
				},
			},
			{
				insertText: "tags",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	// DEPRECATED: "format": {},
	// "queue": {},
	"parameters": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
				"`parameters` Array of object",
				"Job-level parameter definitions",
				"Example",
				"```yaml",
				"parameters:",
				"  - \"table=users\"",
				"  - \"date\"",
				"  - \"{{start_date}}\"",
				"```",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#parameters)",
			),
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n",
					"parameters:",
					"  - \"param-key\"",
					"  - \"param-value\"",
					"  - \"param-key=param-value\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "parameters chunk",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n",
						"Snippet for declare a chunk of parameters",
						"parameters:",
						"  - \"param-key\"",
						"  - \"param-value\"",
						"  - \"param-key=param-value\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n",
					"- \"param-key\"",
					"- \"param-value\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "single parameter (array style)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n",
						"Snippet for declare a single parameters",
						"- \"param-key\"",
						"- \"param-value\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n", "- \"param-key=param-value\""),
				kind:       lsp.CompletionItemKind["Snippet"],
				detail:     "single parameter (equal style)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n",
						"Snippet for declare a single parameters",
						"- \"param-key=param-value\"",
					),
				},
			},
			{
				insertText: "parameters",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	"run_as": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
				"`run_as` object",
				"Write-only setting, available only in Create/Update/Reset and Submit calls. Specifies the user or service principal that the job runs as. If not specified, the job runs as the user who created the job.",
				"Only `user_name` or `service_principal_name` can be specified. If both are specified, an error is thrown.",
				"Example",
				"```yaml",
				"run_as:",
				"  user_name: \"some.one@some.org\"",
				"# or",
				"  service_principle_name: \"some_service_principle\"",
				"```",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#run_as)",
			),
		},
		diag: Diag{
			help:     "`run_as` is missing. Hint: start by typing `run_as`",
			severity: 1,
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n",
					"run_as:",
					"  user_name: \"some.one@some.org\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "run as (user name)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n",
						"Snippet for declare run as",
						"run_as:",
						"  user_name: \"some.one@some.org\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n",
					"run_as:",
					"  service_principle_name: \"some_service_principle\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "run as (service principle)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n",
						"Snippet for declare run as",
						"run_as:",
						"  service_principle_name: \"some_service_principle\"",
					),
				},
			},
			{
				insertText: "run_as",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	"user_name": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n",
				"`user_name` string",
				"Example `\"some.one@some.org\"`",
				"The email of an active workspace user. Non-admin users can only set this field to their own email.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#run_as-user_name)",
			),
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n", "user_name: \"some.one@some.org\""),
				kind:       lsp.CompletionItemKind["Snippet"],
				detail:     "user name",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n",
						"Snippet for declare user name",
						"user_name: \"some.one@some.org\"",
					),
				},
			},
			{
				insertText: "user_name",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	"service_principle_name": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n",
				"`service_principle_name` string",
				"Example `\"some_service_principle\"`",
				"Application ID of an active service principal. Setting this field requires the `servicePrincipal/user` role.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#run_as-service_principle_name)",
			),
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n", "service_principle_name: \"some_service_principle\""),
				kind:       lsp.CompletionItemKind["Snippet"],
				detail:     "service principle",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n",
						"Snippet for declare service principle name",
						"service_principle_name: \"some_service_principle\"",
					),
				},
			},
			{
				insertText: "service_principle_name",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	"edit_mode": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n",
				"`edit_mode` string",
				"Enum: `UI_LOCKED | EDITABLE`",
				"Edit mode of the job.",
				"  `UI_LOCKED`: The job is in a locked UI state and cannot be modified.",
				"  `EDITABLE`: The job is in an editable state and can be modified.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#edit_mode)",
			),
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n", "edit_mode: \"EDITABLE\""),
				kind:       lsp.CompletionItemKind["Snippet"],
				detail:     "edit mode (editable on UI)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n",
						"Snippet for declare edit mode",
						"edit_mode: \"EDITABLE\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n", "edit_mode: \"UI_LOCKED\""),
				kind:       lsp.CompletionItemKind["Snippet"],
				detail:     "edit mode (locked on UI)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n",
						"Snippet for declare edit mode",
						"edit_mode: \"UI_LOCKED\"",
					),
				},
			},
			{
				insertText: "edit_mode",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	// "deployment": {},
	// "environments": {},
	"access_control_list": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n",
				"`access_control_list` Array of object",
				"List of permissions to set on the job.",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#access_control_list)",
			),
		},
		diag: Diag{
			help:     "`access_control_list` is missing. Hint: start by typing `access`",
			severity: 1,
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
					"access_control_list:",
					"  - user_name: \"some.one@some.org\"",
					"    permission_level: \"IS_OWNER\"",
					"  - group_name: \"developer\"",
					"    permission_level: \"CAN_MANAGE\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "access controll list (recommended)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n%s\n",
						"Snippet for declare access control list",
						"access_control_list:",
						"  - user_name: \"some.one@some.org\"",
						"    permission_level: \"IS_OWNER\"",
						"  - group_name: \"developer\"",
						"    permission_level: \"CAN_MANAGE\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n",
					"- user_name: \"some.one@some.org\"",
					"  permission_level: \"IS_OWNER\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "access controll (user name)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n",
						"Snippet for declare user name access control",
						"- user_name: \"some.one@some.org\"",
						"  permission_level: \"IS_OWNER\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n",
					"- group_name: \"developer\"",
					"  permission_level: \"CAN_MANAGE\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "access controll (group name)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n",
						"Snippet for declare group name access control",
						"- group_name: \"developer\"",
						"  permission_level: \"CAN_MANAGE\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n",
					"- service_principle_name: \"some_service_principle\"",
					"  permission_level: \"CAN_RUN\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "access controll (service principle name)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n",
						"Snippet for declare service principle name access control",
						"- service_principle_name: \"some_service_principle\"",
						"  permission_level: \"CAN_RUN\"",
					),
				},
			},
			{
				insertText: "access_control_list",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	"permission_level": {
		hover: lsp.MarkupContent{
			Kind: "markdown",
			Value: fmt.Sprintf("%s\n%s\n%s\n",
				"`permission_level` string",
				"Enum: `Enum: CAN_MANAGE | CAN_RESTART | CAN_ATTACH_TO | IS_OWNER | CAN_MANAGE_RUN | CAN_VIEW | CAN_READ | CAN_RUN | CAN_EDIT | CAN_USE | CAN_MANAGE_STAGING_VERSIONS | CAN_MANAGE_PRODUCTION_VERSIONS | CAN_EDIT_METADATA | CAN_VIEW_METADATA | CAN_BIND | CAN_QUERY`",
				"[See more](https://docs.databricks.com/api/workspace/jobs/create#access_control_list-permission_level)",
			),
		},
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n", "permission_level: \"CAN_RUN\""),
				kind:       lsp.CompletionItemKind["Snippet"],
				detail:     "permisison can run",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n",
						"Snippet for declare can run permission",
						"permisison_level: \"CAN_RUN\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n", "permission_level: \"CAN_MANAGE\""),
				kind:       lsp.CompletionItemKind["Snippet"],
				detail:     "permisison can manage",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n",
						"Snippet for declare can manage permission",
						"permisison_level: \"CAN_MANAGE\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n", "permission_level: \"CAN_VIEW\""),
				kind:       lsp.CompletionItemKind["Snippet"],
				detail:     "permisison can view",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n",
						"Snippet for declare can view permission",
						"permisison_level: \"CAN_VIEW\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n", "permission_level: \"IS_OWNER\""),
				kind:       lsp.CompletionItemKind["Snippet"],
				detail:     "permisison owner",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n",
						"Snippet for declare owner permission",
						"permisison_level: \"IS_OWNER\"",
					),
				},
			},
			{
				insertText: "permission_level",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
}
