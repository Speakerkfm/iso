package command

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Speakerkfm/iso/internal/app/command/adapter/tablewriter"
)

var headers = []string{"service name", "method name", "rule name", "request count"}

func (c *Command) ReportLoad(ctx context.Context) error {
	report, err := c.isoSrv.GetReport(ctx)
	if err != nil {
		return fmt.Errorf("fail to get report from isoserver: %w", err)
	}

	builder := &strings.Builder{}

	_, _ = builder.WriteString("ISO Report\n")

	table := tablewriter.NewWriter(builder)
	table.SetHeader(headers)
	table.SetRowLine(true)

	for serviceName, serviceReport := range report.Service {
		for methodName, methodReport := range serviceReport.Method {
			for ruleName, reqCnt := range methodReport.RuleStat {
				table.Append([]string{
					serviceName,
					methodName,
					ruleName,
					strconv.Itoa(int(reqCnt)),
				})
			}
		}
	}

	table.Render()

	fmt.Print(builder.String())

	return nil
}
