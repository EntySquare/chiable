package main

import (
	"chiable/core"
	"errors"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {

	manualPlottingCMD := manual()
	cmanualPlottingCMD := cmanual()
	tmanualPlottingCMD := tmanual()
	autoPlottingCMD := auto()

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "debug",
				Value: "debug",
				Usage: "Debug commands",
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "manual",
				Usage:       "manual farming",
				Subcommands: []*cli.Command{manualPlottingCMD},
			},
			{
				Name:        "cmanual",
				Usage:       "manual farming",
				Subcommands: []*cli.Command{cmanualPlottingCMD},
			},
			{
				Name:        "tmanual",
				Usage:       "manual farming",
				Subcommands: []*cli.Command{tmanualPlottingCMD},
			},
			{
				Name:        "auto",
				Usage:       "automatic farming",
				Subcommands: []*cli.Command{autoPlottingCMD},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func auto() *cli.Command {
	cmd := &cli.Command{
		Name:  "plot",
		Usage: "plotting plots based on manual options",
		Flags: []cli.Flag{},
		Action: func(context *cli.Context) error {
			return core.NewDynamicStrategy().Run()
		},
	}
	return cmd

}

func manual() *cli.Command {
	cmd := &cli.Command{
		Name:  "plot",
		Usage: "plotting plots based on manual options",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "p",
				Usage: "used for plotting for cert pool",
			},
			&cli.StringFlag{
				Name:  "f",
				Usage: "used for plotting for cert farmer",
			},
			&cli.StringFlag{
				Name:  "d",
				Usage: "user dir use to match move",
			},
			&cli.StringFlag{
				Name:  "i",
				Usage: "docker image name",
			},
			&cli.StringFlag{
				Name:  "k",
				Usage: "k of plot",
			},
			&cli.StringFlag{
				Name:  "m",
				Usage: "manager ip to report",
			},
		},
		Action: func(context *cli.Context) error {
			n := context.Args().First()
			num, err := strconv.ParseInt(n, 10, 64)
			if err != nil {
				return errors.New("can't parse argument correctly")
			}
			return core.NewStaticStrategy(context.String("f"), context.String("p"),
				context.String("d"), context.String("i"), context.String("k"),
				context.String("m")).Run(num)
		},
	}
	return cmd
}

func cmanual() *cli.Command {
	cmd := &cli.Command{
		Name:  "plot",
		Usage: "plotting plots based on manual options",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "p",
				Usage: "used for plotting for cert pool",
			},
			&cli.StringFlag{
				Name:  "f",
				Usage: "used for plotting for cert farmer",
			},
			&cli.StringFlag{
				Name:  "d",
				Usage: "user dir use to match move",
			},
			&cli.StringFlag{
				Name:  "i",
				Usage: "docker image name",
			},
			&cli.StringFlag{
				Name:  "k",
				Usage: "k of plot",
			},
			&cli.StringFlag{
				Name:  "m",
				Usage: "manager ip to report",
			},
		},
		Action: func(context *cli.Context) error {
			n := context.Args().First()
			num, err := strconv.ParseInt(n, 10, 64)
			if err != nil {
				return errors.New("can't parse argument correctly")
			}
			return core.NewStaticStrategy(context.String("f"), context.String("p"),
				context.String("d"), context.String("i"), context.String("k"),
				context.String("m")).ChiaRun(num)
		},
	}
	return cmd
}

func tmanual() *cli.Command {
	cmd := &cli.Command{
		Name:  "plot",
		Usage: "plotting plots based on manual options",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "p",
				Usage: "used for plotting for cert pool",
			},
			&cli.StringFlag{
				Name:  "f",
				Usage: "used for plotting for cert farmer",
			},
			&cli.StringFlag{
				Name:  "d",
				Usage: "user dir use to match move",
			},
			&cli.StringFlag{
				Name:  "i",
				Usage: "docker image name",
			},
			&cli.StringFlag{
				Name:  "k",
				Usage: "k of plot",
			},
			&cli.StringFlag{
				Name:  "m",
				Usage: "manager ip to report",
			},
		},
		Action: func(context *cli.Context) error {
			n := context.Args().First()
			num, err := strconv.ParseInt(n, 10, 64)
			if err != nil {
				return errors.New("can't parse argument correctly")
			}
			return core.NewStaticStrategy(context.String("f"), context.String("p"),
				context.String("d"), context.String("i"), context.String("k"),
				context.String("m")).TestRun(num)
		},
	}
	return cmd
}
