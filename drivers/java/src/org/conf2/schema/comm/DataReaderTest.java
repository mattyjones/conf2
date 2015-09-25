package org.conf2.schema.comm;

import org.conf2.schema.browse.BrowseValue;
import org.junit.Test;

import java.io.IOException;
import java.nio.ByteBuffer;

import static org.junit.Assert.assertEquals;

/**
 *
 */
public class DataReaderTest {
    @Test
    public void testDecodeCStr() throws IOException {
        byte[] b = "a\0b\0".getBytes();
        DataReader r = new DataReader(ByteBuffer.wrap(b));
        assertEquals("a", r.readString());
        assertEquals("b", r.readString());
    }
}
