package draw

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/mmadfox/go-gpsgen/proto"
	"github.com/olekukonko/tablewriter"
)

func Table(s *proto.Device) {
	if s == nil {
		fmt.Println("state is <nil>")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetColMinWidth(0, 25)
	table.SetColMinWidth(1, 30)
	table.SetColMinWidth(2, 30)
	table.SetColMinWidth(3, 30)

	table.SetHeader([]string{
		"Device",
		"Location",
		"Sensors",
		"Custom sensors",
	})

	table.Append([]string{
		model2str(s),
		location2str(s.Location),
		defaultSensors(s),
		sensor2str(s.Sensors),
	})

	table.Render()
}

func defaultSensors(s *proto.Device) string {
	sb := strings.Builder{}
	sb.WriteString("Speed:")
	sb.WriteString(f22s(s.Speed))
	sb.WriteString("m/s\n")
	sb.WriteString("BatteryCharge:")
	sb.WriteString(f22s(s.BatteryCharge))
	sb.WriteString("%\n")
	sb.WriteString("ChargeTime:")
	sb.WriteString(time.Duration(s.BatterChargeTime).String())
	sb.WriteString("\n")
	sb.WriteString("Duration:")
	sb.WriteString(f2dur(s.Duration))
	sb.WriteString("\n")
	sb.WriteString("TotalDist:")
	sb.WriteString(f22s(s.Location.TotalDistance))
	sb.WriteString("m.\n")
	sb.WriteString("CurDist:")
	sb.WriteString(f22s(s.Location.CurrentDistance))
	sb.WriteString("m.\n")
	sb.WriteString("CurSegDist:")
	sb.WriteString(f22s(s.Location.CurrentSegmentDistance))
	sb.WriteString("m.\n")
	sb.WriteString("SegDist:")
	sb.WriteString(f22s(s.Location.SegmentDistance))
	sb.WriteString("m.\n")
	sb.WriteString("Tick:")
	sb.WriteString(f22s(s.Tick))
	sb.WriteString("s\n")
	return sb.String()
}

func sensor2str(sensors []*proto.Sensor) string {
	sb := strings.Builder{}
	for _, sensor := range sensors {
		sb.WriteString(sensor.Name)
		sb.WriteString(":")
		sb.WriteString("valX=")
		sb.WriteString(f2s(sensor.ValX))
		sb.WriteString(" ")
		sb.WriteString("valY=")
		sb.WriteString(f2s(sensor.ValY))
		sb.WriteString("\n")
	}
	return sb.String()
}

func model2str(s *proto.Device) string {
	sb := strings.Builder{}
	sb.WriteString("Model:")
	sb.WriteString(s.Model)
	sb.WriteString("\n")
	sb.WriteString("Status:")
	switch s.Online {
	case true:
		sb.WriteString("Online")
	case false:
		sb.WriteString("Offline")
	}
	sb.WriteString("\n")
	sb.WriteString("RouteIndex:")
	sb.WriteString(d2s(s.Location.RouteIndex))
	sb.WriteString("\n")
	sb.WriteString("TrackIndex:")
	sb.WriteString(d2s(s.Location.TrackIndex))
	sb.WriteString("\n")
	sb.WriteString("SegmentIndex:")
	sb.WriteString(d2s(s.Location.SegmentIndex))
	sb.WriteString("\n")
	sb.WriteString("Descr:")
	sb.WriteString(s.Descr)
	sb.WriteString("\n")
	if len(s.Props) > 0 {
		sb.WriteString("Properties:")
		sb.WriteByte('\n')
		for k, v := range s.Props {
			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(v)
			sb.WriteString(",")
		}
	}
	return sb.String()
}

func location2str(loc *proto.Location) string {
	sb := strings.Builder{}
	sb.WriteString("Lon:")
	sb.WriteString(f2s(loc.Lon))
	sb.WriteString("\n")
	sb.WriteString("Lat:")
	sb.WriteString(f2s(loc.Lat))
	sb.WriteString("\n")
	sb.WriteString("Elevation:")
	sb.WriteString(f2s(loc.Alt))
	sb.WriteString("m\n")
	sb.WriteString("Bearing:")
	sb.WriteString(f2s(loc.Bearing))
	sb.WriteString(" ")
	sb.WriteString("DMSLat:")
	sb.WriteString(navigator.FormatDMS(loc.LatDms))
	sb.WriteString(" ")
	sb.WriteString("DMSLon:")
	sb.WriteString(navigator.FormatDMS(loc.LonDms))
	sb.WriteString("\n")
	sb.WriteString(navigator.FormatUTM(loc.Utm))
	return sb.String()
}

func d2s(v int64) string {
	return fmt.Sprintf("%d", v)
}

func f2s(v float64) string {
	return fmt.Sprintf("%f", v)
}

func f22s(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

func f2dur(v float64) string {
	return time.Duration(time.Duration(v) * time.Second).String()
}
