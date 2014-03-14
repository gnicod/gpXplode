package main

import (
	"fmt"
	"flag"
	"time"
	"io"
	"os"
	"io/ioutil"
	"encoding/xml"
	"sort"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}


type Trk struct{
	XMLName string `xml:"trkseg"`
}

type GpxWpt struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	Ele         float64 `xml:"ele,omitempty"`
	Timestamp   time.Time  `xml:"time,omitempty"`
	MagVar      string  `xml:"magvar,omitempty"`
	GeoIdHeight string  `xml:"geoidheight,omitempty"`
	Name  string    `xml:"name,omitempty"`
	Cmt   string    `xml:"cmt,omitempty"`
	Desc  string    `xml:"desc,omitempty"`
	Src   string    `xml:"src,omitempty"`
	Sym   string    `xml:"sym,omitempty"`
	Type  string    `xml:"type,omitempty"`
	Fix          string  `xml:"fix,omitempty"`
	Sat          int     `xml:"sat,omitempty"`
	Hdop         float64 `xml:"hdop,omitempty"`
	Vdop         float64 `xml:"vdop,omitempty"`
	Pdop         float64 `xml:"pdop,omitempty"`
	AgeOfGpsData float64 `xml:"ageofgpsdata,omitempty"`
	DGpsId       int     `xml:"dgpsid,omitempty"`

}

type GpxTrkseg struct {
	XMLName xml.Name `xml:"trkseg"`
	Points  []GpxWpt `xml:"trkpt"`

}

type GpxTrk struct {
	XMLName  xml.Name    `xml:"trk"`
	Name     string      `xml:"name,omitempty"`
	Cmt      string      `xml:"cmt,omitempty"`
	Desc     string      `xml:"desc,omitempty"`
	Src      string      `xml:"src,omitempty"`
	Number   int         `xml:"number,omitempty"`
	Type     string      `xml:"type,omitempty"`
	Segments []GpxTrkseg `xml:"trkseg"`

}

type Gpx struct {
	XMLName      xml.Name     `xml:"http://www.topografix.com/GPX/1/1 gpx"`
	XmlNsXsi     string       `xml:"xmlns:xsi,attr,omitempty"`
	XmlSchemaLoc string       `xml:"xsi:schemaLocation,attr,omitempty"`
	Version      string       `xml:"version,attr"`
	Creator      string       `xml:"creator,attr"`
	Tracks       []GpxTrk     `xml:"trk"`
}

type GpxWpts  []GpxWpt

func (slice GpxWpts) Len() int {
	return len(slice)
}

func (slice GpxWpts) Less(i, j int) bool {
	return slice[j].Timestamp.After(slice[i].Timestamp);
}

func (slice GpxWpts) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

var gpxFiles []string

func main() {
	outputFolder := flag.String("output", "./", "Folder where the gpx files will be create")
	flag.Parse()

	stat, err := os.Stat(*outputFolder)
	if os.IsNotExist(err){
		fmt.Printf("no such directory: %s", *outputFolder)
		return
	}
	if stat.IsDir() == false{
		fmt.Printf("%s is not a directory", *outputFolder)
		return
	}

	gpxFiles = flag.Args()
	trksegs         := []GpxTrkseg{}
	trkwpts_by_date := make(map[time.Time][]GpxWpt)

	var gpx Gpx

	for _,file := range gpxFiles {
		gpx = Gpx{}
		dat, err := ioutil.ReadFile(file)
		check(err)
		erro := xml.Unmarshal([]byte(dat), &gpx)
		check(erro)
		for _,track := range gpx.Tracks {
			for _,trkseg := range track.Segments {
				trksegs = append(trksegs,trkseg)
				for _,point := range trkseg.Points {
					t := time.Date(point.Timestamp.Year(), point.Timestamp.Month(), point.Timestamp.Day(), 0, 0, 0, 0, time.UTC)
					trkwpts_by_date[t] = append(trkwpts_by_date[t],point)
				}
			}
		}
	}

	var gpxFinal Gpx
	gpxFinal.Version = "1.1"
	gpxFinal.Creator = "gpXplode"
	gpxFinal.XmlSchemaLoc = "http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd"
	gpxFinal.XmlNsXsi = "http://www.w3.org/2001/XMLSchema-instance"

	for date,wptd := range trkwpts_by_date{
		wpts := GpxWpts(wptd)
		sort.Sort(wpts)
		var gpxTrckSeg GpxTrkseg
		var gpxTrck  GpxTrk
		gpxTrckSeg.Points = wptd
		gpxTrck.Name = date.String()
		gpxTrck.Segments = []GpxTrkseg{gpxTrckSeg}
		gpxFinal.Tracks = []GpxTrk{gpxTrck}
		gpXml,_ := xml.Marshal(gpxFinal)

		fName := fmt.Sprintf("%v-%v-%v",date.Day(),date.Month(),date.Year())
		f, err := os.Create(*outputFolder+"/"+fName+".gpx")
		if err != nil {
			fmt.Println(err)

		}
		n, err := io.WriteString(f, string(gpXml))
		if err != nil {
			fmt.Println(n, err)

		}
		f.Close()
	}

}
