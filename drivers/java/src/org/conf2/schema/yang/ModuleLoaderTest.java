package org.conf2.schema.yang;

import org.conf2.schema.MetaUtil;
import org.conf2.schema.SimpleStreamSource;
import org.conf2.schema.Module;
import org.conf2.schema.driver.Driver;
import org.junit.Test;

import java.io.IOException;

import static org.junit.Assert.*;

public class ModuleLoaderTest {

    @Test
    public void testLoadModule() throws IOException {
		ModuleLoader loader = new ModuleLoader(new Driver());
		SimpleStreamSource ds = new SimpleStreamSource(ModuleLoaderTest.class);
		Module simple = loader.loadModule(ds, "simple.yang");
		assertNotNull(simple);
		assertEquals("turing-machine", simple.getIdent());
		assertEquals(4, MetaUtil.collectionLength(simple.getTypedefs()));
		assertEquals(1, MetaUtil.collectionLength(simple.getGroupings()));
		assertEquals(1, MetaUtil.collectionLength(simple));
		assertEquals(2, MetaUtil.collectionLength(simple.getRpcs()));
		assertEquals(1, MetaUtil.collectionLength(simple.getNotifications()));
	}
}
