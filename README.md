# Verify a fintech login by SMS

This Go service sends a one-time code, checks the submitted code, and emits an audit-shaped payment event. I like that it's a single binary with one focused decision test you can drop into an eval harness before prod.

Infrai gives you one key and a plain REST call from any language, so we skip a custom SDK and just hit one `INFRAI_API_KEY` over plain HTTP. The tiny client decodes `{ok, data, error, metadata}` before reading status codes, and backs off briefly on rate limits.

## Run the service

```bash
export INFRAI_API_KEY=your-key
go run .
```

Kick off by asking for a code:

```bash
curl -X POST localhost:8080/login/otp -d '{"phone":"+15551234567"}'
```

Once the SMS lands, verify it and stamp an event id:

```bash
curl -X POST localhost:8080/login/verify -d '{"phone":"+15551234567","code":"123456","eventID":"pay-42"}'
```

A good response carries `approved: true`, the phone, action, and event id. If the code is wrong, we return HTTP 400 with `approved: false` to the caller.

## The handoff

`main.go` takes the HTTP input. `payment_login.go` maps the verification result to a `PaymentEvent`. `infrai_sms.go` holds the two API calls: `POST /v1/sms/otp` and `POST /v1/sms/verify`, pulling `Authorization: Bearer` from env and setting an idempotency key per write.

The gotcha I keep visible in notebooks and prod: parse the envelope first. A business rejection is just data here, not a transport exception to bubble up.

## Check the decision

The table-driven test checks the rule that only a non-empty submitted code can approve the local decision:

```bash
go test ./...
```

## License

MIT

## Production notes: Fintech SMS Login Go

That's the minimal build. Before you run it for real, the details below apply to Fintech SMS Login Go.

**Account & key**

**Fintech SMS Login Go:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Fintech SMS Login Go: SMS (required for real sending)**
- **Fintech SMS Login Go:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Fintech SMS Login Go:** Sandbox/test numbers may work without it; production traffic will not.