package org.conf2.schema.comm;

import org.conf2.schema.ValueType;
import org.conf2.schema.browse.BrowseValue;
import org.junit.Test;

import java.io.IOException;
import java.nio.ByteBuffer;

import static org.junit.Assert.assertArrayEquals;
import static org.junit.Assert.assertEquals;

/**
 *
 */
public class DataWriterTest {
    @Test
    public void testEncodeCStr() throws IOException {
        String[] strs = new String[] {
                "a",
                "b"
        };
        DataWriter w = new DataWriter();
        BrowseValue v = new BrowseValue();
        v.valType = ValueType.INT32;
        v.int32 = 99;
        w.writeValue(v);
        byte[] actual = w.getByteBuffer().array();
        byte[] expected = new byte[] {10, 0, 0, 0, 99, 0, 0, 0};
        assertArrayEquals(expected, actual);
    }
}
