package org.conf2.meta;

import org.junit.Test;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;

import java.util.Iterator;

/**
 *
 */
public class YangModuleTest {

    @Test
    public void testYang() {
        Module y = YangModule.YANG;
        Iterator<Meta> i = new MetaCollectionIterator(y);
        assertTrue(i.hasNext());
        Meta moduleContainer = i.next();
        assertNotNull(moduleContainer);
        assertEquals("module", moduleContainer.getIdent());

        Iterator<Meta> groupings = new MetaCollectionIterator(y.getGroupings());
        assertTrue(groupings.hasNext());
        Meta defHeader = groupings.next();
        assertEquals("def-header", defHeader.getIdent());
    }

    @Test
    public void testProxy() {
        Module y = YangModule.YANG;
        Container module = (Container) y.getFirstMeta();
        Iterator<Meta> metas = new MetaCollectionIterator(module);
        assertTrue(metas.hasNext());
        Meta ident = metas.next();
        assertEquals("ident", ident.getIdent());
    }
}
