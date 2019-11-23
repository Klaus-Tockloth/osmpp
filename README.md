# osmpp

## Purpose

'osmpp' (OSM preprocessing) is a tool to preprocess OSM data. It's typically used within an OSM data processing tool chain.

## Remarks

The master branch is used for program development and may be unstable. See 'Releases' for pre-build binaries.

## Build (master)

go get -u github.com/Klaus-Tockloth/osmpp

make

## Functionality

Duplicates OSM junction point nodes.


## Usage

```txt
Program:
  Name                    : osmpp
  Release                 : v0.1.0 - 2019/11/23
  Purpose                 : OSM data pre-processing
  Info                    : Duplicates OSM junction point nodes.

Usage:
  osmpp -inputOSM=filename -outputNodes=filename -startNode=number

Example:
  osmpp -inputOSM=osmdata.pbf -outputNodes=osmnodes.xml -startNode=10000000000

Options:
  -inputOSM string
    	name of OSM input file (PBF format)
  -outputNodes string
    	name of OSM nodes output file (XML format)
  -startNode int
    	starting ID for new nodes written to nodes output file
```
