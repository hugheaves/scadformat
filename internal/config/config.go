// SCADFormat - Formatter / beautifier for OpenSCAD source code
//
// Copyright (C) 2025  Hugh Eaves
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

package config

type MainConfig struct {
	LogLevel          string // logging level - error, warn, info, etc.
	Watch             bool   // filesystem "watch mode" is enabled
	Recurse           bool   // when target is a directory, apply operation recursively
	NoBackups         bool   // do not create backups of modified files
	PreserveTimestamp bool   // preserve the modified time on reformatted files
	TargetPath        string // the target path of the operation(fie or directory)
}
