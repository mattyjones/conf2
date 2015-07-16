package org.conf2.yang;

import java.io.IOException;
import java.io.InputStream;

/**
 * API to load chunked content from Java into Go
 */
public interface DataSource {
    public InputStream getResource(String resourceId) throws IOException;
}
