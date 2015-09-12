package org.conf2.schema;

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;

/**
 * Read encoding from class path
 */
public class SimpleStreamSource implements StreamSource {
    private Class rootClass;
    private String baseClassPath;
    private File baseDir;

    public SimpleStreamSource(File dir) {
        this.baseDir = dir;
    }

    public SimpleStreamSource(String baseClassPath) {
        this.baseClassPath = baseClassPath;
    }

    public SimpleStreamSource(Class root) {
        this.rootClass = root;
    }

    @Override
    public InputStream getStream(String resourceId) throws IOException {
        if (rootClass != null) {
            InputStream is = rootClass.getResourceAsStream(resourceId);
            return is;
        } else if (this.baseDir != null) {
            return new FileInputStream(new File(baseDir, resourceId));
        }
        return ClassLoader.getSystemResourceAsStream(baseClassPath + resourceId);
    }
}
