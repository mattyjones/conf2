package org.conf2.encoding;

import java.io.IOException;

/**
 *
 */
public interface Renderer {
    void start() throws IOException;
    void enterContainer(String ident) throws IOException;
    void writeValue(String ident, Object value) throws IOException;
    void leaveContainer() throws IOException;
    void end() throws IOException;
}
