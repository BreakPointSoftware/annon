# annon

`annon` is a Go library for anonymising sensitive structured data before logging,
exporting, serialising, or sending it to telemetry.

The public API lives in `github.com/BreakPointSoftware/annon/anonymise`.

## Principles

- Safety
- Predictability
- Extensibility
- Performance
- Clear package boundaries
- Non-mutating behaviour by default
- British English spelling throughout

## Public API

```go
safe, err := anonymise.Copy(customer)

jsonBlob, err := anonymise.JSON(customer)

yamlBlob, err := anonymise.YAML(customer)

safeJSON, err := anonymise.FromJSON(rawJSON)

safeYAML, err := anonymise.FromYAML(rawYAML)
```

All package-level functions accept options, and `anonymise.New(...)` returns a
reusable instance plus an error if configuration is invalid.

Reusable instance:

```go
a := anonymise.New(
    anonymise.WithDefaultStrategies(),
    anonymise.WithDefaultFieldDetection(),
    anonymise.WithValueDetection(true),
)

safeCustomer, err := a.Copy(customer)
safeJSON, err := a.JSON(customer)
```

## Decision Order

1. `anonymise:"false"`
2. Explicit tag strategy
3. Strong field-name match
4. Fallback field-name match
5. Contains field-name match
6. Value-pattern detection
7. No anonymisation

## Tag Behaviour

- `anonymise:"false"`: never anonymise the field
- `anonymise:"true"` and `anonymise:"auto"`: infer a strategy
- `anonymise:"email"`, `"phone"`, `"postcode"`, `"name"`, `"firstName"`, `"surname"`, `"vehicleRegistration"`: force strategy
- `anonymise:"redact"`: use configured redaction text
- `anonymise:"remove"`: zero in copy mode, omit in JSON/YAML mode

Unknown explicit tag strategy names are treated as configuration errors when a
field is walked.

## Default Detection

Default field-name detection includes:

- strong email matches such as `email` and `emailAddress`
- strong phone matches such as `phoneNumber` and `mobileNumber`
- strong postcode matches such as `postcode`
- strong name-part matches such as `firstName` and `surname`
- strong vehicle registration matches such as `vehicleRegistration`, `vehicleReg`, and `vrm`
- fallback matches such as `reg` and `phone`
- contains matches such as `customerName`, excluding false positives like `username`, `fileName`, `hostName`, and `domainName`

Value detection is opt-in and checks strings for:

- email addresses
- UK mobile numbers
- UK postcodes
- UK vehicle registrations

Additional field rules can be added with `anonymise.WithFieldRules(...)` when
the default compiled detector is used.

## Output Modes

- `Copy` returns the same concrete Go type and does not mutate the input
- `JSON` and `YAML` walk Go values into neutral map/list structures before serialising
- `FromJSON` and `FromYAML` anonymise raw blobs without mutating the input bytes

`anonymise:"remove"` behaviour:

- copy mode: zero value
- JSON/YAML mode: omitted field

## Testing

The repository uses same-package tests throughout. Internal implementation
details are tested directly in their owning packages rather than through
external `_test` packages.

## Behaviour Checklist

The implementation is intended to be built and tested in the following order:

1. Field normalisation
2. Strong field-name detection
3. Basic email strategy
4. Copy mode for simple structs
5. JSON mode for simple structs
6. Nested structs, maps, slices, arrays, pointers
7. Preservation config
8. Struct tags
9. Fallback and contains detection
10. Raw JSON input
11. Value detection
12. YAML support
13. Performance caches

Same-package tests cover internal implementation details in `detection`,
`strategy`, `walker`, `encoder`, and `internal` packages. Public API tests in
`anonymise` prove cross-package integration without relying on `_test`
packages.
