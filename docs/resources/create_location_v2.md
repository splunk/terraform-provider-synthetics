---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "synthetics_create_location_v2 Resource - synthetics"
subcategory: ""
description: |-
  
---

# synthetics_create_location_v2 (Resource)



## Example Usage

```terraform
resource "synthetics_create_location_v2" "location_v2_foo" {
  location {
    id = "private-aws-awesome-east"
    label = "awesome aws east location"
  }    
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `location` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--location))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--location"></a>
### Nested Schema for `location`

Required:

- `id` (String)
- `label` (String)

Read-Only:

- `country` (String)
- `default` (Boolean)
- `type` (String)
