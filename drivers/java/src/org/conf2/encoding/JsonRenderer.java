package org.conf2.encoding;

import java.io.IOException;
import java.io.OutputStream;
import java.io.OutputStreamWriter;
import java.io.Writer;

/**
 *
 */
public class JsonRenderer implements Renderer {
    private final static int QUOTE = '\"';
    private final static int COLON = ':';
    private final static int OPEN = '{';
    private final static int CLOSE = '}';
    private final static int COMMA = ',';
    private boolean first = true;
    private Writer out;
    private int depth = Integer.MAX_VALUE;
    private int currentDepth;

    public JsonRenderer(OutputStream out) {
        this.out = new OutputStreamWriter(out);
    }

    public JsonRenderer(Writer out) {
        this.out = out;
    }

    public void setDepth(int depth) {
        this.depth = depth;
    }

    @Override
    public void start() throws IOException {
        out.write(OPEN);
    }

    @Override
    public void enterContainer(String ident) throws IOException {
        this.currentDepth++;
        if (currentDepth >= depth) {
            return;
        }
        writeIdent(ident);
        out.write(OPEN);
        first = true;
    }

    void writeIdent(String ident) throws IOException {
        if (!first) {
            out.write(COMMA);
        }
        first = false;
        out.write(QUOTE);
        out.write(ident);
        out.write(QUOTE);
        out.write(COLON);
    }

    @Override
    public void writeValue(String ident, Object value) throws IOException {
        writeIdent(ident);
        out.write(QUOTE);
        out.write(value.toString());
        out.write(QUOTE);
    }

    @Override
    public void leaveContainer() throws IOException {
        this.currentDepth--;
        out.write(CLOSE);
    }

    @Override
    public void end() throws IOException {
        out.write(CLOSE);
        out.flush();
    }
}
