package bi_shows_model

import (
	"guduo/app/migration/internal"
	"testing"
)

func TestSync(t *testing.T) {
	internal.InitDB()
	Sync()
}