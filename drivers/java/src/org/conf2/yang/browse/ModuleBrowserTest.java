package org.conf2.yang.browse;

import org.conf2.yang.*;
import org.junit.Test;

import static org.junit.Assert.*;

/**
 *
 */
public class ModuleBrowserTest {

    @Test
    public void testGetterMethodNameFromMeta() {
        assertEquals("xFoo", ModuleBrowser.accessorMethodNameFromMeta("x", "foo"));
    }

    @Test
    public void testRootSelector() {
        ModuleBrowser b = new ModuleBrowser(YangModule.YANG);
        Selection s = b.getRootSelector();
        Selection ms = Walk.walk(s, new BrowsePath("module"));
        assertNotNull(ms);
        assertEquals("module", ms.meta.getIdent());
        assertNotNull(ms.meta);
        ms.position = MetaUtil.findByIdent(ms.meta, "prefix");
        assertNotNull(ms.position);
        BrowseValue v = new BrowseValue();
        ms.Read.Read(v);
        assertEquals("yang", v.str);
    }
}
