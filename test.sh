go test -coverprofile=coverage.out .
#go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html





// trace_use_otlp = sdsf[app_id ] or
                nmjsdf[serevice]  or
                trace_use_uid [$.header.uid]