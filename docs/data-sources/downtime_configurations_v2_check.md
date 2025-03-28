---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "synthetics_downtime_configurations_v2_check Data Source - synthetics"
subcategory: ""
description: |-
  
---

# synthetics_downtime_configurations_v2_check (Data Source)



## Example Usage

```terraform
data "synthetics_downtime_configurations_v2_check" "datasource_locations" {
  downtime_configurations {
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `downtime_configurations` (Block Set) (see [below for nested schema](#nestedblock--downtime_configurations))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--downtime_configurations"></a>
### Nested Schema for `downtime_configurations`

Required:

- `id` (Number)

Optional:

- `description` (String)
- `recurrence` (Block Set) (see [below for nested schema](#nestedblock--downtime_configurations--recurrence))
- `timezone` (String)

Read-Only:

- `created_at` (String)
- `end_time` (String)
- `name` (String)
- `rule` (String)
- `start_time` (String)
- `status` (String)
- `test_count` (Number)
- `tests_updated_at` (String)
- `updated_at` (String)

<a id="nestedblock--downtime_configurations--recurrence"></a>
### Nested Schema for `downtime_configurations.recurrence`

Optional:

- `end` (Block Set) (see [below for nested schema](#nestedblock--downtime_configurations--recurrence--end))
- `repeats` (Block Set) (see [below for nested schema](#nestedblock--downtime_configurations--recurrence--repeats))

<a id="nestedblock--downtime_configurations--recurrence--end"></a>
### Nested Schema for `downtime_configurations.recurrence.end`

Optional:

- `type` (String)
- `value` (String)


<a id="nestedblock--downtime_configurations--recurrence--repeats"></a>
### Nested Schema for `downtime_configurations.recurrence.repeats`

Optional:

- `custom_frequency` (String)
- `custom_value` (Number)
- `type` (String)
