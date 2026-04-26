.PHONY: abigen

## abigen: regenerate Go contract bindings from forge build artifacts.
##   Runs forge build then abigen for every contract.
##   Pass SKIP_BUILD=1 to skip the forge step (artifacts must already exist).
abigen:
	SKIP_BUILD=$(SKIP_BUILD) scripts/abigen.sh
