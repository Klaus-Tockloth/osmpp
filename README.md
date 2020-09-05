# osmpp

## Purpose

'osmpp' (OSM preprocessing) is a tool to pre-process OSM data. It's typically used within an OSM data processing tool chain.

## Remarks

The master branch is used for program development and may be unstable. See 'Releases' for pre-build binaries.

## Build (master)

go get -u github.com/Klaus-Tockloth/osmpp

make

## Functionality

Processes node_network objects.

Processes turning_circle/loop objects.


## Usage

```txt
Program:
  Name                    : main
  Release                 : v0.2.0 - 2020/09/05
  Purpose                 : OSM data pre-processing
  Info                    : Processes node_network and turning_circle objects.

Usage:
  main -inputOSM=filename -outputNodes=filename -startNode=number

Example:
  main -inputOSM=osmdata.pbf -outputNodes=osmpp.xml -startNode=1000000000000

Options:
  -inputOSM string
    	name of OSM input file (PBF format)
  -outputNodes string
    	name of OSM nodes output file (XML format)
  -startNode int
    	starting ID for new nodes written to nodes output file
```
