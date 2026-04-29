# annon

`annon` is a Go library for redacting sensitive data before logging,
exporting, serialising, or sending it to telemetry.

The primary public package is:

- `github.com/BreakPointSoftware/annon/redact`

Everything else lives under `internal/...` and is not part of the public
compatibility contract.

## Stability

This library should be treated as pre-`v1`.

- the public API may still change
- internal packages are not part of the compatibility contract
- breaking changes may happen between releases while the design settles

## Public API

Use `redact` as the public entrypoint for defensive redaction.

```go
safeValue := redact.Data(customer)

safeJSON := redact.JSON(customer)
safeYAML := redact.YAML(customer)

safeJSONBytes := redact.JSONBytes(rawJSON)
safeYAMLBytes := redact.YAMLBytes(rawYAML)
```

These APIs are intended for logging and export use cases:

- they do not return errors
- they are designed not to panic
- JSON/YAML helpers always return valid fallback payloads on failure
- caller input must not be mutated

## Direct Value Redactors

Use the direct string helpers when you already know the value type.

```go
safeString := redact.String("greg@example.com")
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

Preservation configuration is public through the `redact` package.

```go
cfg := redact.Config{
    RedactionText: "[hidden]",
    RedactChar:    'x',
    Email: redact.EmailConfig{
        KeepLocalPrefix: 2,
        KeepDomain:      true,
    },
    Phone: redact.PhoneConfig{
        KeepLast: 3,
    },
    Name: redact.NameConfig{
        KeepPrefix: 1,
    },
    Postcode: redact.PostcodeConfig{
        KeepOutward: true,
    },
    VehicleRegistration: redact.VehicleRegistrationConfig{
        KeepPrefix: 2,
    },
}
```

Apply it with:

- `redact.WithConfig(cfg)`

## Detection Order

Structured redaction follows this order:

1. `anonymise:"false"`
2. explicit tag strategy
3. strong field-name match
4. fallback field-name match
5. contains field-name match
6. value-pattern detection
7. no redaction

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

Unknown explicit tag strategy names are treated as errors internally.

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

Value detection is opt-in internally and checks strings for:

- email addresses
- UK mobile numbers
- UK postcodes
- UK vehicle registrations

## Fallback Behaviour

- malformed JSON bytes return `{"redaction_error":true}`
- malformed YAML bytes return `redaction_error: true`
- unrecoverable value-level failures fall back conservatively rather than leaking the original value

## Internal Structure

```text
annon/
├── redact/
└── internal/
	    ├── walk/
	    ├── decision/
	    ├── detection/
	    ├── encode/
    ├── engine/
    ├── output/
    ├── redactcore/
    └── support/
        └── normalise/
```

### Internal responsibilities

- `internal/walk`: typed traversal, deep-copy logic, and struct metadata caching
- `internal/decision`: tag parsing and redaction decision flow
- `internal/detection`: field and value matching
- `internal/encode`: JSON and YAML decode/encode helpers
- `internal/engine`: no-error orchestration and fallback handling for public `redact` APIs
- `internal/output`: neutral output construction for JSON/YAML and raw blobs
- `internal/redactcore`: concrete redaction implementations and shared config
- `internal/support/normalise`: low-level field-name normalisation helpers

## Testing

The repository uses same-package tests only.

- public integration tests live in `redact`
- internal implementation tests live in the owning `internal/...` packages
- benchmark coverage exists for detection, walk, output, and public redaction entrypoints

Run the full suite with:

```bash
go test ./...
```

Run benchmarks with:

```bash
go test -bench=. ./...
```
