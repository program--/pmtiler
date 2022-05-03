package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	gj "pmtiler/geojson"
	ptio "pmtiler/io"

	"github.com/paulmach/orb/maptile"
	"github.com/urfave/cli/v2"
)

func main() {
	var (
		output string
		xcol   string
		ycol   string
		zoom   int
	)

	app := &cli.App{
		Name:  "pmtiler",
		Usage: "Write PMTiles directly from various formats.",
		Authors: []*cli.Author{
			{
				Name:  "Justin Singh-Mohudpur",
				Email: "justin@justinsingh.me",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Value:       ".",
				Usage:       "Directory or file to output PMTiles to/as.",
				Destination: &output,
			},
			&cli.StringFlag{
				Name:        "xcol",
				Aliases:     []string{"x"},
				Value:       "X",
				Usage:       "X-coordinate column within the input file.",
				Destination: &xcol,
			},
			&cli.StringFlag{
				Name:        "ycol",
				Aliases:     []string{"y"},
				Value:       "Y",
				Usage:       "Y-coordinate column within the input file.",
				Destination: &ycol,
			},
			&cli.IntFlag{
				Name:        "maxzoom",
				Aliases:     []string{"z"},
				Value:       15,
				Usage:       "Maximum Zoom to write.",
				Destination: &zoom,
			},
		},
		Action: func(c *cli.Context) error {
			input := c.Args().Get(0)
			base := path.Base(input)
			ext := path.Ext(input)
			base = base[0 : len(base)-len(ext)]

			if output == "." {
				output = filepath.Join(".", base+".pmtiles")
			}

			fc, err := ptio.ParquetToGeoJSON(input, xcol, ycol)
			if err != nil {
				log.Fatal(err)
			}

			err = gj.GeoJSONToTiles(output, fc, base, maptile.Zoom(zoom))
			if err != nil {
				log.Fatal(err)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// pmtiler [global options] command [command options] [arguments...]
// pmtiler -o tiles.pmtiles s3://...
