## 12. Database design

### `states`

```text
id
code
name
region
division
created_at
updated_at
```

### `categories`

```text
id
slug
name
description
default_weight
display_order
```

### `metrics`

```text
id
category_id
slug
name
description
unit
higher_is_better
normalization_method
default_weight
source_id
active
created_at
updated_at
```

### `metric_values`

```text
id
state_id
metric_id
year
value
source_record_id
import_id
created_at
```

Unique constraint:

```text
state_id + metric_id + year + import_id
```

### `data_sources`

```text
id
name
publisher
source_url
license
format
description
created_at
updated_at
```

### `imports`

```text
id
source_id
status
started_at
completed_at
records_read
records_inserted
records_rejected
checksum
error_summary
```

### `import_errors`

```text
id
import_id
row_number
field_name
raw_value
error_message
```

### `scoring_profiles`

```text
id
name
description
is_default
is_system
created_at
updated_at
```

### `profile_category_weights`

```text
profile_id
category_id
weight
```

### `profile_metric_weights`

```text
profile_id
metric_id
weight
```

### `score_snapshots`

```text
id
profile_id
state_id
year
overall_score
completeness
calculated_at
calculation_version
```

### `category_score_snapshots`

```text
score_snapshot_id
category_id
score
completeness
```

### `application_settings`

```text
key
value
updated_at
```

---
