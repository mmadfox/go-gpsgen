package draw

import (
	"fmt"
	"os"
	"strings"

	"github.com/mmadfox/go-gpsgen"
	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/olekukonko/tablewriter"
)

func Table(s *gpsgen.State) {
	if s == nil {
		fmt.Println("state is <nil>")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetColMinWidth(0, 30)
	table.SetColMinWidth(1, 30)
	table.SetColMinWidth(4, 40)
	table.SetNewLine("\n")

	table.SetHeader([]string{
		"Device",
		"Location",
		"Speed (m/s)",
		"Distance (meters)",
		"Sensors",
	})

	table.Append([]string{
		model2str(s.Model, s.Descr, s.Props),
		location2str(&s.Location),
		f2s(s.Speed),
		dist2str(s.Location.TotalDistance, s.Location.CurrentDistance),
		sensor2str(s.Sensors),
	})

	table.Render()
}

func sensor2str(sensors map[string][2]float64) string {
	sb := strings.Builder{}
	for name, val := range sensors {
		sb.WriteString(name)
		sb.WriteString(":")
		sb.WriteString("valX=")
		sb.WriteString(f2s(val[0]))
		sb.WriteString(" ")
		sb.WriteString("valY=")
		sb.WriteString(f2s(val[1]))
		sb.WriteString("\n")
	}
	return sb.String()
}

func model2str(model, descr string, props gpsgen.Properties) string {
	sb := strings.Builder{}
	sb.WriteString("Model:")
	sb.WriteString(model)
	sb.WriteString("\n")
	sb.WriteString("Descr:")
	sb.WriteString(descr)
	sb.WriteString("\n")
	sb.WriteString("Properties:")
	sb.WriteByte('\n')
	for k, v := range props {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(v)
		sb.WriteString(",")
	}
	return sb.String()
}

func dist2str(total, current float64) string {
	sb := strings.Builder{}
	sb.WriteString("Current:")
	sb.WriteString(f2s(current))
	sb.WriteString("\n")
	sb.WriteString("Total:")
	sb.WriteString(f2s(total))
	sb.WriteString("\n")
	return sb.String()
}

func location2str(loc *navigator.Location) string {
	sb := strings.Builder{}
	sb.WriteString("Lon:")
	sb.WriteString(f2s(loc.Lon))
	sb.WriteString("\n")
	sb.WriteString("Lat:")
	sb.WriteString(f2s(loc.Lat))
	sb.WriteString("\n")
	sb.WriteString("Elevation:")
	sb.WriteString(f2s(loc.Alt))
	sb.WriteString("\n")
	sb.WriteString("Bearing:")
	sb.WriteString(f2s(loc.Bearing))
	sb.WriteString("\n")
	sb.WriteString("DMSLat:")
	sb.WriteString(loc.LatDMS.String())
	sb.WriteString("\n")
	sb.WriteString("DMSLon:")
	sb.WriteString(loc.LonDMS.String())
	sb.WriteString("\n")
	return sb.String()
}

func f2s(v float64) string {
	return fmt.Sprintf("%f", v)
}
