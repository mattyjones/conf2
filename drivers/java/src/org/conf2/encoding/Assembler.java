package org.conf2.encoding;

import java.io.IOException;

/**
 *
 */
public interface Assembler {
    void write(Renderer wtr) throws IOException;
}
