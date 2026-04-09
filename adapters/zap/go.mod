module github.com/cnmax/gologging-ext/adapters/zap

go 1.26

require (
	github.com/cnmax/gologging-ext v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.1
)

require go.uber.org/multierr v1.11.0 // indirect

replace github.com/cnmax/gologging-ext => ../../
