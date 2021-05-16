#!/usr/bin/env sh

export {{.AppIdName}}="{{.AppIdValue}}"
export {{.AppKeyName}}="{{.AppKeyValue}}"

. "/root/.acme.sh/acme.sh.env"

acme.sh --upgrade

acme.sh --issue --dns {{.DnsApi}} -d {{.MainDomain}} {{.ExtraDomain}}