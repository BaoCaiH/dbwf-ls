package analysis

import (
	"dbwf-ls/lsp"
	"fmt"
	"strings"
)

//	lsp.CompletionItem{
//			Label:         word,
//			Detail:        "Current typing word",
//			Documentation: "Nothing to document here",
//		}
type Keyword struct {
	completions []Completion
}

type Completion struct {
	insertText, detail string
	kind               int
	documentation      lsp.MarkupContent
}

var Keywords = map[string]Keyword{
	"name": {
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
	"email_notification": {
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
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n",
					"health:",
					"  rules:",
					"     -",
					"        metric: \"RUN_DURATION_SECONDS\"",
					"        op: \"GREATER_THAN\"",
					"        value: 10800 # 3 hours",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "workflow health",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n%s\n%s\n",
						"Snippet for setting workflow health alert",
						"health:",
						"  rules:",
						"     -",
						"        metric: \"RUN_DURATION_SECONDS\"",
						"        op: \"GREATER_THAN\"",
						"        value: 10800 # 3 hours",
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
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n",
					"schedule:",
					"    quartz_cron_expression: \"0 0 0 * * ?\" # Everyday at 0am",
					"    timezone_id: \"UTC\"",
					"    pause_status: \"UNPAUSED\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "schedule (unpaused)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n",
						"Snippet to schedule workflow (unpaused)",
						"schedule:",
						"    quartz_cron_expression: \"0 0 0 * * ?\" # Everyday at 0am",
						"    timezone_id: \"UTC\"",
						"    pause_status: \"UNPAUSED\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n",
					"schedule:",
					"    quartz_cron_expression: \"0 0 0 * * ?\" # Everyday at 0am",
					"    timezone_id: \"UTC\"",
					"    pause_status: \"PAUSED\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "schedule (paused)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n",
						"Snippet to schedule workflow (paused)",
						"schedule:",
						"    quartz_cron_expression: \"0 0 0 * * ?\" # Everyday at 0am",
						"    timezone_id: \"UTC\"",
						"    pause_status: \"PAUSED\"",
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
	// TODO "tasks": {},
	"job_clusters": {
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n",
					"job_clusters:",
					"  -",
					"    job_cluster_key: \"job_cluster\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "job clusters (name only)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n",
						"Snippet for declaring job clusters",
						"job_clusters:",
						"  -",
						"    job_cluster_key: \"job_cluster\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
					"job_clusters:",
					"  -",
					"    job_cluster_key: \"job_cluster\"",
					"    new_cluster:",
					"      autoscale:",
					"        min_workers: 5",
					"        max_workers: 15",
					"      spark_conf:",
					"        spark.sql.shuffle.partitions: \"auto\"",
					"      runtime_engine: \"PHOTON\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "job clusters (new, suggested settings)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
						"Snippet for declaring new job clusters",
						"job_clusters:",
						"  -",
						"    job_cluster_key: \"job_cluster\"",
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
				insertText: "job_clusters",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	"job_cluster_key": {
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n", "job_cluster_key: \"job_cluster\""),
				kind:       lsp.CompletionItemKind["Snippet"],
				detail:     "job cluster key",
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
				detail: "job clusters key (new, suggested settings)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
						"Snippet for declaring new job clusters",
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
			{
				insertText: "job_cluster_key",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
		},
	},
	"new_cluster": {
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
					"new_cluster:",
					"  autoscale:",
					"    min_workers: 5",
					"    max_workers: 15",
					"  spark_conf:",
					"    spark.sql.shuffle.partitions: \"auto\"",
					"  runtime_engine: \"PHOTON\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "new job cluster (suggested settings)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
						"Snippet for declare new job cluster",
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
			{
				insertText: "new_cluster",
				kind:       lsp.CompletionItemKind["Keyword"],
			},
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
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n",
					"tags:",
					"    tag-key: \"tag-value\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "workflow tags",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n",
						"Snippet for declare workflow tags",
						"tags:",
						"    tag-key: \"tag-value\"",
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
		completions: []Completion{
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
					"access_control_list:",
					"  -",
					"    user_name: \"some.one@some.org\"",
					"    permission_level: \"IS_OWNER\"",
					"  -",
					"    group_name: \"developer\"",
					"    permission_level: \"CAN_MANAGE\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "access controll list (recommended)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
						"Snippet for declare access control list",
						"access_control_list:",
						"  -",
						"    user_name: \"some.one@some.org\"",
						"    permission_level: \"IS_OWNER\"",
						"  -",
						"    group_name: \"developer\"",
						"    permission_level: \"CAN_MANAGE\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n",
					"-",
					"  user_name: \"some.one@some.org\"",
					"  permission_level: \"IS_OWNER\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "access controll (user name)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n",
						"Snippet for declare user name access control",
						"-",
						"  user_name: \"some.one@some.org\"",
						"  permission_level: \"IS_OWNER\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n",
					"-",
					"  group_name: \"developer\"",
					"  permission_level: \"CAN_MANAGE\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "access controll (group name)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n",
						"Snippet for declare group name access control",
						"-",
						"  group_name: \"developer\"",
						"  permission_level: \"CAN_MANAGE\"",
					),
				},
			},
			{
				insertText: fmt.Sprintf("%s\n%s\n%s\n",
					"-",
					"  service_principle_name: \"some_service_principle\"",
					"  permission_level: \"CAN_RUN\"",
				),
				kind:   lsp.CompletionItemKind["Snippet"],
				detail: "access controll (service principle name)",
				documentation: lsp.MarkupContent{
					Kind: "markdown",
					Value: fmt.Sprintf("---\n%s\n\n%s\n%s\n%s\n",
						"Snippet for declare service principle name access control",
						"-",
						"  service_principle_name: \"some_service_principle\"",
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

func hammingRatio(input, keyword string) float32 {
	compareLength := len(input)
	if len(input) > len(keyword) {
		compareLength = len(keyword)
	}

	dist := 0
	for i := range compareLength {
		if input[i] != keyword[i] {
			dist += 1
		}
	}

	return (float32(compareLength) - float32(dist)) / float32(compareLength)
}

func complete(word, leading string) []lsp.CompletionItem {
	options := []lsp.CompletionItem{}
	for kw, completions := range Keywords {
		if hammingRatio(word, kw) >= 0.75 {
			for _, completion := range completions.completions {
				options = append(options, lsp.CompletionItem{
					Label:         kw,
					Kind:          completion.kind,
					Detail:        completion.detail,
					Documentation: completion.documentation,
					InsertText:    strings.ReplaceAll(completion.insertText, "\n", "\n"+leading),
				})
			}
		}
	}
	return options
}
