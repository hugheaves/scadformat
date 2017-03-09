/*
 * Copyright (C) 2017  Hugh Eaves
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package com.hugheaves.scadformat;

import java.util.ArrayList;
import java.util.List;

import com.beust.jcommander.Parameter;
import com.beust.jcommander.Parameters;

/**
 * The Class CLIOptions.
 */
@Parameters()
public class CLIOptions {

    /**
     * The debug.
     */
    @Parameter(names = { "-D", "--debug" }, required = false, description = "Display debugging information")
    protected boolean debug;

    @Parameter(description = "Files", required = true)
    protected List<String> files = new ArrayList<>();

}
