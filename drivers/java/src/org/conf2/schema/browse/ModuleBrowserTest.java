package org.conf2.schema.browse;

import org.conf2.schema.*;
import org.conf2.schema.yang.YangModule;
import org.junit.Test;

import static org.junit.Assert.*;

/**
 *
 */
public class ModuleBrowserTest {

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
        BrowseValue v = ms.Read.Read();
        assertEquals("yang", v.str);
    }
}
