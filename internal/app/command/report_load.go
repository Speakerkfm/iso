package command

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Speakerkfm/iso/internal/app/command/adapter/tablewriter"
)

var headers = []string{"service name", "method name", "success count", "error count"}

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
			table.Append([]string{
				serviceName,
				methodName,
				strconv.Itoa(methodReport.Stat.SuccessCount),
				strconv.Itoa(methodReport.Stat.ErrorCount),
			})
		}
	}

	table.Render()

	fmt.Print(builder.String())

	return nil
}
