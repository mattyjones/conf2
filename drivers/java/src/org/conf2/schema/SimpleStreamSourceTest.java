package org.conf2.schema;

import org.junit.Test;

import java.io.IOException;
import java.io.InputStream;
import static org.junit.Assert.*;

/**
 *
 */
public class SimpleStreamSourceTest {

    @Test
    public void testClassLoad() throws IOException {
        SimpleStreamSource ds = new SimpleStreamSource(SimpleStreamSourceTest.class);
        InputStream is = ds.getStream("SimpleStreamSourceTest.class");
        assertNotNull(is);
    }
}
