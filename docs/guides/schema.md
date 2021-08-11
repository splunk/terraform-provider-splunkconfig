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
- **apps** (List of Object) Apps defined. (see [schema for app](#app))
- **indexes** (List of Object) Indexes defined. (see [schema for index](#index))
- **lookups** (List of Object) Lookups defined. (see [schema for lookup](#lookup))
- **roles** (List of Object) Roles defined. (see [schema for role](#role))
- **saml_groups** (List of Object) SAML Groups defined. (see [schema for saml_group](#saml_group))
- **users** (List of Object) Users defined. (see [schema for user](#user))

<a id="app"></a>
## Schema for `app`

- **name** (String, required) App name.
- **description** (String) App description.
- **id** (String, required) App ID. As per the `app.conf` specification:
```
* id must adhere to these cross-platform folder name restrictions:
* must contain only letters, numbers, "." (dot), and "_" (underscore) characters.
```
- **author** (String) App author.
- **isvisible** (Bool) App visibility.
- **version** (String or Object) App version. If given as a string, must be in `<major>.<minor>.<patch>` format.
Can also be a `version` object. If not defined, defaults to `0.0.0`. (see [schema for version](#version))
- **indexes** (Bool or List of Object) If `true`, include the global `indexes` configuration in this app. Can also
be a list of index objects to include in the app. (see [schema for index](#index))
- **lookups** (List of String or List of Object) If defined as a list of strings, include the referenced global
`lookup` objects in this app. Can also be a list of lookup objects to include in the app. (see [schema for lookup](#lookup))
- **roles** (Bool or List of Object) If `true`, include the global `roles` configuration in this app. Can also be a
list of role objects to include in the app. (see [schema for role](#role))

<a id="index"></a>
## Schema for `index`

- **name** (String, required) Index name. As per the `indexes.conf` specification:
```
Index names must consist of only numbers, lowercase letters, underscores,
and hyphens. They cannot begin with an underscore or hyphen, or contain
the word "kvstore".
```
- **frozenTimePeriod** (Object) Frozen time period. (see [schema for timeperiod](#timeperiod))
- **srchRolesAllowed** (List of String) Names of roles that can search this index. List values must be valid role
names. As per the `authorize.conf` specification:
```
* Role names cannot have uppercase characters.
* Role names cannot contain spaces, colons, semicolons, or forward slashes.
```
- **lookup_rows** (List of Object) Add rows of values for this index to lookups. (see [schema for lookup_row](#lookup_row))
- **homePath** (String) homePath of the index. Defaults to `$SPLUNK_DB/<index name>/db`.
- **coldPath** (String) coldPath of the index. Defaults to `$SPLUNK_DB/<index name>/colddb`.
- **thawedPath** (String) thawedPath of the index. Defaults to `$SPLUNK_DB/<index name>/thaweddb`.

<a id="lookup"></a>
## Schema for `lookup`

- **name** (String, required) Lookup name. The resulting CSV file will be `<name>.csv`.
- **fields** (List of Object) Fields included in the lookup. (see [schema for lookup_field](#lookup_field))
- **rows** (List of Object) Rows included in the lookup. (see [schema for lookup_row](#lookup_row))

<a id="lookup_field"></a>
## Schema for `lookup_field`

- **name** (String, required) Name of the field, placed in the header row.
- **default** (String) Default value for this field.
- **default_row_field** (Bool) If true, an object type that sets a value for this field will result in a default row
being automatically created for the associated lookup. This is useful when you want to automatically create rows for
every index or role.
- **required** (Bool) If true, a value for this field must exist for every row, or the lookup will fail validation.

<a id="lookup_row"></a>
## Schema for `lookup_row`

- **lookup_name** (String) Name of lookup the row belongs to. Not used when defined directly in a lookup object.
- **values** (Map) Lookup values to create. Map keys must match fields that exist in the lookup.

<a id="role"></a>
## Schema for `role`

- **name** (String, required) Name of the role.  As per the `authorize.conf` specification:
```
* Role names cannot have uppercase characters.
* Role names cannot contain spaces, colons, semicolons, or forward slashes.
```
- **saml_groups** (List of String) List of SAML groups that this role should be added to.
- **srchIndexesAllowed** (List of String) List of indexes that this role can search. Listed index names must be
valid.  As per the `indexes.conf` specification:
```
Index names must consist of only numbers, lowercase letters, underscores,
and hyphens. They cannot begin with an underscore or hyphen, or contain
the word "kvstore".
```
- **importRoles** (List of String) List of Roles that this role will import. Listed role names must be valid. As
per the `authorize.conf` specification:
```
* Role names cannot have uppercase characters.
* Role names cannot contain spaces, colons, semicolons, or forward slashes.
```
- **capabilities** (Map of String to Bool) Capabilities defined for this role. Key (String) is the capability name,
value (Bool) is true if the capability is enabled, false if the capability is disabled. Capability names must be
valid. As per the `authorize.conf` specification:
```
Only alphanumeric characters and "_" (underscore) are allowed in capability names.
```
- **lookup_rows** (List of Object) Lookup rows to create for this role. (see [schema for lookup_row](#lookup_row))
- **srchFilter** (String) Search filter defined for the role.
- **srchTimeWin** (Integer) srchTimeWin for the role.
- **srchDiskQuota** (Integer) srchDiskQuota for the role.
- **srchJobsQuota** (Integer) srchJobsQuota for the role.
- **rtSrchJobsQuota** (Integer) rtSrchJobsQuota for the role.
- **cumulativeSrchJobsQuota** (Integer) cumulativeSrchJobsQuota for the role.
- **cumulativeRTSrchJobsQuota** (Integer) cumulativeRTSrchJobsQuota for the role.

<a id="saml_group"></a>
## Schema for `saml_group`

- **name** (String) Name of the SAML group.
- **roles** (List of String) Roles to apply to the SAML group. Listed role names must be valid. As per the
`authorize.conf` specification:
```
* Role names cannot have uppercase characters.
* Role names cannot contain spaces, colons, semicolons, or forward slashes.
```

<a id="timeperiod"></a>
## Schema for `timeperiod`

- **seconds** (Integer) Seconds.
- **minutes** (Integer) Minutes.
- **hours** (Integer) Hours.
- **days** (Integer) Days.

<a id="user"></a>
## Schema for `user`

- **name** (String) Name of the user (username used to log in).
- **email** (String) Email address of the user.
- **password** (String, sensitive) Password of the user. Avoid using this for any real password value. It can
instead be used to trigger rotation of a randomly generated password whenever this value changes.
- **force_change_pass** (Bool) True if the user should be forced to change their password after logging in.
- **realname** (String) Real name of the user.
- **roles** (List of String) Roles to apply to the user.

<a id="version"></a>
## Schema for `version`

- **major** (Integer) Major version.
- **minor** (Integer) Minor version.
- **patch** (Integer) Patch version.
