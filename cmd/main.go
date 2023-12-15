// SCADFormat - Formatter / beautifier for OpenSCAD source code
//
// Copyright (C) 2023  Hugh Eaves
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.

package main

import (
	_ "embed"
	"strings"

	"github.com/hugheaves/scadformat/internal/formatter"
	"github.com/hugheaves/scadformat/internal/logutil"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

//go:generate sh -c "git describe > version.txt"
//go:embed version.txt
var gitVersion string

func main() {
	err := logutil.ConfigureLogging("error")
	if err != nil {
		panic(err)
	}

	var logLevel string
	pflag.StringVar(&logLevel, "log-level", "info", "Logging level (one of debug, info, warn, or error)")
	pflag.Parse()

	err = logutil.ConfigureLogging(logLevel)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	zap.L().Info("SCADFormat " + strings.TrimSpace(gitVersion))

	var fileName string
	if len(pflag.Args()) == 1 {
		fileName = pflag.Arg(0)
	} else if len(pflag.Args()) > 1 {
		zap.L().Fatal("only a single filename may be specified on the command line")
	}

	formatter := formatter.NewFormatter(fileName)
	err = formatter.Format()
	if err != nil {
		zap.L().Fatal(err.Error())
	}

}
