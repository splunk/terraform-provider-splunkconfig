---
page_title: "splunkconfig YAML Schema"
subcategory: ""
description: |-
  Define your Splunk configuration
---

# YAML Schema

## Description

The Splunk configuration will be defined in YAML. In general, configuration options which exist in a Splunk .conf
specification will use the same field name as in that specification.

## Fields

- **anchors** (Freeform) Any valid YAML can be placed here for the purpose of defining YAML anchors.
- **indexes** (List of Object) Indexes defined. (see [schema for index](#index))
- **lookups** (List of Object) Lookups defined. (see [schema for lookup](#lookup))

<a id="index"></a>
## Schema for `index` object

- **name** (String) Index name.
- **frozenTimePeriod** (Object) Frozen time period. (see [schema for timeperiod](#timeperiod))
- **srchRolesAllowed** (List of String) Names of roles that can search this index.
- **lookup_rows** (List of Object) Add rows of values for this index to lookups. (see [schema for lookup_row](#lookup_row))
- **homePath** (String) homePath of the index. Defaults to `$SPLUNK_DB/<index name>/db`.
- **coldPath** (String) coldPath of the index. Defaults to `$SPLUNK_DB/<index name>/colddb`.
- **thawedPath** (String) thawedPath of the index. Defaults to `$SPLUNK_DB/<index name>/thaweddb`.

<a id="lookup"></a>
## Schema for `lookup` object

- **name** (String) Lookup name. The resulting CSV file will be `<name>.csv`.
- **fields** (List of Object) Fields included in the lookup. (see [schema for lookup_field](#lookup_field))
- **rows** (List of Object) Rows included in the lookup. (see [schema for lookup_row](#lookup_row))

<a id="lookup_field"></a>
## Schema for `lookup_field` object

- **name** (String) Name of the field, placed in the header row.
- **default** (String) Default value for this field.
- **default_row_field** (Bool) If true, an object type that sets a value for this field will result in a default row
being automatically created for the associated lookup. This is useful when you want to automatically create rows for
every index or role.
- **required** (Bool) If true, a value for this field must exist for every row, or the lookup will fail validation.

<a id="lookup_row"></a>
## Schema for `lookup_row` object

- **lookup_name** (String) Name of lookup the row belongs to. Not used when defined directly in a lookup object.
- **values** (Map) Lookup values to create.

<a id="timeperiod"></a>
## Schema for `timeperiod`

- **seconds** (Integer) Seconds.
- **minutes** (Integer) Minutes.
- **hours** (Integer) Hours.
- **days** (Integer) Days.

