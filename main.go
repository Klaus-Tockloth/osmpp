/*
Purpose:
- OSM data pre-processing

Description:
- Duplicates OSM junction point nodes.

Releases:
- v0.1.0 - 2019/11/21 : initial release

Author:
- Klaus Tockloth

Copyright and license:
- Copyright (c) 2019 Klaus Tockloth
- MIT license

Permission is hereby granted, free of charge, to any person obtaining a copy of this software
and associated documentation files (the Software), to deal in the Software without restriction,
including without limitation the rights to use, copy, modify, merge, publish, distribute,
sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or
substantial portions of the Software.

The software is provided 'as is', without warranty of any kind, express or implied, including
but not limited to the warranties of merchantability, fitness for a particular purpose and
noninfringement. In no event shall the authors or copyright holders be liable for any claim,
damages or other liability, whether in an action of contract, tort or otherwise, arising from,
out of or in connection with the software or the use or other dealings in the software.

Contact (eMail):
- freizeitkarte@googlemail.com

Remarks:
- NN

Links:
- https://github.com/paulmach/osm
- https://github.com/paulmach/osm/blob/master/osmpbf/example_stats_test.go#L6
*/

package main

import (
	"bufio"
	"context"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

// general program info
var (
	_, progName = filepath.Split(os.Args[0])
	progVersion = "v0.1.0"
	progDate    = "2019/11/23"
	progPurpose = "OSM data pre-processing"
	progInfo    = "Duplicates OSM junction point nodes."
)

/*
init initializes this program
*/
func init() {

	// initialize logger
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

/*
main starts this program
*/
func main() {

	fmt.Printf("\nProgram:\n")
	fmt.Printf("  Name                    : %s\n", progName)
	fmt.Printf("  Release                 : %s - %s\n", progVersion, progDate)
	fmt.Printf("  Purpose                 : %s\n", progPurpose)
	fmt.Printf("  Info                    : %s\n", progInfo)

	// command line options
	inputOSM := flag.String("inputOSM", "", "name of OSM input file (PBF format)")
	outputNodes := flag.String("outputNodes", "", "name of OSM nodes output file (XML format)")
	startNode := flag.Int("startNode", 0, "starting ID for new nodes written to nodes output file")

	flag.Usage = printProgUsage
	flag.Parse()

	if *inputOSM == "" || *outputNodes == "" || *startNode == 0 {
		printProgUsage()
	}

	fmt.Printf("\nProcessing:\n")
	fmt.Printf("  OSM input file          : %s\n", *inputOSM)
	fmt.Printf("  Nodes output file       : %s\n", *outputNodes)
	fmt.Printf("  Starting node ID        : %d\n", *startNode)

	fileInput, err := os.Open(*inputOSM)
	if err != nil {
		log.Fatalf("could not open file: %v", err)
	}

	fileOutput, err := os.OpenFile(*outputNodes, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("could not open file: %v", err)
	}
	writer := bufio.NewWriter(fileOutput)
	_, err = fmt.Fprintf(writer, "<?xml version='1.0' encoding='UTF-8'?>\n")
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}
	_, err = fmt.Fprintf(writer, "<osm version='0.6' generator='%s'>\n", progName)
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}

	nodes, ways, relations := 0, 0, 0
	stats := newElementStats()

	nodeID := osm.NodeID(*startNode)
	junctionPointsFound := 0
	juctionPointsWritten := 0

	minLat, maxLat := math.MaxFloat64, -math.MaxFloat64
	minLon, maxLon := math.MaxFloat64, -math.MaxFloat64

	minTS, maxTS := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC), time.Time{}

	var (
		maxNodeRefs   int
		maxNodeRefsID osm.WayID
	)

	var (
		maxRelRefs   int
		maxRelRefsID osm.RelationID
	)

	scanner := osmpbf.New(context.Background(), fileInput, 3)
	defer scanner.Close()

	for scanner.Scan() {
		var ts time.Time

		switch e := scanner.Object().(type) {
		case *osm.Node:
			nodes++
			ts = e.Timestamp
			stats.Add(e.ElementID(), e.Tags)

			if e.Lat > maxLat {
				maxLat = e.Lat
			}

			if e.Lat < minLat {
				minLat = e.Lat
			}

			if e.Lon > maxLon {
				maxLon = e.Lon
			}

			if e.Lon < minLon {
				minLon = e.Lon
			}

			tags := e.TagMap()
			// id := e.ElementID()
			if len(tags) > 0 {
				tagValue, found := tags["network:type"]
				if found && tagValue == "node_network" {
					junctionPointsFound++
					nameKey := "name"
					nameValue, _ := tags[nameKey]

					refKey := "rcn_ref" // cycling
					refValue, found := tags[refKey]
					if found {
						duplicateNetworkJunctionPoint(writer, e, nodeID, refKey, refValue, nameKey, nameValue)
						nodeID++
						juctionPointsWritten++
					}

					refKey = "rwn_ref" // walking
					refValue, found = tags[refKey]
					if found {
						duplicateNetworkJunctionPoint(writer, e, nodeID, refKey, refValue, nameKey, nameValue)
						nodeID++
						juctionPointsWritten++
					}

					refKey = "rin_ref" // inline skating
					refValue, found = tags[refKey]
					if found {
						duplicateNetworkJunctionPoint(writer, e, nodeID, refKey, refValue, nameKey, nameValue)
						nodeID++
						juctionPointsWritten++
					}

					refKey = "rhn_ref" // horse riding
					refValue, found = tags[refKey]
					if found {
						duplicateNetworkJunctionPoint(writer, e, nodeID, refKey, refValue, nameKey, nameValue)
						nodeID++
						juctionPointsWritten++
					}

					refKey = "rpn_ref" // canoeing
					refValue, found = tags[refKey]
					if found {
						duplicateNetworkJunctionPoint(writer, e, nodeID, refKey, refValue, nameKey, nameValue)
						nodeID++
						juctionPointsWritten++
					}

					refKey = "rmn_ref" // motorboat driving
					refValue, found = tags[refKey]
					if found {
						duplicateNetworkJunctionPoint(writer, e, nodeID, refKey, refValue, nameKey, nameValue)
						nodeID++
						juctionPointsWritten++
					}
				}
			}
		case *osm.Way:
			ways++
			ts = e.Timestamp
			stats.Add(e.ElementID(), e.Tags)

			if l := len(e.Nodes); l > maxNodeRefs {
				maxNodeRefs = l
				maxNodeRefsID = e.ID
			}
		case *osm.Relation:
			relations++
			ts = e.Timestamp
			stats.Add(e.ElementID(), e.Tags)

			if l := len(e.Members); l > maxRelRefs {
				maxRelRefs = l
				maxRelRefsID = e.ID
			}
		}

		if ts.After(maxTS) {
			maxTS = ts
		}

		if ts.Before(minTS) {
			minTS = ts
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("scanner returned error: %v", err)
		os.Exit(1)
	}

	fmt.Printf("\nJunction point statistics:\n")
	fmt.Printf("  Points found            : %v\n", junctionPointsFound)
	fmt.Printf("  Nodes written           : %v\n", juctionPointsWritten)

	fmt.Printf("\nOSM data statistics:\n")
	fmt.Printf("  Timestamp min           : %v\n", minTS.Format(time.RFC3339))
	fmt.Printf("  Timestamp max           : %v\n", maxTS.Format(time.RFC3339))
	fmt.Printf("  Lon min                 : %0.7f\n", minLon)
	fmt.Printf("  Lon max                 : %0.7f\n", maxLon)
	fmt.Printf("  Lat min                 : %0.7f\n", minLat)
	fmt.Printf("  Lat max                 : %0.7f\n", maxLat)
	fmt.Printf("  Nodes                   : %v\n", nodes)
	fmt.Printf("  Ways                    : %v\n", ways)
	fmt.Printf("  Relations               : %v\n", relations)
	fmt.Printf("  Version max             : %v\n", stats.MaxVersion)
	fmt.Printf("  Node ID min             : %v\n", stats.Ranges[osm.TypeNode].Min)
	fmt.Printf("  Node ID max             : %v\n", stats.Ranges[osm.TypeNode].Max)
	fmt.Printf("  Way ID min              : %v\n", stats.Ranges[osm.TypeWay].Min)
	fmt.Printf("  Way ID max              : %v\n", stats.Ranges[osm.TypeWay].Max)
	fmt.Printf("  Relation ID min         : %v\n", stats.Ranges[osm.TypeRelation].Min)
	fmt.Printf("  Relation ID max         : %v\n", stats.Ranges[osm.TypeRelation].Max)
	fmt.Printf("  Keyval pairs max        : %v\n", stats.MaxTags)
	fmt.Printf("  Keyval pairs max object : %v %v\n", stats.MaxTagsID.Type(), stats.MaxTagsID.Ref())
	fmt.Printf("  Noderefs max            : %v\n", maxNodeRefs)
	fmt.Printf("  Noderefs max object     : way %v\n", maxNodeRefsID)
	fmt.Printf("  Relrefs max             : %v\n", maxRelRefs)
	fmt.Printf("  Relrefs max object      : relation %v\n", maxRelRefsID)

	_, err = fmt.Fprintf(writer, "</osm>\n")
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatalf("could not flush file buffer: %v", err)
	}
	err = fileOutput.Close()
	if err != nil {
		log.Fatalf("could not close file: %v", err)
	}
	err = fileInput.Close()
	if err != nil {
		log.Fatalf("could not close file: %v", err)
	}

	fmt.Printf("\n")
	os.Exit(0)
}

