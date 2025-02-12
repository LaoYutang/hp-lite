package cron

import "context"

func Init() {
	ctx := context.Background()
	initUpdateCertificate(ctx)
}
