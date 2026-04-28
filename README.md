# annon

`annon` is a Go library for anonymising sensitive data before logging,
exporting, serialising, or sending it to telemetry.

The repository now exposes two public packages only:

- `github.com/BreakPointSoftware/annon`
- `github.com/BreakPointSoftware/annon/redact`

Everything else lives under `internal/...`.

## Public Packages

### `annon`

Use `annon` to walk structured data, apply field and value detection, and
return anonymised copies or serialised output.

```go
safeCustomer, err := annon.Copy(customer)

jsonBlob, err := annon.JSON(customer)

yamlBlob, err := annon.YAML(customer)

safeJSON, err := annon.FromJSON(rawJSON)

safeYAML, err := annon.FromYAML(rawYAML)
```

Reusable instance:

```go
a, err := annon.New(
    annon.WithValueDetection(true),
)

safeAny, err := a.Copy(customer)
jsonBlob, err := a.JSON(customer)
```

### `redact`

Use `redact` for direct single-value redaction without structured walking.

```go
safeEmail := redact.Email("greg@example.com")
safePhone := redact.Phone("07700 900123")
safePostcode := redact.Postcode("TN9 1XA")
safeReg := redact.VehicleRegistration("AB12 CDE")
safeName := redact.Name("Greg Bryant")
safeText := redact.Redact("secret")
```

Reusable configured redactor:

```go
r, err := redact.New(
    redact.WithRedactChar('x'),
)

safeEmail := r.Email("greg@example.com")
```

## Public Configuration

Preservation configuration is public through the root `annon` package.

```go
cfg := annon.PreservationConfig{
    RedactionText: "[hidden]",
    RedactChar:    'x',
    Email: annon.EmailConfig{
        KeepLocalPrefix: 2,
        KeepDomain:      true,
    },
    Phone: annon.PhoneConfig{
        KeepLast: 3,
    },
    Name: annon.NameConfig{
        KeepPrefix: 1,
    },
    Postcode: annon.PostcodeConfig{
        KeepOutward: true,
    },
    VehicleRegistration: annon.VehicleRegistrationConfig{
        KeepPrefix: 2,
    },
}
```

Apply it to structured anonymisation with:

- `annon.WithPreservation(cfg)`

Apply equivalent redaction settings to direct redaction with:

- `redact.WithConfig(redact.Config(cfg))`

## Detection Order

Structured anonymisation follows this order:

1. `anonymise:"false"`
2. explicit tag strategy
3. strong field-name match
4. fallback field-name match
5. contains field-name match
6. value-pattern detection
7. no anonymisation

## Supported Tags

- `anonymise:"false"`
- `anonymise:"true"`
- `anonymise:"auto"`
- `anonymise:"email"`
- `anonymise:"phone"`
- `anonymise:"postcode"`
- `anonymise:"name"`
- `anonymise:"firstName"`
- `anonymise:"surname"`
- `anonymise:"vehicleRegistration"`
- `anonymise:"redact"`
- `anonymise:"remove"`

Unknown explicit tag strategy names are treated as errors.

## Field Detection

Default field detection includes:

- strong email matches such as `email` and `emailAddress`
- strong phone matches such as `phoneNumber` and `mobileNumber`
- strong postcode matches such as `postcode`
- strong name-part matches such as `firstName` and `surname`
- strong vehicle registration matches such as `vehicleRegistration`, `vehicleReg`, and `vrm`
- fallback matches such as `reg` and `phone`
- contains matches such as `customerName`
- built-in exclusions for `username`, `fileName`, `hostName`, and `domainName`

Value detection is opt-in and checks strings for:

- email addresses
- UK mobile numbers
- UK postcodes
- UK vehicle registrations

Additional field rules can be added with:

- `annon.WithFieldRules(...)`

## Output Behaviour

- `Copy` returns the same concrete Go type and does not mutate the input
- `JSON` and `YAML` walk Go values into neutral map/list structures before serialising
- `FromJSON` and `FromYAML` anonymise raw blobs without mutating the input bytes
- `anonymise:"remove"` zeroes a field in copy mode and omits it from JSON/YAML output

## Internal Structure

The repo is organised around two public packages and internal domain packages:

```text
annon/
├── *.go
├── redact/
└── internal/
    ├── detection/
    ├── encode/
    ├── redactcore/
    ├── support/
    │   └── normalise/
    └── walk/
```

### Internal responsibilities

- `internal/detection`: field and value matching
- `internal/encode`: JSON and YAML decode/encode helpers
- `internal/redactcore`: concrete redaction implementations and shared config
- `internal/support/normalise`: low-level field-name normalisation helpers
- `internal/walk`: typed and blob traversal, tag handling, and metadata caching

## Testing

The repository uses same-package tests only.

- public integration tests live in `anonymise` and `redact`
- internal implementation tests live in the owning `internal/...` packages
- benchmark coverage exists for detection, walking, and public anonymisation entrypoints

Run the full suite with:

```bash
go test ./...
```

Run benchmarks with:

```bash
go test -bench=. ./...
```
