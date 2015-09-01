package org.conf2.yang.browse;

import org.junit.Test;
import static org.junit.Assert.*;

/**
 *
 */
public class BrowseUtilTest {

    @Test
    public void testGetterMethodNameFromMeta() {
        assertEquals("xFoo", BrowseUtil.accessorMethodNameFromMeta("x", "foo"));
        assertEquals("xFooBar", BrowseUtil.accessorMethodNameFromMeta("x", "foo-bar"));
        assertEquals("fooBar", BrowseUtil.accessorMethodNameFromMeta("", "foo-bar"));
    }
}
