## 1.6.0 (TBD)

FEATURES:

* **New Data Source**: `splunkconfig_app_package`
* **New Resource**: `splunkconfig_app_auto_version`
* **Deprecated Resource**: `splunkconfig_app_package`

## 1.5.0 (October 26, 2021)

FEATURES:

* **Data Source Enhancement**: `splunkconfig_lookup_attributes` implements `row_number_field`
* **Data Source Enhancement**: `splunkconfig_app_ids` implements `exclude_tag`

## 1.4.0 (October 19, 2021)

FEATURES:

* **Provider Change**: New configuration argument `configuration_path`.
* **Schema Change**: `Apps` can have `collections`.
* **New Tool**: `template-lookup-csv`

## 1.3.2 (September 28, 2021)

FIXES:

* **Schema Change**: Index names can have asterisks (*) for `srchIndexesAllowed`.

## 1.3.1 (September 27, 2021)

FIXES:

* **Schema Change**: `App` IDs can have dashes

## 1.3.0 (September 21, 2021)

FEATURES:

* **Schema Change**: `Apps` can have tags
* **Data Source Enhancement**: `splunkconfig_app_ids` implements `require_tag`

## 1.2.0 (September 3, 2021)

FEATURES:

* **Resource Enhancement**: `Lookups` generate transforms.conf when packaged with `splunkconfig_app_package`
* **New Data Source:** `app_ids`
* **New Data Source:** `app_attribtes`
* **Schema Change**: `Apps` can have ACLs

## 1.1.0 (August 23, 2021)

FEATURES:

* **New Data Source:** `splunkconfig_lookup_attributes`

## 1.0.0 (August 11, 2021)

Initial version
