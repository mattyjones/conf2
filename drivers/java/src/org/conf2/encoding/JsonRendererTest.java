package org.conf2.encoding;

import org.junit.Before;
import org.junit.Test;
import static org.junit.Assert.assertEquals;

import java.io.IOException;
import java.io.StringWriter;

/**
 *
 */
public class JsonRendererTest {
    StringWriter actual;
    JsonRenderer wtr;

    @Before
    public void setUp() {
        actual = new StringWriter();
        wtr = new JsonRenderer(actual);
    }

    @Test
    public void testWrite() throws IOException {
        wtr.start();
        wtr.end();
        assertEquals("{}", actual.toString());
    }

    @Test
    public void testContainer() throws IOException {
        wtr.start();
        wtr.enterContainer("x");
        wtr.leaveContainer();
        wtr.end();
        assertEquals("{\"x\":{}}", actual.toString());
    }

    @Test
    public void testContainerWithValues() throws IOException {
        wtr.start();
        wtr.enterContainer("x");
        wtr.writeValue("y", 1);
        wtr.writeValue("z", 2);
        wtr.leaveContainer();
        wtr.end();
        assertEquals("{\"x\":{\"y\":\"1\",\"z\":\"2\"}}", actual.toString());
    }
}
