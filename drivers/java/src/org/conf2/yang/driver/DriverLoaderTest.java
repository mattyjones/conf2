package org.conf2.yang.driver;

import org.conf2.yang.MetaUtil;
import org.conf2.yang.SimpleDataSource;
import org.conf2.yang.Module;
import org.junit.Test;

import java.io.IOException;
import java.io.InputStream;

import static org.junit.Assert.*;

public class DriverLoaderTest {

    @Test
    public void testLoadModule() throws IOException {
        Driver d = new Driver();
		DriverLoader loader = new DriverLoader(d);
		SimpleDataSource ds = new SimpleDataSource(DriverLoaderTest.class);
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
