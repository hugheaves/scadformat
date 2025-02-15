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

package logutil

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ConfigureLogging(logLevel string) error {
	if logLevel == "" {
		logLevel = os.Getenv("LOG_LEVEL")
	}
	if logLevel == "" {
		logLevel = "error"
	}

	level, err := zap.ParseAtomicLevel(logLevel)
	if err != nil {
		return err
	}
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.TimeKey = zapcore.OmitKey
	if !level.Enabled(zap.DebugLevel) {
		encoderConfig.StacktraceKey = zapcore.OmitKey
		encoderConfig.CallerKey = zapcore.OmitKey
	}

	config := zap.Config{
		Level:            level,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var logger *zap.Logger
	logger, err = config.Build()
	zap.ReplaceGlobals(logger)
	return nil
}