// Stats is a shared bit of code to accumulate stats from the element ids.
type elementStats struct {
	Ranges     map[osm.Type]*idRange
	MaxVersion int

	MaxTags   int
	MaxTagsID osm.ElementID
}

type idRange struct {
	Min, Max int64
}

func newElementStats() *elementStats {

	return &elementStats{
		Ranges: map[osm.Type]*idRange{
			osm.TypeNode:     {Min: math.MaxInt64},
			osm.TypeWay:      {Min: math.MaxInt64},
			osm.TypeRelation: {Min: math.MaxInt64},
		},
	}
}

func (s *elementStats) Add(id osm.ElementID, tags osm.Tags) {

	s.Ranges[id.Type()].Add(id.Ref())

	if v := id.Version(); v > s.MaxVersion {
		s.MaxVersion = v
	}

	if l := len(tags); l > s.MaxTags {
		s.MaxTags = l
		s.MaxTagsID = id
	}
}

func (r *idRange) Add(ref int64) {

	if ref > r.Max {
		r.Max = ref
	}

	if ref < r.Min {
		r.Min = ref
	}
}

/*
duplicateNetworkJunctionPoint duplicates node from junction point network.

<node id="355939532" lat="52.2220383" lon="7.022982600000001" user="" uid="0" visible="true" version="8" changeset="0" timestamp="2019-09-13T06:50:45Z">
  <tag k="expected_rcn_route_relations" v="3"></tag>
  <tag k="network:type" v="node_network"></tag>
  <tag k="rcn:name" v="Spechtholtshook"></tag>
  <tag k="rcn_ref" v="53"></tag>
  <tag k="rwn_ref" v="X32"></tag>
</node>
... will be transformed to:
<node id="xxxxxxx001" lat="52.2220383" lon="7.022982600000001" user="" uid="0" visible="true" version="8" changeset="0" timestamp="2019-09-13T06:50:45Z">
  <tag k="fzk_network:type" v="node_network"></tag>
  <tag k="rcn_ref" v="53"></tag>
  <tag k="name" v="Spechtholtshook"></tag>
</node>
<node id="xxxxxxx002" lat="52.2220383" lon="7.022982600000001" user="" uid="0" visible="true" version="8" changeset="0" timestamp="2019-09-13T06:50:45Z">
  <tag k="fzk_network:type" v="fzk_network:type"></tag>
  <tag k="rwn_ref" v="X32"></tag>
  <tag k="name" v="Spechtholtshook"></tag>
</node>
*/
func duplicateNetworkJunctionPoint(writer *bufio.Writer, sourceOsmNode *osm.Node, nodeID osm.NodeID, refKey, refValue, nameKey, nameValue string) {

	newOsmNode := sourceOsmNode
	newOsmNode.ID = nodeID
	newOsmNode.Tags = []osm.Tag{}

	tag := osm.Tag{Key: "fzk_network:type", Value: "node_network"}
	newOsmNode.Tags = append(newOsmNode.Tags, tag)

	tag.Key = refKey
	tag.Value = refValue
	newOsmNode.Tags = append(newOsmNode.Tags, tag)

	if nameValue != "" {
		tag.Key = nameKey
		tag.Value = nameValue
		newOsmNode.Tags = append(newOsmNode.Tags, tag)
	}

	data, err := xml.MarshalIndent(newOsmNode, "  ", "  ")
	if err != nil {
		log.Fatalf("error <%v> at xml.MarshalIndent()", err)
	}

	_, err = fmt.Fprintf(writer, "%s\n", string(data))
	if err != nil {
		log.Fatalf("error writing output file: %v", err)
	}
}

/*
Print program usage.
*/
func printProgUsage() {

	fmt.Printf("\nUsage:\n")
	fmt.Printf("  %s -inputOSM=filename -outputNodes=filename -startNode=number\n", progName)

	fmt.Printf("\nExample:\n")
	fmt.Printf("  %s -inputOSM=osmdata.pbf -outputNodes=osmnodes.xml -startNode=10000000000\n", progName)

	fmt.Printf("\nOptions:\n")
	flag.PrintDefaults()

	os.Exit(1)
}
