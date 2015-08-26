package org.conf2.yang;

import org.junit.Test;

import java.io.IOException;
import java.io.InputStream;
import static org.junit.Assert.*;

/**
 *
 */
public class SimpleDataSourceTest {

    @Test
    public void testClassLoad() throws IOException {
        SimpleDataSource ds = new SimpleDataSource(SimpleDataSourceTest.class);
        InputStream is = ds.getResource("SimpleDataSourceTest.class");
        assertNotNull(is);
    }
}
