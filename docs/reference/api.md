# API Reference

The API is mounted at `/api/v1`. Success responses use `{ "data": ... }`; errors use `{ "error": { "code", "message", "details" } }`. The machine-readable contract is [spec/openapi.yaml](../../spec/openapi.yaml).

## Catalogs

- `GET /states?region=` and `GET /states/{code}`
- `GET /regions`, `/categories`, and `/metrics?category_id=`
- `GET /profiles`, `/profiles/{profileId}`, and `/profiles/default`

## Values

`GET /values` returns values for every state. Optional `state_id` and `year` parameters narrow the result. Collection screens should use the bulk form.

## Sources, imports, and scores

- `GET|POST /sources`; `PUT /sources/{sourceId}`
- `GET|POST /imports`; `GET /imports/{importId}`
- `GET /scores?profile_id=&year=`
- `POST /scores/recalculate` with `{ "profileId": 0, "year": 0 }`

Imports accept multipart `source_id` and `file`. CSV requires `state_code`, `metric_slug`, `year`, and `value`; `source_record_id` is optional.

`GET /health` reports process health. `GET /status` reports database readiness and active managed imports.
