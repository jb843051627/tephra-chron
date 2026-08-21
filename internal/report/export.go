package report

import (
	"fmt"
	"github.com/jb843051627/tephra-chron/internal/model"
	"strings"
)

func Text(r model.ChronologyReport) string {
	return fmt.Sprintf("sequence=%s accepted=%t mean=%.3f notes=%s", r.SequenceID, r.Accepted, r.MeanAge, strings.Join(r.Notes, ","))
}
