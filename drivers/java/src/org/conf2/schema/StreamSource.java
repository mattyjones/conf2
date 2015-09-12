package org.conf2.schema;

import java.io.IOException;
import java.io.InputStream;

/**
 * API to load chunked content from Java into Go
 */
public interface StreamSource {
    public InputStream getStream(String resourceId) throws IOException;
}
