package org.conf2.yang;

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;

/**
 * Read encoding from class path
 */
public class SimpleDataSource implements DataSource {
    private Class rootClass;
    private String baseClassPath;
    private File baseDir;

    public SimpleDataSource(File dir) {
        this.baseDir = dir;
    }

    public SimpleDataSource(String baseClassPath) {
        this.baseClassPath = baseClassPath;
    }

    public SimpleDataSource(Class root) {
        this.rootClass = root;
    }

    @Override
    public InputStream getResource(String resourceId) throws IOException {
System.out.println("SimpleDataSource.getResource");
        if (rootClass != null) {
System.out.println("SimpleDataSource - rootClass");
            InputStream is = rootClass.getResourceAsStream(resourceId);
System.out.println("SimpleDataSource - GOT IS");
            return is;
        } else if (this.baseDir != null) {
            return new FileInputStream(new File(baseDir, resourceId));
        }
        return ClassLoader.getSystemResourceAsStream(baseClassPath + resourceId);
    }
}
