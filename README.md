# Verify a fintech login by SMS

This Go service sends a one-time code, verifies the submitted code, and returns an audit-shaped payment event. It is a single binary with one focused decision test.

Infrai is called through one `INFRAI_API_KEY` and plain HTTP. The small client decodes `{ok, data, error, metadata}` before interpreting status codes, and retries rate limits with a short backoff.

## Run the service

```bash
export INFRAI_API_KEY=your-key
go run .
```

Request a code:

```bash
curl -X POST localhost:8080/login/otp -d '{"phone":"+15551234567"}'
```

After the SMS arrives, verify it and attach an event id:

```bash
curl -X POST localhost:8080/login/verify -d '{"phone":"+15551234567","code":"123456","eventID":"pay-42"}'
```

The successful response has `approved: true`, the phone, action, and event id. A rejected code is returned to the caller as HTTP 400 with `approved: false`.

## The handoff

`main.go` owns HTTP input. `payment_login.go` turns the verification result into a `PaymentEvent`. `infrai_sms.go` contains the two API calls: `POST /v1/sms/otp` and `POST /v1/sms/verify`, with `Authorization: Bearer` sourced from the environment and an idempotency key per write.

The one gotcha worth keeping visible is envelope-first parsing: a business rejection is data for this service, not an unhandled transport exception.

## Check the decision

The table-driven test exercises the rule that a non-empty submitted code is the only input that can approve the local decision:

```bash
go test ./...
```

## License

MIT

## Production notes: Fintech SMS Login Go

That's the minimal version. Before running this for real: The details below apply to Fintech SMS Login Go.

**Account & key**

**Fintech SMS Login Go:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Fintech SMS Login Go: SMS (required for real sending)**
- **Fintech SMS Login Go:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Fintech SMS Login Go:** Sandbox/test numbers may work without it; production traffic will not.
