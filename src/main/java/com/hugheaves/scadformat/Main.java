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

import java.io.ByteArrayOutputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.PrintStream;
import java.io.UncheckedIOException;
import java.nio.file.DirectoryStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;

import org.antlr.v4.runtime.ANTLRFileStream;
import org.antlr.v4.runtime.ANTLRInputStream;
import org.antlr.v4.runtime.CommonTokenStream;
import org.antlr.v4.runtime.misc.ParseCancellationException;
import org.antlr.v4.runtime.tree.ParseTree;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.beust.jcommander.JCommander;
import com.hugheaves.scadformat.antlr.SCADLexer;
import com.hugheaves.scadformat.antlr.SCADParser;

/**
 * The Class Main.
 */
public class Main {
    /** The logger. */
    private static Logger logger = LoggerFactory.getLogger(Main.class);

    /**
     * Instantiates a new main.
     */
    private Main() {

    }

    /**
     * Scan dir.
     *
     * @param fileNames
     *            the file names
     * @param path
     *            the path
     */
    protected void scanDir(final List<Path> fileNames, final Path path) {
        try {
            final DirectoryStream<Path> directoryStream = Files.newDirectoryStream(path);
            for (final Path subPath : directoryStream) {
                if (Files.isDirectory(subPath)) {
                    scanDir(fileNames, subPath);
                } else {
                    processFile(fileNames, subPath);
                }
            }
            directoryStream.close();
        } catch (final IOException e) {
            throw new UncheckedIOException(e);
        }
    }

    /**
     * Process file.
     *
     * @param fileNames
     *            the file names
     * @param subPath
     *            the sub path
     */
    private void processFile(final List<Path> fileNames, final Path subPath) {
        fileNames.add(subPath);
    }

    /**
     * The main method.
     *
     * @param args
     *            the arguments
     * @throws Exception
     *             the exception
     */
    public static void main(final String[] args) throws Exception {
        final Main main = new Main();
        main.run(args);
    }

    /**
     * Run.
     *
     * @param args
     *            the args
     */
    private void run(final String[] args) {
        System.out.println("scadformat 0.1");

        final CLIOptions options = new CLIOptions();
        final JCommander jCommander = new JCommander(options);
        jCommander.setProgramName("java -jar scadformat.jar");

        jCommander.parse(args);

        options.files.forEach(fileName -> formatFile(fileName));
    }

    /**
     * Parses the file.
     *
     * @param fileName
     *            the file name
     */
    private void formatFile(final String fileName) {
        try {
            logger.info("Parsing {} ", fileName);

            ANTLRInputStream antlrStream;
            try {
                antlrStream = new ANTLRFileStream(fileName);
            } catch (final IOException e) {
                throw new UncheckedIOException(e);
            }

            final SCADLexer lexer = new SCADLexer(antlrStream);

            final CommonTokenStream tokenStream = new CommonTokenStream(lexer);

            final SCADParser parser = new SCADParser(tokenStream);

            parser.removeErrorListeners();
            parser.addErrorListener(new ErrorListener());

            final ParseTree parseTree = parser.start();

            logger.info("Reformatting");
            final ByteArrayOutputStream byteStream = new ByteArrayOutputStream();
            final PrintStream printStream = new PrintStream(byteStream);
            final TokenFormatter formatter = new TokenFormatter(tokenStream, printStream);

            final SCADRenderingVisitor visitor = new SCADRenderingVisitor(tokenStream, formatter);

            visitor.visit(parseTree);

            final Path path = Paths.get(fileName);
            final Path newFile = path.normalize().getParent().resolve(path.getFileName().toString());

            logger.info("Writing {}", newFile.toString());
            final FileOutputStream outputStream = new FileOutputStream(newFile.toFile());
            outputStream.write(byteStream.toByteArray());
            outputStream.close();

        } catch (

        final ParseCancellationException e) {
            logger.error("Error parsing {}: {}", fileName, e);
        } catch (final IOException e) {
            logger.error("Error processing {}: {}", fileName, e);
        }
    }
}
