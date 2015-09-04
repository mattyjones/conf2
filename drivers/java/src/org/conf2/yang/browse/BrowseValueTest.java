package org.conf2.yang.browse;

import org.junit.Test;

import java.io.IOException;
import java.nio.ByteBuffer;
import static org.junit.Assert.*;

/**
 *
 */
public class BrowseValueTest {

    @Test
    public void testEncodeCStr() throws IOException {
        String[] strs = new String[] {
                "a",
                "b"
        };
        ByteBuffer buff = BrowseValue.encodeCStrArray(strs);
        assertEquals(4, buff.capacity());
        buff.rewind();
        byte[] b = new byte[4];
        buff.get(b);
        assertEquals("a\0b\0", new String(b));
    }

    @Test
    public void testDecodeCStr() throws IOException {
        byte[] b = "a\0b\0".getBytes();
        ByteBuffer buff = ByteBuffer.wrap(b);
        String[] strs = BrowseValue.decodeCStrArray(buff, 2);
        assertEquals(2, strs.length);
        assertEquals("a", strs[0]);
        assertEquals("b", strs[1]);
    }

    @Test
    public void testCInt() throws IOException {
        int[] intlist = new int[] { 10, 20 };
        ByteBuffer buff = BrowseValue.encodeCIntArray(intlist);
        assertEquals(8, buff.capacity());
        buff.rewind();
        int[] roundtrip = BrowseValue.decodeCIntArray(buff, 2);
        assertArrayEquals(intlist, roundtrip);
    }

    @Test
    public void testCBool() throws IOException {
        boolean[] boollist = new boolean[] { true, false };
        ByteBuffer buff = BrowseValue.encodeCBoolArray(boollist);
        assertEquals(4, buff.capacity());
        buff.rewind();
        boolean[] roundtrip = BrowseValue.decodeCBoolArray(buff, 2);
        assertArrayEquals(boollist, roundtrip);
    }
}
