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

<a id="index"></a>
## Schema for `index` object

- **name** (String) Index name.
- **frozenTimePeriod** (Object) Frozen time period. (see [schema for timeperiod](#timeperiod))
- **srchRolesAllowed** (List of String) Names of roles that can search this index.
- **lookup_rows** (List of Object) Add rows of values for this index to lookups. (see [schema for lookup_row](#lookup_row))
- **homePath** (String) homePath of the index. Defaults to `$SPLUNK_DB/<index name>/db`.
- **coldPath** (String) coldPath of the index. Defaults to `$SPLUNK_DB/<index name>/colddb`.
- **thawedPath** (String) thawedPath of the index. Defaults to `$SPLUNK_DB/<index name>/thaweddb`.

<a id="timeperiod"></a>
## Schema for `timeperiod`

- **seconds** (Integer) Seconds.
- **minutes** (Integer) Minutes.
- **hours** (Integer) Hours.
- **days** (Integer) Days.

<a id="lookup_row"></a>
## Schema for `lookup_row`

- **lookup_name** (String) Name of lookup the row belongs to.
- **values** (Map) Lookup values to create.
